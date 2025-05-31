package expr

import (
	"testing"
)

func TestLogicalEval(t *testing.T) {
	// 测试逻辑表达式 and
	logicalExpr := LogicalExpr{
		Left: ConstantExpr{
			Type:      "const",
			Value:     "true",
			ValueType: "boolean",
		},
		Operator: And,
		Right: ConstantExpr{
			Type:      "const",
			Value:     "true",
			ValueType: "boolean",
		},
	}
	result, err := logicalExpr.Eval(nil)
	if err != nil {
		t.Errorf("failed to evaluate logical expression: %v", err)
	}
	if result.Value != true {
		t.Errorf("expected true, got %v", result)
	}

	// 测试逻辑表达式 or
	orExpr := LogicalExpr{
		Left: ConstantExpr{
			Type:      "const",
			Value:     "true",
			ValueType: "boolean",
		},
		Operator: Or,
		Right: ConstantExpr{
			Type:      "const",
			Value:     "true",
			ValueType: "boolean",
		},
	}
	result, err = orExpr.Eval(nil)
	if err != nil {
		t.Errorf("failed to evaluate logical expression: %v", err)
	}
	if result.Value != true {
		t.Errorf("expected true, got %v", result)
	}
}

func TestUnmarshalExpr(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Expr
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				data: []byte(`{"left":{"left":{"selector":{"id":"ODnYSOXB6HQP2_vz6JcZE","name":"certificate.validity","type":"boolean"},"type":"var"},"operator":"is","right":{"type":"const","value":true,"valueType":"boolean"},"type":"comparison"},"operator":"and","right":{"left":{"selector":{"id":"ODnYSOXB6HQP2_vz6JcZE","name":"certificate.daysLeft","type":"number"},"type":"var"},"operator":"eq","right":{"type":"const","value":2,"valueType":"number"},"type":"comparison"},"type":"logical"}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalExpr(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalExpr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("UnmarshalExpr() got = nil, want %v", tt.want)
				return
			}
		})
	}
}

func TestExpr_Eval(t *testing.T) {
	type args struct {
		variables map[string]map[string]any
		data      []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *EvalResult
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				variables: map[string]map[string]any{
					"ODnYSOXB6HQP2_vz6JcZE": {
						"certificate.validity": true,
						"certificate.daysLeft": 2,
					},
				},
				data: []byte(`{"left":{"left":{"selector":{"id":"ODnYSOXB6HQP2_vz6JcZE","name":"certificate.validity","type":"boolean"},"type":"var"},"operator":"is","right":{"type":"const","value":true,"valueType":"boolean"},"type":"comparison"},"operator":"and","right":{"left":{"selector":{"id":"ODnYSOXB6HQP2_vz6JcZE","name":"certificate.daysLeft","type":"number"},"type":"var"},"operator":"eq","right":{"type":"const","value":2,"valueType":"number"},"type":"comparison"},"type":"logical"}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := UnmarshalExpr(tt.args.data)
			if err != nil {
				t.Errorf("UnmarshalExpr() error = %v", err)
				return
			}
			got, err := c.Eval(tt.args.variables)
			t.Log("got:", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConstExpr.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Value != true {
				t.Errorf("ConstExpr.Eval() got = %v, want %v", got.Value, true)
			}
		})
	}
}
