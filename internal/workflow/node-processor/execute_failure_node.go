package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type executeFailureNode struct {
	node *domain.WorkflowNode
	*nodeLogger
}

func NewExecuteFailureNode(node *domain.WorkflowNode) *executeFailureNode {
	return &executeFailureNode{
		node:       node,
		nodeLogger: NewNodeLogger(node),
	}
}

func (e *executeFailureNode) Run(ctx context.Context) error {
	e.AddOutput(ctx,
		e.node.Name,
		"进入执行失败分支",
	)
	return nil
}
