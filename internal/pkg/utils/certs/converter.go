package certs

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"

	xerrors "github.com/pkg/errors"
)

// 将 ecdsa.PrivateKey 对象转换为 PEM 编码的字符串。
//
// 入参:
//   - privkey: ecdsa.PrivateKey 对象。
//
// 出参:
//   - privkeyPem: 私钥 PEM 内容。
//   - err: 错误。
func ConvertECPrivateKeyToPEM(privkey *ecdsa.PrivateKey) (privkeyPem string, err error) {
	data, err := x509.MarshalECPrivateKey(privkey)
	if err != nil {
		return "", xerrors.Wrap(err, "failed to marshal EC private key")
	}

	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: data,
	}

	return string(pem.EncodeToMemory(block)), nil
}
