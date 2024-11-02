package x509

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	xerrors "github.com/pkg/errors"
)

// 从 PEM 编码的证书字符串解析并返回一个 x509.Certificate 对象。
//
// 入参:
//   - certPem: 证书 PEM 内容。
//
// 出参:
//   - cert: x509.Certificate 对象。
//   - err: 错误。
func ParseCertificateFromPEM(certPem string) (cert *x509.Certificate, err error) {
	pemData := []byte(certPem)

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to parse certificate")
	}

	return cert, nil
}

// 从 PEM 编码的私钥字符串解析并返回一个 ecdsa.PrivateKey 对象。
//
// 入参:
//   - privkeyPem: 私钥 PEM 内容。
//
// 出参:
//   - privkey: ecdsa.PrivateKey 对象。
//   - err: 错误。
func ParseECPrivateKeyFromPEM(privkeyPem string) (privkey *ecdsa.PrivateKey, err error) {
	pemData := []byte(privkeyPem)

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	privkey, err = x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to parse private key")
	}

	return privkey, nil
}

// 从 PEM 编码的私钥字符串解析并返回一个 rsa.PrivateKey 对象。
//
// 入参:
//   - privkeyPem: 私钥 PEM 内容。
//
// 出参:
//   - privkey: rsa.PrivateKey 对象。
//   - err: 错误。
func ParsePKCS1PrivateKeyFromPEM(privkeyPem string) (privkey *rsa.PrivateKey, err error) {
	pemData := []byte(privkeyPem)

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	privkey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to parse private key")
	}

	return privkey, nil
}
