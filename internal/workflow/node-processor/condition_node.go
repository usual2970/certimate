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
		n.logger.Info("no condition found, continue to next node")
		return nil
	}

	rs, err := n.eval(ctx, nodeConfig.Expression)
	if err != nil {
		n.logger.Warn("failed to eval expression: " + err.Error())
		return err
	}

	if rs.Value == false {
		n.logger.Info("condition not met, skip this branch")
		return errors.New("condition not met")
	}

	n.logger.Info("condition met, continue to next node")
	return nil
}

func (n *conditionNode) eval(ctx context.Context, expression domain.Expr) (*domain.EvalResult, error) {
	variables := GetNodeOutputs(ctx)
	return expression.Eval(variables)
}
