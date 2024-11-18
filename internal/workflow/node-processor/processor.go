package nodeprocessor

import (
	"context"
	"errors"

	"github.com/usual2970/certimate/internal/domain"
)

type RunLog struct {
	NodeName string         `json:"node_name"`
	Err      string         `json:"err"`
	Outputs  []RunLogOutput `json:"outputs"`
}

type RunLogOutput struct {
	Time    string `json:"time"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Error   string `json:"error"`
}

type NodeProcessor interface {
	Run(ctx context.Context) error
	Log(ctx context.Context) *RunLog
	AddOutput(ctx context.Context, time, title, content string, err ...string)
}

type Logger struct {
	log *RunLog
}

func NewLogger(node *domain.WorkflowNode) *Logger {
	return &Logger{
		log: &RunLog{
			NodeName: node.Name,
			Outputs:  make([]RunLogOutput, 0),
		},
	}
}

func (l *Logger) Log(ctx context.Context) *RunLog {
	return l.log
}

func (l *Logger) AddOutput(ctx context.Context, time, title, content string, err ...string) {
	output := RunLogOutput{
		Time:    time,
		Title:   title,
		Content: content,
	}
	if len(err) > 0 {
		output.Error = err[0]
		l.log.Err = err[0]
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
	}
	return nil, errors.New("not implemented")
}
