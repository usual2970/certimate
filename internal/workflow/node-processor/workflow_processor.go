package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
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
	ctx = WithWorkflowId(ctx, w.workflow.Id)
	return w.runNode(ctx, w.workflow.Content)
}

func (w *workflowProcessor) runNode(ctx context.Context, node *domain.WorkflowNode) error {
	current := node
	for current != nil {
		if current.Type == domain.WorkflowNodeTypeBranch {
			for _, branch := range current.Branches {
				if err := w.runNode(ctx, &branch); err != nil {
					continue
				}
			}
		}

		if current.Type != domain.WorkflowNodeTypeBranch {
			processor, err := GetProcessor(current)
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

func WithWorkflowId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, "workflow_id", id)
}

func GetWorkflowId(ctx context.Context) string {
	return ctx.Value("workflow_id").(string)
}
