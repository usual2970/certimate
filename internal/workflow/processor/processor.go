package processor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
	nodes "github.com/usual2970/certimate/internal/workflow/node-processor"
)

type workflowProcessor struct {
	workflow *domain.Workflow
	logs     []domain.WorkflowRunLog
}

func NewWorkflowProcessor(workflow *domain.Workflow) *workflowProcessor {
	return &workflowProcessor{
		workflow: workflow,
		logs:     make([]domain.WorkflowRunLog, 0),
	}
}

func (w *workflowProcessor) Run(ctx context.Context) error {
	ctx = setContextWorkflowId(ctx, w.workflow.Id)
	return w.processNode(ctx, w.workflow.Content)
}

func (w *workflowProcessor) GetRunLogs() []domain.WorkflowRunLog {
	return w.logs
}

func (w *workflowProcessor) processNode(ctx context.Context, node *domain.WorkflowNode) error {
	current := node
	for current != nil {
		if current.Type == domain.WorkflowNodeTypeBranch || current.Type == domain.WorkflowNodeTypeExecuteResultBranch {
			for _, branch := range current.Branches {
				if err := w.processNode(ctx, &branch); err != nil {
					continue
				}
			}
		}

		var processor nodes.NodeProcessor
		var runErr error
		for {
			if current.Type != domain.WorkflowNodeTypeBranch && current.Type != domain.WorkflowNodeTypeExecuteResultBranch {
				processor, runErr = nodes.GetProcessor(current)
				if runErr != nil {
					break
				}

				runErr = processor.Run(ctx)
				log := processor.Log(ctx)
				if log != nil {
					w.logs = append(w.logs, *log)
				}
				if runErr != nil {
					break
				}
			}

			break
		}

		if runErr != nil && current.Next != nil && current.Next.Type != domain.WorkflowNodeTypeExecuteResultBranch {
			return runErr
		} else if runErr != nil && current.Next != nil && current.Next.Type == domain.WorkflowNodeTypeExecuteResultBranch {
			current = getBranchByType(current.Next.Branches, domain.WorkflowNodeTypeExecuteFailure)
		} else if runErr == nil && current.Next != nil && current.Next.Type == domain.WorkflowNodeTypeExecuteResultBranch {
			current = getBranchByType(current.Next.Branches, domain.WorkflowNodeTypeExecuteSuccess)
		} else {
			current = current.Next
		}
	}

	return nil
}

func setContextWorkflowId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, "workflow_id", id)
}

func getBranchByType(branches []domain.WorkflowNode, nodeType domain.WorkflowNodeType) *domain.WorkflowNode {
	for _, branch := range branches {
		if branch.Type == nodeType {
			return &branch
		}
	}
	return nil
}
