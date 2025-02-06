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

func (n *conditionNode) Process(ctx context.Context) error {
	// 此类型节点不需要执行任何操作，直接返回
	n.AddOutput(ctx, n.node.Name, "完成")

	return nil
}
