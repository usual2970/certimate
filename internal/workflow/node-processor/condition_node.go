package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type conditionNode struct {
	node *domain.WorkflowNode
	*nodeProcessor
}

func NewConditionNode(node *domain.WorkflowNode) *conditionNode {
	return &conditionNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),
	}
}

func (n *conditionNode) Process(ctx context.Context) error {
	// 此类型节点不需要执行任何操作，直接返回
	return nil
}
