package nodeprocessor

import (
	"context"
	"errors"
	"fmt"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/domain/expr"
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
	nodeCfg := n.node.GetConfigForCondition()
	if nodeCfg.Expression == nil {
		n.logger.Info("without any conditions, enter this branch")
		return nil
	}

	rs, err := n.evalExpr(ctx, nodeCfg.Expression)
	if err != nil {
		n.logger.Warn(fmt.Sprintf("failed to eval condition expression: %w", err))
		return err
	}

	if rs.Value == false {
		n.logger.Info("condition not met, skip this branch")
		return errors.New("condition not met") // TODO: 错误处理
	} else {
		n.logger.Info("condition met, enter this branch")
	}

	return nil
}

func (n *conditionNode) evalExpr(ctx context.Context, expression expr.Expr) (*expr.EvalResult, error) {
	variables := GetNodeOutputs(ctx)
	return expression.Eval(variables)
}
