package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type startNode struct {
	node *domain.WorkflowNode
	*Logger
}

func NewStartNode(node *domain.WorkflowNode) *startNode {
	return &startNode{
		node:   node,
		Logger: NewLogger(node),
	}
}

// 开始节点没有任何操作
func (s *startNode) Run(ctx context.Context) error {
	s.AddOutput(ctx,
		s.node.Name,
		"完成",
	)
	return nil
}
