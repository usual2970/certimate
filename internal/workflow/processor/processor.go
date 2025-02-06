package processor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
	nodes "github.com/usual2970/certimate/internal/workflow/node-processor"
)

type workflowProcessor struct {
	workflowId      string
	workflowContent *domain.WorkflowNode
	runId           string
	runLogs         []domain.WorkflowRunLog
}

func NewWorkflowProcessor(workflowId string, workflowContent *domain.WorkflowNode, workflowRunId string) *workflowProcessor {
	return &workflowProcessor{
		workflowId:      workflowId,
		workflowContent: workflowContent,
		runId:           workflowRunId,
		runLogs:         make([]domain.WorkflowRunLog, 0),
	}
}

func (w *workflowProcessor) Process(ctx context.Context) error {
	ctx = context.WithValue(ctx, "workflow_id", w.workflowId)
	ctx = context.WithValue(ctx, "workflow_run_id", w.runId)
	return w.processNode(ctx, w.workflowContent)
}

func (w *workflowProcessor) GetLogs() []domain.WorkflowRunLog {
	return w.runLogs
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

				runErr = processor.Process(ctx)
				log := processor.GetLog(ctx)
				if log != nil {
					w.runLogs = append(w.runLogs, *log)
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
			current = w.getBranchByType(current.Next.Branches, domain.WorkflowNodeTypeExecuteFailure)
		} else if runErr == nil && current.Next != nil && current.Next.Type == domain.WorkflowNodeTypeExecuteResultBranch {
			current = w.getBranchByType(current.Next.Branches, domain.WorkflowNodeTypeExecuteSuccess)
		} else {
			current = current.Next
		}
	}

	return nil
}

func (w *workflowProcessor) getBranchByType(branches []domain.WorkflowNode, nodeType domain.WorkflowNodeType) *domain.WorkflowNode {
	for _, branch := range branches {
		if branch.Type == nodeType {
			return &branch
		}
	}
	return nil
}
