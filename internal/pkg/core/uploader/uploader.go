﻿package uploader

import "context"

// 表示定义证书上传者的抽象类型接口。
// 云服务商通常会提供 SSL 证书管理服务，可供用户集中管理证书。
// 注意与 `Deployer` 区分，“上传”通常为“部署”的前置操作。
type Uploader interface {
	// 上传证书。
	//
	// 入参：
	//   - ctx：
	//   - certPem：证书 PEM 内容
	//   - privkeyPem：私钥 PEM 内容
	//
	// 出参：
	//   - res：
	//   - err：
	Upload(ctx context.Context, certPem string, privkeyPem string) (res *UploadResult, err error)
}

// 表示证书上传结果的数据结构，包含上传后的证书 ID、名称和其他数据。
type UploadResult struct {
	CertId   string                 `json:"certId"`
	CertName string                 `json:"certName"`
	CertData map[string]interface{} `json:"certData,omitempty"`
}