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

func (w *workflowProcessor) Log(ctx context.Context) []domain.WorkflowRunLog {
	return w.logs
}

func (w *workflowProcessor) Run(ctx context.Context) error {
	ctx = setContextWorkflowId(ctx, w.workflow.Id)
	return w.processNode(ctx, w.workflow.Content)
}

func (w *workflowProcessor) processNode(ctx context.Context, node *domain.WorkflowNode) error {
	current := node
	for current != nil {
		if current.Type == domain.WorkflowNodeTypeBranch {
			for _, branch := range current.Branches {
				if err := w.processNode(ctx, &branch); err != nil {
					continue
				}
			}
		}

		if current.Type != domain.WorkflowNodeTypeBranch {
			processor, err := nodes.GetProcessor(current)
			if err != nil {
				return err
			}

			err = processor.Run(ctx)

			log := processor.Log(ctx)
			if log != nil {
				w.logs = append(w.logs, *log)
			}

			if err != nil {
				return err
			}
		}
		current = current.Next

	}
	return nil
}

func setContextWorkflowId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, "workflow_id", id)
}
