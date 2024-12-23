package nodeprocessor

import (
	"context"
	"errors"
	"time"

	"github.com/usual2970/certimate/internal/domain"
)

type NodeProcessor interface {
	Run(ctx context.Context) error
	Log(ctx context.Context) *domain.RunLog
	AddOutput(ctx context.Context, title, content string, err ...string)
}

type Logger struct {
	log *domain.RunLog
}

func NewLogger(node *domain.WorkflowNode) *Logger {
	return &Logger{
		log: &domain.RunLog{
			NodeName: node.Name,
			Outputs:  make([]domain.RunLogOutput, 0),
		},
	}
}

func (l *Logger) Log(ctx context.Context) *domain.RunLog {
	return l.log
}

func (l *Logger) AddOutput(ctx context.Context, title, content string, err ...string) {
	output := domain.RunLogOutput{
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
	case domain.WorkflowNodeTypeDeploy:
		return NewDeployNode(node), nil
	case domain.WorkflowNodeTypeNotify:
		return NewNotifyNode(node), nil
	}
	return nil, errors.New("not implemented")
}
