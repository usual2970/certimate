package deployer

import (
	"context"
	"log/slog"
)

// 表示定义证书部署器的抽象类型接口。
// 注意与 `Uploader` 区分，“部署”通常为“上传”的后置操作。
type Deployer interface {
	WithLogger(logger *slog.Logger) Deployer

	// 部署证书。
	//
	// 入参：
	//   - ctx：上下文。
	//   - certPem：证书 PEM 内容。
	//   - privkeyPem：私钥 PEM 内容。
	//
	// 出参：
	//   - res：部署结果。
	//   - err: 错误。
	Deploy(ctx context.Context, certPem string, privkeyPem string) (res *DeployResult, err error)
}

// 表示证书部署结果的数据结构。
type DeployResult struct {
	ExtendedData map[string]any `json:"extendedData,omitempty"`
}
