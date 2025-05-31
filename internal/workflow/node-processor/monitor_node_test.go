package nodeprocessor_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/usual2970/certimate/internal/domain"
	nodeprocessor "github.com/usual2970/certimate/internal/workflow/node-processor"
)

func Test_MonitorNode(t *testing.T) {
	t.Run("Monitor", func(t *testing.T) {
		node := nodeprocessor.NewMonitorNode(&domain.WorkflowNode{
			Id:   "test",
			Type: domain.WorkflowNodeTypeMonitor,
			Name: "test",
			Config: map[string]any{
				"host": "baidu.com",
				"port": 443,
			},
		})
		node.SetLogger(slog.Default())
		if err := node.Process(context.Background()); err != nil {
			t.Errorf("err: %+v", err)
		}
	})
}
