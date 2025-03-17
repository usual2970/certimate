package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type executeSuccessNode struct {
	node *domain.WorkflowNode
	*nodeProcessor
}

func NewExecuteSuccessNode(node *domain.WorkflowNode) *executeSuccessNode {
	return &executeSuccessNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),
	}
}

func (n *executeSuccessNode) Process(ctx context.Context) error {
	// 此类型节点不需要执行任何操作，直接返回
	n.logger.Info("the previous node execution was succeeded")

	return nil
}
