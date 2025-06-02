package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type executeSuccessNode struct {
	node *domain.WorkflowNode
	*nodeProcessor
	*nodeOutputer
}

func NewExecuteSuccessNode(node *domain.WorkflowNode) *executeSuccessNode {
	return &executeSuccessNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),
		nodeOutputer:  newNodeOutputer(),
	}
}

func (n *executeSuccessNode) Process(ctx context.Context) error {
	// 此类型节点不需要执行任何操作，直接返回

	return nil
}
