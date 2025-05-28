package nodeprocessor

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/usual2970/certimate/internal/domain"
)

type NodeProcessor interface {
	GetLogger() *slog.Logger
	SetLogger(*slog.Logger)

	Process(ctx context.Context) error

	GetOutputs() map[string]any
}

type nodeProcessor struct {
	logger *slog.Logger
}

func (n *nodeProcessor) GetLogger() *slog.Logger {
	return n.logger
}

func (n *nodeProcessor) SetLogger(logger *slog.Logger) {
	if logger == nil {
		panic("logger is nil")
	}

	n.logger = logger
}

type nodeOutputer struct {
	outputs map[string]any
}

func newNodeOutputer() *nodeOutputer {
	return &nodeOutputer{
		outputs: make(map[string]any),
	}
}

func (n *nodeOutputer) GetOutputs() map[string]any {
	return n.outputs
}

type certificateRepository interface {
	GetByWorkflowNodeId(ctx context.Context, workflowNodeId string) (*domain.Certificate, error)
	GetByWorkflowRunId(ctx context.Context, workflowRunId string) (*domain.Certificate, error)
	Save(ctx context.Context, certificate *domain.Certificate) (*domain.Certificate, error)
}

type workflowOutputRepository interface {
	GetByNodeId(ctx context.Context, workflowNodeId string) (*domain.WorkflowOutput, error)
	Save(ctx context.Context, workflowOutput *domain.WorkflowOutput) (*domain.WorkflowOutput, error)
	SaveWithCertificate(ctx context.Context, workflowOutput *domain.WorkflowOutput, certificate *domain.Certificate) (*domain.WorkflowOutput, error)
}

type settingsRepository interface {
	GetByName(ctx context.Context, name string) (*domain.Settings, error)
}

func newNodeProcessor(node *domain.WorkflowNode) *nodeProcessor {
	return &nodeProcessor{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func GetProcessor(node *domain.WorkflowNode) (NodeProcessor, error) {
	switch node.Type {
	case domain.WorkflowNodeTypeStart:
		return NewStartNode(node), nil
	case domain.WorkflowNodeTypeApply:
		return NewApplyNode(node), nil
	case domain.WorkflowNodeTypeUpload:
		return NewUploadNode(node), nil
	case domain.WorkflowNodeTypeMonitor:
		return NewMonitorNode(node), nil
	case domain.WorkflowNodeTypeDeploy:
		return NewDeployNode(node), nil
	case domain.WorkflowNodeTypeNotify:
		return NewNotifyNode(node), nil
	case domain.WorkflowNodeTypeCondition:
		return NewConditionNode(node), nil
	case domain.WorkflowNodeTypeExecuteSuccess:
		return NewExecuteSuccessNode(node), nil
	case domain.WorkflowNodeTypeExecuteFailure:
		return NewExecuteFailureNode(node), nil
	}

	return nil, fmt.Errorf("unsupported node type: %s", string(node.Type))
}

func getContextWorkflowId(ctx context.Context) string {
	return ctx.Value("workflow_id").(string)
}

func getContextWorkflowRunId(ctx context.Context) string {
	return ctx.Value("workflow_run_id").(string)
}
