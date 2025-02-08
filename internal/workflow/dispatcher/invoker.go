package dispatcher

import (
	"context"
	"errors"

	"github.com/usual2970/certimate/internal/domain"
	nodes "github.com/usual2970/certimate/internal/workflow/node-processor"
)

type workflowInvoker struct {
	workflowId      string
	workflowContent *domain.WorkflowNode
	runId           string
	runLogs         []domain.WorkflowRunLog
}

func newWorkflowInvoker(data *WorkflowWorkerData) *workflowInvoker {
	if data == nil {
		panic("worker data is nil")
	}

	return &workflowInvoker{
		workflowId:      data.WorkflowId,
		workflowContent: data.WorkflowContent,
		runId:           data.RunId,
		runLogs:         make([]domain.WorkflowRunLog, 0),
	}
}

func (w *workflowInvoker) Invoke(ctx context.Context) error {
	ctx = context.WithValue(ctx, "workflow_id", w.workflowId)
	ctx = context.WithValue(ctx, "workflow_run_id", w.runId)
	return w.processNode(ctx, w.workflowContent)
}

func (w *workflowInvoker) GetLogs() []domain.WorkflowRunLog {
	return w.runLogs
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
					break
				}

				procErr = processor.Process(ctx)
				log := processor.GetLog(ctx)
				if log != nil {
					w.runLogs = append(w.runLogs, *log)
				}
				if procErr != nil {
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
