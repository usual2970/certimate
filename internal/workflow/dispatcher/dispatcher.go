package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/utils/slices"
)

var maxWorkers = 16

func init() {
	envMaxWorkers := os.Getenv("CERTIMATE_WORKFLOW_MAX_WORKERS")
	if n, err := strconv.Atoi(envMaxWorkers); err != nil && n > 0 {
		maxWorkers = n
	}
}

type workflowWorker struct {
	Data   *WorkflowWorkerData
	Cancel context.CancelFunc
}

type WorkflowWorkerData struct {
	WorkflowId      string
	WorkflowContent *domain.WorkflowNode
	RunId           string
}

type WorkflowDispatcher struct {
	semaphore chan struct{}

	queue      []*WorkflowWorkerData
	queueMutex sync.Mutex

	workers     map[string]*workflowWorker // key: WorkflowId
	workerIdMap map[string]string          // key: RunId, value: WorkflowId
	workerMutex sync.Mutex

	chWork  chan *WorkflowWorkerData
	chCandi chan struct{}

	wg sync.WaitGroup

	workflowRepo    workflowRepository
	workflowRunRepo workflowRunRepository
}

func newWorkflowDispatcher(workflowRepo workflowRepository, workflowRunRepo workflowRunRepository) *WorkflowDispatcher {
	dispatcher := &WorkflowDispatcher{
		semaphore: make(chan struct{}, maxWorkers),

		queue:      make([]*WorkflowWorkerData, 0),
		queueMutex: sync.Mutex{},

		workers:     make(map[string]*workflowWorker),
		workerIdMap: make(map[string]string),
		workerMutex: sync.Mutex{},

		chWork:  make(chan *WorkflowWorkerData),
		chCandi: make(chan struct{}, 1),

		workflowRepo:    workflowRepo,
		workflowRunRepo: workflowRunRepo,
	}

	go func() {
		for {
			select {
			case <-dispatcher.chWork:
				dispatcher.dequeueWorker()

			case <-dispatcher.chCandi:
				dispatcher.dequeueWorker()
			}
		}
	}()

	return dispatcher
}

func (w *WorkflowDispatcher) Dispatch(data *WorkflowWorkerData) {
	if data == nil {
		panic("worker data is nil")
	}

	w.enqueueWorker(data)

	select {
	case w.chWork <- data:
	default:
	}
}

func (w *WorkflowDispatcher) Cancel(runId string) {
	hasWorker := false

	// 取消正在执行的 WorkflowRun
	w.workerMutex.Lock()
	if workflowId, ok := w.workerIdMap[runId]; ok {
		if worker, ok := w.workers[workflowId]; ok {
			hasWorker = true
			worker.Cancel()
			delete(w.workers, workflowId)
			delete(w.workerIdMap, runId)
		}
	}
	w.workerMutex.Unlock()

	// 移除排队中的 WorkflowRun
	w.queueMutex.Lock()
	w.queue = slices.Filter(w.queue, func(d *WorkflowWorkerData) bool {
		return d.RunId != runId
	})
	w.queueMutex.Unlock()

	// 已挂起，查询 WorkflowRun 并更新其状态为 Canceled
	if !hasWorker {
		if run, err := w.workflowRunRepo.GetById(context.Background(), runId); err == nil {
			if run.Status == domain.WorkflowRunStatusTypePending || run.Status == domain.WorkflowRunStatusTypeRunning {
				run.Status = domain.WorkflowRunStatusTypeCanceled
				w.workflowRunRepo.Save(context.Background(), run)
			}
		}
	}
}

func (w *WorkflowDispatcher) Shutdown() {
	// 清空排队中的 WorkflowRun
	w.queueMutex.Lock()
	w.queue = make([]*WorkflowWorkerData, 0)
	w.queueMutex.Unlock()

	// 等待所有正在执行的 WorkflowRun 完成
	w.workerMutex.Lock()
	for _, worker := range w.workers {
		worker.Cancel()
		delete(w.workers, worker.Data.WorkflowId)
		delete(w.workerIdMap, worker.Data.RunId)
	}
	w.workerMutex.Unlock()
	w.wg.Wait()
}

