package nodeprocessor

import (
	"context"
	"errors"

	"github.com/usual2970/certimate/internal/domain"
)

type conditionNode struct {
	node *domain.WorkflowNode
	*nodeProcessor
	*nodeOutputer
}

func NewConditionNode(node *domain.WorkflowNode) *conditionNode {
	return &conditionNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),
		nodeOutputer:  newNodeOutputer(),
	}
}

func (n *conditionNode) Process(ctx context.Context) error {
	n.logger.Info("enter condition node: " + n.node.Name)

	nodeConfig := n.node.GetConfigForCondition()
	if nodeConfig.Expression == nil {
		return nil
	}
	return nil
}

func (n *conditionNode) eval(ctx context.Context, expression domain.Expr) (any, error) {
	switch expr:=expression.(type) {
	case domain.CompareExpr:
		left,err:= n.eval(ctx, expr.Left)
		if err != nil {
			return nil, err
		}
		right,err:= n.eval(ctx, expr.Right)
		if err != nil {
			return nil, err
		}

	case domain.LogicalExpr:
	case domain.NotExpr:
	case domain.VarExpr:
	case domain.ConstExpr:
	}
	return false, errors.New("unknown expression type")
}
