package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type startNode struct {
	node *domain.WorkflowNode
	*nodeProcessor
	*nodeOutputer
}

func NewStartNode(node *domain.WorkflowNode) *startNode {
	return &startNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),
		nodeOutputer:  newNodeOutputer(),
	}
}

func (n *startNode) Process(ctx context.Context) error {
	// 此类型节点不需要执行任何操作，直接返回
	n.logger.Info("workflow is started")

	return nil
}
