package nodeprocessor

import (
	"context"
	"testing"

	"github.com/usual2970/certimate/internal/domain"
)

func Test_inspectWebsiteCertificateNode_inspect(t *testing.T) {
	type args struct {
		ctx        context.Context
		nodeConfig domain.WorkflowNodeConfigForInspect
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				ctx: context.Background(),
				nodeConfig: domain.WorkflowNodeConfigForInspect{
					Domain: "baidu.com",
					Port:   "443",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewInspectNode(&domain.WorkflowNode{})
			if err := n.inspect(tt.args.ctx, tt.args.nodeConfig); (err != nil) != tt.wantErr {
				t.Errorf("inspectWebsiteCertificateNode.inspect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
