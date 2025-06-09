package nodeprocessor

import (
	"context"
	"sync"
)

// 定义上下文键类型，避免键冲突
type workflowContextKey string

const (
	nodeOutputsKey workflowContextKey = "node_outputs"
)

// 带互斥锁的节点输出容器
type nodeOutputsContainer struct {
	sync.RWMutex
	outputs map[string]map[string]any
}

// 创建新的并发安全的节点输出容器
func newNodeOutputsContainer() *nodeOutputsContainer {
	return &nodeOutputsContainer{
		outputs: make(map[string]map[string]any),
	}
}

// 获取节点输出容器
func getNodeOutputsContainer(ctx context.Context) *nodeOutputsContainer {
	value := ctx.Value(nodeOutputsKey)
	if value == nil {
		return nil
	}
	return value.(*nodeOutputsContainer)
}

// 添加节点输出到上下文
func AddNodeOutput(ctx context.Context, nodeId string, output map[string]any) context.Context {
	container := getNodeOutputsContainer(ctx)
	if container == nil {
		container = newNodeOutputsContainer()
	}

	container.Lock()
	defer container.Unlock()

	// 创建输出的深拷贝
	// TODO: 暂时使用浅拷贝，等后续值类型扩充后修改
	outputCopy := make(map[string]any, len(output))
	for k, v := range output {
		outputCopy[k] = v
	}

	container.outputs[nodeId] = outputCopy
	return context.WithValue(ctx, nodeOutputsKey, container)
}

// 从上下文获取节点输出
func GetNodeOutput(ctx context.Context, nodeId string) map[string]any {
	container := getNodeOutputsContainer(ctx)
	if container == nil {
		container = newNodeOutputsContainer()
	}

	container.RLock()
	defer container.RUnlock()

	output, exists := container.outputs[nodeId]
	if !exists {
		return nil
	}

	outputCopy := make(map[string]any, len(output))
	for k, v := range output {
		outputCopy[k] = v
	}

	return outputCopy
}

// 获取所有节点输出
func GetAllNodeOutputs(ctx context.Context) map[string]map[string]any {
	container := getNodeOutputsContainer(ctx)
	if container == nil {
		container = newNodeOutputsContainer()
	}

	container.RLock()
	defer container.RUnlock()

	// 创建所有输出的深拷贝
	// TODO: 暂时使用浅拷贝，等后续值类型扩充后修改
	allOutputs := make(map[string]map[string]any, len(container.outputs))
	for nodeId, output := range container.outputs {
		nodeCopy := make(map[string]any, len(output))
		for k, v := range output {
			nodeCopy[k] = v
		}
		allOutputs[nodeId] = nodeCopy
	}

	return allOutputs
}
