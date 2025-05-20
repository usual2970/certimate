package domain

import (
	"encoding/json"
	"fmt"
)

type Value any

type (
	ComparisonOperator string
	LogicalOperator    string
)

const (
	GreaterThan    ComparisonOperator = ">"
	LessThan       ComparisonOperator = "<"
	GreaterOrEqual ComparisonOperator = ">="
	LessOrEqual    ComparisonOperator = "<="
	Equal          ComparisonOperator = "=="
	NotEqual       ComparisonOperator = "!="
	Is             ComparisonOperator = "is"

	And LogicalOperator = "and"
	Or  LogicalOperator = "or"
	Not LogicalOperator = "not"
)

type Expr interface {
	GetType() string
	Eval(variables map[string]map[string]any) (any, error)
}

type ConstExpr struct {
	Type  string `json:"type"`
	Value Value  `json:"value"`
}

func (c ConstExpr) GetType() string { return c.Type }

type VarExpr struct {
	Type     string                      `json:"type"`
	Selector WorkflowNodeIOValueSelector `json:"selector"`
}

func (v VarExpr) GetType() string { return v.Type }

func (v VarExpr) Eval(variables map[string]map[string]any) (any, error) {
	if v.Selector.Id == "" {
		return nil, fmt.Errorf("node id is empty")
	}
	if v.Selector.Name == "" {
		return nil, fmt.Errorf("name is empty")
	}

	if _, ok := variables[v.Selector.Id]; !ok {
		return nil, fmt.Errorf("node %s not found", v.Selector.Id)
	}

	if _, ok := variables[v.Selector.Id][v.Selector.Name]; !ok {
		return nil, fmt.Errorf("variable %s not found in node %s", v.Selector.Name, v.Selector.NodeId)
	}

	return variables[v.Selector.Id][v.Selector.Name], nil
}

type CompareExpr struct {
	Type  string             `json:"type"` // compare
	Op    ComparisonOperator `json:"op"`
	Left  Expr               `json:"left"`
	Right Expr               `json:"right"`
}

func (c CompareExpr) GetType() string { return c.Type }

func (c CompareExpr) Eval(variables map[string]map[string]any) (any, error) {
	left, err := c.Left.Eval(variables)
	if err != nil {
		return nil, err
	}
	right, err := c.Right.Eval(variables)
	if err != nil {
		return nil, err
	}

	switch c.Op {
	case GreaterThan:
		return left.(float64) > right.(float64), nil
	case LessThan:
		return left.(float64) < right.(float64), nil
	case GreaterOrEqual:
		return left.(float64) >= right.(float64), nil
	case LessOrEqual:
		return left.(float64) <= right.(float64), nil
	case Equal:
		return left == right, nil
	case NotEqual:
		return left != right, nil
	case Is:
		return left == right, nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", c.Op)
	}
}

type LogicalExpr struct {
	Type  string          `json:"type"` // logical
	Op    LogicalOperator `json:"op"`
	Left  Expr            `json:"left"`
	Right Expr            `json:"right"`
}

func (l LogicalExpr) GetType() string { return l.Type }

func (l LogicalExpr) Eval(variables map[string]map[string]any) (any, error) {
	left, err := l.Left.Eval(variables)
	if err != nil {
		return nil, err
	}
	right, err := l.Right.Eval(variables)
	if err != nil {
		return nil, err
	}

	switch l.Op {
	case And:
		return left.(bool) && right.(bool), nil
	case Or:
		return left.(bool) || right.(bool), nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", l.Op)
	}
}

type NotExpr struct {
	Type string `json:"type"` // not
	Expr Expr   `json:"expr"`
}

func (n NotExpr) GetType() string { return n.Type }

func (n NotExpr) Eval(variables map[string]map[string]any) (any, error) {
	inner, err := n.Expr.Eval(variables)
	if err != nil {
		return nil, err
	}
	return !inner.(bool), nil
}

type rawExpr struct {
	Type string `json:"type"`
}

func MarshalExpr(e Expr) ([]byte, error) {
	return json.Marshal(e)
}

func UnmarshalExpr(data []byte) (Expr, error) {
	var typ rawExpr
	if err := json.Unmarshal(data, &typ); err != nil {
		return nil, err
	}

	switch typ.Type {
	case "const":
		var e ConstExpr
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return e, nil
	case "var":
		var e VarExpr
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return e, nil
	case "compare":
		var e CompareExprRaw
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return e.ToCompareExpr()
	case "logical":
		var e LogicalExprRaw
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return e.ToLogicalExpr()
	case "not":
		var e NotExprRaw
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return e.ToNotExpr()
	default:
		return nil, fmt.Errorf("unknown expr type: %s", typ.Type)
	}
}

type CompareExprRaw struct {
	Type  string             `json:"type"`
	Op    ComparisonOperator `json:"op"`
	Left  json.RawMessage    `json:"left"`
	Right json.RawMessage    `json:"right"`
}

func (r CompareExprRaw) ToCompareExpr() (CompareExpr, error) {
	leftExpr, err := UnmarshalExpr(r.Left)
	if err != nil {
		return CompareExpr{}, err
	}
	rightExpr, err := UnmarshalExpr(r.Right)
	if err != nil {
		return CompareExpr{}, err
	}
	return CompareExpr{
		Type:  r.Type,
		Op:    r.Op,
		Left:  leftExpr,
		Right: rightExpr,
	}, nil
}

type LogicalExprRaw struct {
	Type  string          `json:"type"`
	Op    LogicalOperator `json:"op"`
	Left  json.RawMessage `json:"left"`
	Right json.RawMessage `json:"right"`
}

func (r LogicalExprRaw) ToLogicalExpr() (LogicalExpr, error) {
	left, err := UnmarshalExpr(r.Left)
	if err != nil {
		return LogicalExpr{}, err
	}
	right, err := UnmarshalExpr(r.Right)
	if err != nil {
		return LogicalExpr{}, err
	}
	return LogicalExpr{
		Type:  r.Type,
		Op:    r.Op,
		Left:  left,
		Right: right,
	}, nil
}

type NotExprRaw struct {
	Type string          `json:"type"`
	Expr json.RawMessage `json:"expr"`
}

func (r NotExprRaw) ToNotExpr() (NotExpr, error) {
	inner, err := UnmarshalExpr(r.Expr)
	if err != nil {
		return NotExpr{}, err
	}
	return NotExpr{
		Type: r.Type,
		Expr: inner,
	}, nil
}
