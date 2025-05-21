package domain

import (
	"encoding/json"
	"fmt"
)

type Value any

type (
	ComparisonOperator string
	LogicalOperator    string
	ValueType          string
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

	Number  ValueType = "number"
	String  ValueType = "string"
	Boolean ValueType = "boolean"
)

type EvalResult struct {
	Type  ValueType
	Value any
}

func (e *EvalResult) GetFloat64() (float64, error) {
	if e.Type != Number {
		return 0, fmt.Errorf("type mismatch: %s", e.Type)
	}
	switch v := e.Value.(type) {
	case int:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}

func (e *EvalResult) GreaterThan(other *EvalResult) (*EvalResult, error) {
	if e.Type != other.Type {
		return nil, fmt.Errorf("type mismatch: %s vs %s", e.Type, other.Type)
	}
	switch e.Type {
	case Number:

		left, err := e.GetFloat64()
		if err != nil {
			return nil, err
		}
		right, err := other.GetFloat64()
		if err != nil {
			return nil, err
		}

		return &EvalResult{
			Type:  Boolean,
			Value: left > right,
		}, nil
	case String:
		return &EvalResult{
			Type:  Boolean,
			Value: e.Value.(string) > other.Value.(string),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported type: %s", e.Type)
	}
}

func (e *EvalResult) GreaterOrEqual(other *EvalResult) (*EvalResult, error) {
	if e.Type != other.Type {
		return nil, fmt.Errorf("type mismatch: %s vs %s", e.Type, other.Type)
	}
	switch e.Type {
	case Number:
		left, err := e.GetFloat64()
		if err != nil {
			return nil, err
		}
		right, err := other.GetFloat64()
		if err != nil {
			return nil, err
		}
		return &EvalResult{
			Type:  Boolean,
			Value: left >= right,
		}, nil
	case String:
		return &EvalResult{
			Type:  Boolean,
			Value: e.Value.(string) >= other.Value.(string),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported type: %s", e.Type)
	}
}

func (e *EvalResult) LessThan(other *EvalResult) (*EvalResult, error) {
	if e.Type != other.Type {
		return nil, fmt.Errorf("type mismatch: %s vs %s", e.Type, other.Type)
	}
	switch e.Type {
	case Number:
		left, err := e.GetFloat64()
		if err != nil {
			return nil, err
		}
		right, err := other.GetFloat64()
		if err != nil {
			return nil, err
		}
		return &EvalResult{
			Type:  Boolean,
			Value: left < right,
		}, nil
	case String:
		return &EvalResult{
			Type:  Boolean,
			Value: e.Value.(string) < other.Value.(string),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported type: %s", e.Type)
	}
}

func (e *EvalResult) LessOrEqual(other *EvalResult) (*EvalResult, error) {
	if e.Type != other.Type {
		return nil, fmt.Errorf("type mismatch: %s vs %s", e.Type, other.Type)
	}
	switch e.Type {
	case Number:
		left, err := e.GetFloat64()
		if err != nil {
			return nil, err
		}
		right, err := other.GetFloat64()
		if err != nil {
			return nil, err
		}
		return &EvalResult{
			Type:  Boolean,
			Value: left <= right,
		}, nil
	case String:
		return &EvalResult{
			Type:  Boolean,
			Value: e.Value.(string) <= other.Value.(string),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported type: %s", e.Type)
	}
}

func (e *EvalResult) Equal(other *EvalResult) (*EvalResult, error) {
	if e.Type != other.Type {
		return nil, fmt.Errorf("type mismatch: %s vs %s", e.Type, other.Type)
	}
	switch e.Type {
	case Number:
		left, err := e.GetFloat64()
		if err != nil {
			return nil, err
		}
		right, err := other.GetFloat64()
		if err != nil {
			return nil, err
		}
		return &EvalResult{
			Type:  Boolean,
			Value: left == right,
		}, nil
	case String:
		return &EvalResult{
			Type:  Boolean,
			Value: e.Value.(string) == other.Value.(string),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported type: %s", e.Type)
	}
}

func (e *EvalResult) NotEqual(other *EvalResult) (*EvalResult, error) {
	if e.Type != other.Type {
		return nil, fmt.Errorf("type mismatch: %s vs %s", e.Type, other.Type)
	}
	switch e.Type {
	case Number:
		left, err := e.GetFloat64()
		if err != nil {
			return nil, err
		}
		right, err := other.GetFloat64()
		if err != nil {
			return nil, err
		}
		return &EvalResult{
			Type:  Boolean,
			Value: left != right,
		}, nil
	case String:
		return &EvalResult{
			Type:  Boolean,
			Value: e.Value.(string) != other.Value.(string),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported type: %s", e.Type)
	}
}

func (e *EvalResult) And(other *EvalResult) (*EvalResult, error) {
	if e.Type != other.Type {
		return nil, fmt.Errorf("type mismatch: %s vs %s", e.Type, other.Type)
	}
	switch e.Type {
	case Boolean:
		return &EvalResult{
			Type:  Boolean,
			Value: e.Value.(bool) && other.Value.(bool),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", e.Type)
	}
}

func (e *EvalResult) Or(other *EvalResult) (*EvalResult, error) {
	if e.Type != other.Type {
		return nil, fmt.Errorf("type mismatch: %s vs %s", e.Type, other.Type)
	}
	switch e.Type {
	case Boolean:
		return &EvalResult{
			Type:  Boolean,
			Value: e.Value.(bool) || other.Value.(bool),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", e.Type)
	}
}

func (e *EvalResult) Not() (*EvalResult, error) {
	if e.Type != Boolean {
		return nil, fmt.Errorf("type mismatch: %s", e.Type)
	}
	return &EvalResult{
		Type:  Boolean,
		Value: !e.Value.(bool),
	}, nil
}

func (e *EvalResult) Is(other *EvalResult) (*EvalResult, error) {
	if e.Type != other.Type {
		return nil, fmt.Errorf("type mismatch: %s vs %s", e.Type, other.Type)
	}
	switch e.Type {
	case Boolean:
		return &EvalResult{
			Type:  Boolean,
			Value: e.Value.(bool) == other.Value.(bool),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", e.Type)
	}
}

type Expr interface {
	GetType() string
	Eval(variables map[string]map[string]any) (*EvalResult, error)
}

type ConstExpr struct {
	Type      string    `json:"type"`
	Value     Value     `json:"value"`
	ValueType ValueType `json:"valueType"`
}

func (c ConstExpr) GetType() string { return c.Type }

func (c ConstExpr) Eval(variables map[string]map[string]any) (*EvalResult, error) {
	return &EvalResult{
		Type:  c.ValueType,
		Value: c.Value,
	}, nil
}

type VarExpr struct {
	Type     string                      `json:"type"`
	Selector WorkflowNodeIOValueSelector `json:"selector"`
}

func (v VarExpr) GetType() string { return v.Type }

func (v VarExpr) Eval(variables map[string]map[string]any) (*EvalResult, error) {
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
		return nil, fmt.Errorf("variable %s not found in node %s", v.Selector.Name, v.Selector.Id)
	}
	return &EvalResult{
		Type:  v.Selector.Type,
		Value: variables[v.Selector.Id][v.Selector.Name],
	}, nil
}

type CompareExpr struct {
	Type  string             `json:"type"` // compare
	Op    ComparisonOperator `json:"op"`
	Left  Expr               `json:"left"`
	Right Expr               `json:"right"`
}

func (c CompareExpr) GetType() string { return c.Type }

func (c CompareExpr) Eval(variables map[string]map[string]any) (*EvalResult, error) {
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
		return left.GreaterThan(right)
	case LessThan:
		return left.LessThan(right)
	case GreaterOrEqual:
		return left.GreaterOrEqual(right)
	case LessOrEqual:
		return left.LessOrEqual(right)
	case Equal:
		return left.Equal(right)
	case NotEqual:
		return left.NotEqual(right)
	case Is:
		return left.Is(right)
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

func (l LogicalExpr) Eval(variables map[string]map[string]any) (*EvalResult, error) {
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
		return left.And(right)
	case Or:
		return left.Or(right)
	default:
		return nil, fmt.Errorf("unknown operator: %s", l.Op)
	}
}

type NotExpr struct {
	Type string `json:"type"` // not
	Expr Expr   `json:"expr"`
}

func (n NotExpr) GetType() string { return n.Type }

func (n NotExpr) Eval(variables map[string]map[string]any) (*EvalResult, error) {
	inner, err := n.Expr.Eval(variables)
	if err != nil {
		return nil, err
	}
	return inner.Not()
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
