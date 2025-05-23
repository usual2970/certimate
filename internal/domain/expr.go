package domain

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Value any

type (
	ComparisonOperator string
	LogicalOperator    string
	ValueType          string
	ExprType           string
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

	ConstExprType   ExprType = "const"
	VarExprType     ExprType = "var"
	CompareExprType ExprType = "compare"
	LogicalExprType ExprType = "logical"
	NotExprType     ExprType = "not"
)

type EvalResult struct {
	Type  ValueType
	Value any
}

func (e *EvalResult) GetFloat64() (float64, error) {
	if e.Type != Number {
		return 0, fmt.Errorf("type mismatch: %s", e.Type)
	}

	stringValue, ok := e.Value.(string)
	if !ok {
		return 0, fmt.Errorf("value is not a string: %v", e.Value)
	}

	floatValue, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse float64: %v", err)
	}
	return floatValue, nil
}

func (e *EvalResult) GetBool() (bool, error) {
	if e.Type != Boolean {
		return false, fmt.Errorf("type mismatch: %s", e.Type)
	}

	strValue, ok := e.Value.(string)
	if ok {
		if strValue == "true" {
			return true, nil
		} else if strValue == "false" {
			return false, nil
		}
		return false, fmt.Errorf("value is not a boolean: %v", e.Value)
	}

	boolValue, ok := e.Value.(bool)
	if !ok {
		return false, fmt.Errorf("value is not a boolean: %v", e.Value)
	}

	return boolValue, nil
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
		left, err := e.GetBool()
		if err != nil {
			return nil, err
		}
		right, err := other.GetBool()
		if err != nil {
			return nil, err
		}
		return &EvalResult{
			Type:  Boolean,
			Value: left && right,
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
		left, err := e.GetBool()
		if err != nil {
			return nil, err
		}
		right, err := other.GetBool()
		if err != nil {
			return nil, err
		}
		return &EvalResult{
			Type:  Boolean,
			Value: left || right,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", e.Type)
	}
}

func (e *EvalResult) Not() (*EvalResult, error) {
	if e.Type != Boolean {
		return nil, fmt.Errorf("type mismatch: %s", e.Type)
	}
	boolValue, err := e.GetBool()
	if err != nil {
		return nil, err
	}
	return &EvalResult{
		Type:  Boolean,
		Value: !boolValue,
	}, nil
}

func (e *EvalResult) Is(other *EvalResult) (*EvalResult, error) {
	if e.Type != other.Type {
		return nil, fmt.Errorf("type mismatch: %s vs %s", e.Type, other.Type)
	}
	switch e.Type {
	case Boolean:
		left, err := e.GetBool()
		if err != nil {
			return nil, err
		}
		right, err := other.GetBool()
		if err != nil {
			return nil, err
		}
		return &EvalResult{
			Type:  Boolean,
			Value: left == right,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", e.Type)
	}
}

type Expr interface {
	GetType() ExprType
	Eval(variables map[string]map[string]any) (*EvalResult, error)
}

type ConstExpr struct {
	Type      ExprType  `json:"type"`
	Value     Value     `json:"value"`
	ValueType ValueType `json:"valueType"`
}

func (c ConstExpr) GetType() ExprType { return c.Type }

func (c ConstExpr) Eval(variables map[string]map[string]any) (*EvalResult, error) {
	return &EvalResult{
		Type:  c.ValueType,
		Value: c.Value,
	}, nil
}

type VarExpr struct {
	Type     ExprType                    `json:"type"`
	Selector WorkflowNodeIOValueSelector `json:"selector"`
}

func (v VarExpr) GetType() ExprType { return v.Type }

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
	Type  ExprType           `json:"type"` // compare
	Op    ComparisonOperator `json:"op"`
	Left  Expr               `json:"left"`
	Right Expr               `json:"right"`
}

func (c CompareExpr) GetType() ExprType { return c.Type }

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
	Type  ExprType        `json:"type"` // logical
	Op    LogicalOperator `json:"op"`
	Left  Expr            `json:"left"`
	Right Expr            `json:"right"`
}

func (l LogicalExpr) GetType() ExprType { return l.Type }

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
	Type ExprType `json:"type"` // not
	Expr Expr     `json:"expr"`
}

func (n NotExpr) GetType() ExprType { return n.Type }

func (n NotExpr) Eval(variables map[string]map[string]any) (*EvalResult, error) {
	inner, err := n.Expr.Eval(variables)
	if err != nil {
		return nil, err
	}
	return inner.Not()
}

type rawExpr struct {
	Type ExprType `json:"type"`
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
	case ConstExprType:
		var e ConstExpr
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return e, nil
	case VarExprType:
		var e VarExpr
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return e, nil
	case CompareExprType:
		var e CompareExprRaw
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return e.ToCompareExpr()
	case LogicalExprType:
		var e LogicalExprRaw
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return e.ToLogicalExpr()
	case NotExprType:
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
	Type  ExprType           `json:"type"`
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
	Type  ExprType        `json:"type"`
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
	Type ExprType        `json:"type"`
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
