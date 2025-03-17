package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type executeFailureNode struct {
	node *domain.WorkflowNode
	*nodeProcessor
}

func NewExecuteFailureNode(node *domain.WorkflowNode) *executeFailureNode {
	return &executeFailureNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),
	}
}

func (n *executeFailureNode) Process(ctx context.Context) error {
	// 此类型节点不需要执行任何操作，直接返回
	n.logger.Info("the previous node execution was failed")

	return nil
}
