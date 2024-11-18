package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/xtime"
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
	s.AddOutput(ctx, xtime.BeijingTimeStr(),
		s.node.Name,
		"完成",
	)
	return nil
}
