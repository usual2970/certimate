package dispatcher

import (
	"context"
	"errors"
	"log/slog"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/logging"
	nodes "github.com/usual2970/certimate/internal/workflow/node-processor"
)

type workflowInvoker struct {
	workflowId      string
	workflowContent *domain.WorkflowNode
	runId           string
	logs            []domain.WorkflowLog

	workflowLogRepo workflowLogRepository
}

func newWorkflowInvokerWithData(workflowLogRepo workflowLogRepository, data *WorkflowWorkerData) *workflowInvoker {
	if data == nil {
		panic("worker data is nil")
	}

	return &workflowInvoker{
		workflowId:      data.WorkflowId,
		workflowContent: data.WorkflowContent,
		runId:           data.RunId,
		logs:            make([]domain.WorkflowLog, 0),

		workflowLogRepo: workflowLogRepo,
	}
}

func (w *workflowInvoker) Invoke(ctx context.Context) error {
	ctx = context.WithValue(ctx, "workflow_id", w.workflowId)
	ctx = context.WithValue(ctx, "workflow_run_id", w.runId)
	return w.processNode(ctx, w.workflowContent)
}

func (w *workflowInvoker) GetLogs() domain.WorkflowLogs {
	return w.logs
}

func (w *workflowInvoker) processNode(ctx context.Context, node *domain.WorkflowNode) error {
	current := node
	for current != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if current.Type == domain.WorkflowNodeTypeBranch || current.Type == domain.WorkflowNodeTypeExecuteResultBranch {
			for _, branch := range current.Branches {
				if err := w.processNode(ctx, &branch); err != nil {
					// 并行分支的某一分支发生错误时，忽略此错误，继续执行其他分支
					if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
						continue
					}
					return err
				}
			}
		}

		var processor nodes.NodeProcessor
		var procErr error
		for {
			if current.Type != domain.WorkflowNodeTypeBranch && current.Type != domain.WorkflowNodeTypeExecuteResultBranch {
				processor, procErr = nodes.GetProcessor(current)
				if procErr != nil {
					panic(procErr)
				}

				processor.SetLogger(slog.New(logging.NewHookHandler(&logging.HookHandlerOptions{
					Level: slog.LevelDebug,
					WriteFunc: func(ctx context.Context, record *logging.Record) error {
						log := domain.WorkflowLog{}
						log.WorkflowId = w.workflowId
						log.RunId = w.runId
						log.NodeId = current.Id
						log.NodeName = current.Name
						log.Timestamp = record.Time.UnixMilli()
						log.Level = record.Level.String()
						log.Message = record.Message
						log.Data = record.Data
						log.CreatedAt = record.Time
						if _, err := w.workflowLogRepo.Save(ctx, &log); err != nil {
							return err
						}

						w.logs = append(w.logs, log)
						return nil
					},
				})))

				procErr = processor.Process(ctx)
				if procErr != nil {
					processor.GetLogger().Error(procErr.Error())
					break
				}
			}

			break
		}

		// TODO: 优化可读性
		if procErr != nil && current.Next != nil && current.Next.Type != domain.WorkflowNodeTypeExecuteResultBranch {
			return procErr
		} else if procErr != nil && current.Next != nil && current.Next.Type == domain.WorkflowNodeTypeExecuteResultBranch {
			current = w.getBranchByType(current.Next.Branches, domain.WorkflowNodeTypeExecuteFailure)
		} else if procErr == nil && current.Next != nil && current.Next.Type == domain.WorkflowNodeTypeExecuteResultBranch {
			current = w.getBranchByType(current.Next.Branches, domain.WorkflowNodeTypeExecuteSuccess)
		} else {
			current = current.Next
		}
	}

	return nil
}

func (w *workflowInvoker) getBranchByType(branches []domain.WorkflowNode, nodeType domain.WorkflowNodeType) *domain.WorkflowNode {
	for _, branch := range branches {
		if branch.Type == nodeType {
			return &branch
		}
	}
	return nil
}
