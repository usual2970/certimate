package certutil

import (
	"encoding/pem"
	"errors"
)

// 从 PEM 编码的证书字符串解析并提取服务器证书和中间证书。
//
// 入参:
//   - certPem: 证书 PEM 内容。
//
// 出参:
//   - serverCertPem: 服务器证书的 PEM 内容。
//   - interCertPem: 中间证书的 PEM 内容。
//   - err: 错误。
func ExtractCertificatesFromPEM(certPem string) (serverCertPem string, interCertPem string, err error) {
	pemBlocks := make([]*pem.Block, 0)
	pemData := []byte(certPem)
	for {
		block, rest := pem.Decode(pemData)
		if block == nil || block.Type != "CERTIFICATE" {
			break
		}

		pemBlocks = append(pemBlocks, block)
		pemData = rest
	}

	serverCertPem = ""
	interCertPem = ""

	if len(pemBlocks) == 0 {
		return "", "", errors.New("failed to decode PEM block")
	}

	if len(pemBlocks) > 0 {
		serverCertPem = string(pem.EncodeToMemory(pemBlocks[0]))
	}

	if len(pemBlocks) > 1 {
		for i := 1; i < len(pemBlocks); i++ {
			interCertPem += string(pem.EncodeToMemory(pemBlocks[i]))
		}
	}

	return serverCertPem, interCertPem, nil
}
