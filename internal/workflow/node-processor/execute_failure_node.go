package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type executeFailureNode struct {
	node *domain.WorkflowNode
	*nodeProcessor
	*nodeOutputer
}

func NewExecuteFailureNode(node *domain.WorkflowNode) *executeFailureNode {
	return &executeFailureNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),
		nodeOutputer:  newNodeOutputer(),
	}
}

func (n *executeFailureNode) Process(ctx context.Context) error {
	// 此类型节点不需要执行任何操作，直接返回

	return nil
}
