package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/xtime"
)

type conditionNode struct {
	node *domain.WorkflowNode
	*Logger
}

func NewConditionNode(node *domain.WorkflowNode) *conditionNode {
	return &conditionNode{
		node:   node,
		Logger: NewLogger(node),
	}
}

// 条件节点没有任何操作
func (c *conditionNode) Run(ctx context.Context) error {
	c.AddOutput(ctx, xtime.BeijingTimeStr(),
		c.node.Name,
		"完成",
	)
	return nil
}
