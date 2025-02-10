package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type startNode struct {
	node *domain.WorkflowNode
	*nodeLogger
}

func NewStartNode(node *domain.WorkflowNode) *startNode {
	return &startNode{
		node:       node,
		nodeLogger: newNodeLogger(node),
	}
}

func (n *startNode) Process(ctx context.Context) error {
	// 此类型节点不需要执行任何操作，直接返回
	n.AppendLogRecord(ctx, domain.WorkflowRunLogLevelInfo, "进入开始节点")

	return nil
}
