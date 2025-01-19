package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type executeSuccessNode struct {
	node *domain.WorkflowNode
	*nodeLogger
}

func NewExecuteSuccessNode(node *domain.WorkflowNode) *executeSuccessNode {
	return &executeSuccessNode{
		node:       node,
		nodeLogger: NewNodeLogger(node),
	}
}

func (e *executeSuccessNode) Run(ctx context.Context) error {
	e.AddOutput(ctx,
		e.node.Name,
		"进入执行成功分支",
	)
	return nil
}
