package nodeprocessor

import (
	"context"
	"errors"
	"time"

	"github.com/usual2970/certimate/internal/domain"
)

type NodeProcessor interface {
	Run(ctx context.Context) error
	Log(ctx context.Context) *domain.WorkflowRunLog
	AddOutput(ctx context.Context, title, content string, err ...string)
}

type nodeLogger struct {
	log *domain.WorkflowRunLog
}

type certificateRepository interface {
	GetByWorkflowNodeId(ctx context.Context, workflowNodeId string) (*domain.Certificate, error)
}

type workflowOutputRepository interface {
	GetByNodeId(ctx context.Context, nodeId string) (*domain.WorkflowOutput, error)
	Save(ctx context.Context, output *domain.WorkflowOutput, certificate *domain.Certificate, cb func(id string) error) error
}

type settingsRepository interface {
	GetByName(ctx context.Context, name string) (*domain.Settings, error)
}

func NewNodeLogger(node *domain.WorkflowNode) *nodeLogger {
	return &nodeLogger{
		log: &domain.WorkflowRunLog{
			NodeId:   node.Id,
			NodeName: node.Name,
			Outputs:  make([]domain.WorkflowRunLogOutput, 0),
		},
	}
}

func (l *nodeLogger) Log(ctx context.Context) *domain.WorkflowRunLog {
	return l.log
}

func (l *nodeLogger) AddOutput(ctx context.Context, title, content string, err ...string) {
	output := domain.WorkflowRunLogOutput{
		Time:    time.Now().UTC().Format(time.RFC3339),
		Title:   title,
		Content: content,
	}
	if len(err) > 0 {
		output.Error = err[0]
		l.log.Error = err[0]
	}
	l.log.Outputs = append(l.log.Outputs, output)
}

func GetProcessor(node *domain.WorkflowNode) (NodeProcessor, error) {
	switch node.Type {
	case domain.WorkflowNodeTypeStart:
		return NewStartNode(node), nil
	case domain.WorkflowNodeTypeCondition:
		return NewConditionNode(node), nil
	case domain.WorkflowNodeTypeApply:
		return NewApplyNode(node), nil
	case domain.WorkflowNodeTypeUpload:
		return NewUploadNode(node), nil
	case domain.WorkflowNodeTypeDeploy:
		return NewDeployNode(node), nil
	case domain.WorkflowNodeTypeNotify:
		return NewNotifyNode(node), nil
	case domain.WorkflowNodeTypeExecuteSuccess:
		return NewExecuteSuccessNode(node), nil
	case domain.WorkflowNodeTypeExecuteFailure:
		return NewExecuteFailureNode(node), nil
	}
	return nil, errors.New("not implemented")
}

func getContextWorkflowId(ctx context.Context) string {
	return ctx.Value("workflow_id").(string)
}
