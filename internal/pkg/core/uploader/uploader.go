package uploader

import (
	"context"
	"log/slog"
)

// 表示定义证书上传器的抽象类型接口。
// 云服务商通常会提供 SSL 证书管理服务，可供用户集中管理证书。
// 注意与 `Deployer` 区分，“上传”通常为“部署”的前置操作。
type Uploader interface {
	WithLogger(logger *slog.Logger) Uploader

	// 上传证书。
	//
	// 入参：
	//   - ctx：上下文。
	//   - certPEM：证书 PEM 内容。
	//   - privkeyPEM：私钥 PEM 内容。
	//
	// 出参：
	//   - res：上传结果。
	//   - err: 错误。
	Upload(ctx context.Context, certPEM string, privkeyPEM string) (_res *UploadResult, _err error)
}

// 表示证书上传结果的数据结构，包含上传后的证书 ID、名称和其他数据。
type UploadResult struct {
	CertId       string         `json:"certId"`
	CertName     string         `json:"certName,omitzero"`
	ExtendedData map[string]any `json:"extendedData,omitempty"`
}
