package certutil

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	xerrors "github.com/pkg/errors"
)

// 将 x509.Certificate 对象转换为 PEM 编码的字符串。
//
// 入参:
//   - cert: x509.Certificate 对象。
//
// 出参:
//   - certPem: 证书 PEM 内容。
//   - err: 错误。
func ConvertCertificateToPEM(cert *x509.Certificate) (certPem string, err error) {
	if cert == nil {
		return "", errors.New("`cert` is nil")
	}

	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}

	return string(pem.EncodeToMemory(block)), nil
}

// 将 ecdsa.PrivateKey 对象转换为 PEM 编码的字符串。
//
// 入参:
//   - privkey: ecdsa.PrivateKey 对象。
//
// 出参:
//   - privkeyPem: 私钥 PEM 内容。
//   - err: 错误。
func ConvertECPrivateKeyToPEM(privkey *ecdsa.PrivateKey) (privkeyPem string, err error) {
	if privkey == nil {
		return "", errors.New("`privkey` is nil")
	}

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
