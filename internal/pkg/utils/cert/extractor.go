package certutil

import (
	"encoding/pem"
	"errors"
)

// 从 PEM 编码的证书字符串解析并提取服务器证书和中间证书。
//
// 入参:
//   - certPEM: 证书 PEM 内容。
//
// 出参:
//   - serverCertPEM: 服务器证书的 PEM 内容。
//   - intermediaCertPEM: 中间证书的 PEM 内容。
//   - err: 错误。
func ExtractCertificatesFromPEM(certPEM string) (_serverCertPEM string, _intermediaCertPEM string, _err error) {
	pemBlocks := make([]*pem.Block, 0)
	pemData := []byte(certPEM)
	for {
		block, rest := pem.Decode(pemData)
		if block == nil || block.Type != "CERTIFICATE" {
			break
		}

		pemBlocks = append(pemBlocks, block)
		pemData = rest
	}

	serverCertPEM := ""
	intermediaCertPEM := ""

	if len(pemBlocks) == 0 {
		return "", "", errors.New("failed to decode PEM block")
	}

	if len(pemBlocks) > 0 {
		serverCertPEM = string(pem.EncodeToMemory(pemBlocks[0]))
	}

	if len(pemBlocks) > 1 {
		for i := 1; i < len(pemBlocks); i++ {
			intermediaCertPEM += string(pem.EncodeToMemory(pemBlocks[i]))
		}
	}

	return serverCertPEM, intermediaCertPEM, nil
}