func (w *WorkflowDispatcher) enqueueWorker(data *WorkflowWorkerData) {
	w.queueMutex.Lock()
	defer w.queueMutex.Unlock()
	w.queue = append(w.queue, data)
}

func (w *WorkflowDispatcher) dequeueWorker() {
	for {
		select {
		case w.semaphore <- struct{}{}:
		default:
			// 达到最大并发数
			return
		}

		w.queueMutex.Lock()
		if len(w.queue) == 0 {
			w.queueMutex.Unlock()
			<-w.semaphore
			return
		}

		data := w.queue[0]
		w.queue = w.queue[1:]
		w.queueMutex.Unlock()

		// 检查是否有相同 WorkflowId 的 WorkflowRun 正在执行
		// 如果有，则重新排队，以保证同一个工作流同一时间内只有一个正在执行
		// 即不同 WorkflowId 的任务并行化，相同 WorkflowId 的任务串行化
		w.workerMutex.Lock()
		if _, exists := w.workers[data.WorkflowId]; exists {
			w.queueMutex.Lock()
			w.queue = append(w.queue, data)
			w.queueMutex.Unlock()
			w.workerMutex.Unlock()

			<-w.semaphore

			continue
		}

		ctx, cancel := context.WithCancel(context.Background())
		w.workers[data.WorkflowId] = &workflowWorker{data, cancel}
		w.workerIdMap[data.RunId] = data.WorkflowId
		w.workerMutex.Unlock()

		w.wg.Add(1)
		go w.work(ctx, data)
	}
}

func (w *WorkflowDispatcher) work(ctx context.Context, data *WorkflowWorkerData) {
	defer func() {
		<-w.semaphore
		w.workerMutex.Lock()
		delete(w.workers, data.WorkflowId)
		delete(w.workerIdMap, data.RunId)
		w.workerMutex.Unlock()

		w.wg.Done()

		// 尝试取出排队中的其他 WorkflowRun 继续执行
		select {
		case w.chCandi <- struct{}{}:
		default:
		}
	}()

	// 查询 WorkflowRun
	run, err := w.workflowRunRepo.GetById(ctx, data.RunId)
	if err != nil {
		if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			app.GetLogger().Error(fmt.Sprintf("failed to get workflow run #%s", data.RunId), "err", err)
		}
		return
	} else if run.Status != domain.WorkflowRunStatusTypePending {
		return
	} else if ctx.Err() != nil {
		run.Status = domain.WorkflowRunStatusTypeCanceled
		w.workflowRunRepo.Save(ctx, run)
		return
	}

	// 更新 WorkflowRun 状态为 Running
	run.Status = domain.WorkflowRunStatusTypeRunning
	if _, err := w.workflowRunRepo.Save(ctx, run); err != nil {
		if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			panic(err)
		}
		return
	}

	// 执行工作流
	invoker := newWorkflowInvokerWithData(w.workflowRunRepo, data)
	if runErr := invoker.Invoke(ctx); runErr != nil {
		if errors.Is(runErr, context.Canceled) {
			run.Status = domain.WorkflowRunStatusTypeCanceled
			run.Logs = invoker.GetLogs()
		} else {
			run.Status = domain.WorkflowRunStatusTypeFailed
			run.EndedAt = time.Now()
			run.Logs = invoker.GetLogs()
			run.Error = runErr.Error()
		}

		if _, err := w.workflowRunRepo.Save(ctx, run); err != nil {
			if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
				panic(err)
			}
		}

		return
	}

	// 更新 WorkflowRun 状态为 Succeeded/Failed
	run.EndedAt = time.Now()
	run.Logs = invoker.GetLogs()
	run.Error = domain.WorkflowRunLogs(invoker.GetLogs()).ErrorString()
	if run.Error == "" {
		run.Status = domain.WorkflowRunStatusTypeSucceeded
	} else {
		run.Status = domain.WorkflowRunStatusTypeFailed
	}
	if _, err := w.workflowRunRepo.Save(ctx, run); err != nil {
		if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			panic(err)
		}
	}
}
