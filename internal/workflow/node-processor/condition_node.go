package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type conditionNode struct {
	node *domain.WorkflowNode
	*nodeLogger
}

func NewConditionNode(node *domain.WorkflowNode) *conditionNode {
	return &conditionNode{
		node:       node,
		nodeLogger: NewNodeLogger(node),
	}
}

// 条件节点没有任何操作
func (c *conditionNode) Run(ctx context.Context) error {
	c.AddOutput(ctx,
		c.node.Name,
		"完成",
	)
	return nil
}
