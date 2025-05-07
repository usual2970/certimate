package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	sliceutil "github.com/usual2970/certimate/internal/pkg/utils/slice"
)

var maxWorkers = 1

func init() {
	envMaxWorkers := os.Getenv("CERTIMATE_WORKFLOW_MAX_WORKERS")
	if n, err := strconv.Atoi(envMaxWorkers); err != nil && n > 0 {
		maxWorkers = n
	} else {
		maxWorkers = runtime.GOMAXPROCS(0)
		if maxWorkers == 0 {
			maxWorkers = max(1, runtime.NumCPU())
		}
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
	workflowLogRepo workflowLogRepository
}

func newWorkflowDispatcher(workflowRepo workflowRepository, workflowRunRepo workflowRunRepository, workflowLogRepo workflowLogRepository) *WorkflowDispatcher {
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
		workflowLogRepo: workflowLogRepo,
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

func (d *WorkflowDispatcher) Dispatch(data *WorkflowWorkerData) {
	if data == nil {
		panic("worker data is nil")
	}

	d.enqueueWorker(data)

	select {
	case d.chWork <- data:
	default:
	}
}

func (d *WorkflowDispatcher) Cancel(runId string) {
	hasWorker := false

	// 取消正在执行的 WorkflowRun
	d.workerMutex.Lock()
	if workflowId, ok := d.workerIdMap[runId]; ok {
		if worker, ok := d.workers[workflowId]; ok {
			hasWorker = true
			worker.Cancel()
			delete(d.workers, workflowId)
			delete(d.workerIdMap, runId)
		}
	}
	d.workerMutex.Unlock()

	// 移除排队中的 WorkflowRun
	d.queueMutex.Lock()
	d.queue = sliceutil.Filter(d.queue, func(d *WorkflowWorkerData) bool {
		return d.RunId != runId
	})
	d.queueMutex.Unlock()

	// 已挂起，查询 WorkflowRun 并更新其状态为 Canceled
	if !hasWorker {
		if run, err := d.workflowRunRepo.GetById(context.Background(), runId); err == nil {
			if run.Status == domain.WorkflowRunStatusTypePending || run.Status == domain.WorkflowRunStatusTypeRunning {
				run.Status = domain.WorkflowRunStatusTypeCanceled
				d.workflowRunRepo.Save(context.Background(), run)
			}
		}
	}
}

func (d *WorkflowDispatcher) Shutdown() {
	// 清空排队中的 WorkflowRun
	d.queueMutex.Lock()
	d.queue = make([]*WorkflowWorkerData, 0)
	d.queueMutex.Unlock()

	// 等待所有正在执行的 WorkflowRun 完成
	d.workerMutex.Lock()
	for _, worker := range d.workers {
		worker.Cancel()
		delete(d.workers, worker.Data.WorkflowId)
		delete(d.workerIdMap, worker.Data.RunId)
	}
	d.workerMutex.Unlock()
	d.wg.Wait()
}

func (d *WorkflowDispatcher) enqueueWorker(data *WorkflowWorkerData) {
	d.queueMutex.Lock()
	defer d.queueMutex.Unlock()
	d.queue = append(d.queue, data)
}

func (d *WorkflowDispatcher) dequeueWorker() {
	for {
		select {
		case d.semaphore <- struct{}{}:
		default:
			// 达到最大并发数
			return
		}

		d.queueMutex.Lock()
		if len(d.queue) == 0 {
			d.queueMutex.Unlock()
			<-d.semaphore
			return
		}

		data := d.queue[0]
		d.queue = d.queue[1:]
		d.queueMutex.Unlock()

		// 检查是否有相同 WorkflowId 的 WorkflowRun 正在执行
		// 如果有，则重新排队，以保证同一个工作流同一时间内只有一个正在执行
		// 即不同 WorkflowId 的任务并行化，相同 WorkflowId 的任务串行化
		d.workerMutex.Lock()
		if _, exists := d.workers[data.WorkflowId]; exists {
			d.queueMutex.Lock()
			d.queue = append(d.queue, data)
			d.queueMutex.Unlock()
			d.workerMutex.Unlock()

			<-d.semaphore

			continue
		}

		ctx, cancel := context.WithCancel(context.Background())
		d.workers[data.WorkflowId] = &workflowWorker{data, cancel}
		d.workerIdMap[data.RunId] = data.WorkflowId
		d.workerMutex.Unlock()

		d.wg.Add(1)
		go d.work(ctx, data)
	}
}

func (d *WorkflowDispatcher) work(ctx context.Context, data *WorkflowWorkerData) {
	defer func() {
		<-d.semaphore
		d.workerMutex.Lock()
		delete(d.workers, data.WorkflowId)
		delete(d.workerIdMap, data.RunId)
		d.workerMutex.Unlock()

		d.wg.Done()

		// 尝试取出排队中的其他 WorkflowRun 继续执行
		select {
		case d.chCandi <- struct{}{}:
		default:
		}
	}()

	// 查询 WorkflowRun
	run, err := d.workflowRunRepo.GetById(ctx, data.RunId)
	if err != nil {
		if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			app.GetLogger().Error(fmt.Sprintf("failed to get workflow run #%s", data.RunId), "err", err)
		}
		return
	} else if run.Status != domain.WorkflowRunStatusTypePending {
		return
	} else if ctx.Err() != nil {
		run.Status = domain.WorkflowRunStatusTypeCanceled
		d.workflowRunRepo.Save(ctx, run)
		return
	}

	// 更新 WorkflowRun 状态为 Running
	run.Status = domain.WorkflowRunStatusTypeRunning
	if _, err := d.workflowRunRepo.Save(ctx, run); err != nil {
		if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			panic(err)
		}
		return
	}

	// 执行工作流
	invoker := newWorkflowInvokerWithData(d.workflowLogRepo, data)
	if runErr := invoker.Invoke(ctx); runErr != nil {
		if errors.Is(runErr, context.Canceled) {
			run.Status = domain.WorkflowRunStatusTypeCanceled
		} else {
			run.Status = domain.WorkflowRunStatusTypeFailed
			run.EndedAt = time.Now()
			run.Error = runErr.Error()
		}

		if _, err := d.workflowRunRepo.Save(ctx, run); err != nil {
			if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
				panic(err)
			}
		}

		return
	}

	// 更新 WorkflowRun 状态为 Succeeded/Failed
	run.EndedAt = time.Now()
	run.Error = invoker.GetLogs().ErrorString()
	if run.Error == "" {
		run.Status = domain.WorkflowRunStatusTypeSucceeded
	} else {
		run.Status = domain.WorkflowRunStatusTypeFailed
	}
	if _, err := d.workflowRunRepo.Save(ctx, run); err != nil {
		if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			panic(err)
		}
	}
}
