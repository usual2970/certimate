package certutil

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/go-acme/lego/v4/certcrypto"
)

// 从 PEM 编码的证书字符串解析并返回一个 x509.Certificate 对象。
// PEM 内容可能是包含多张证书的证书链，但只返回第一个证书（即服务器证书）。
//
// 入参:
//   - certPEM: 证书 PEM 内容。
//
// 出参:
//   - cert: x509.Certificate 对象。
//   - err: 错误。
func ParseCertificateFromPEM(certPEM string) (cert *x509.Certificate, err error) {
	pemData := []byte(certPEM)

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert, nil
}

// 从 PEM 编码的私钥字符串解析并返回一个 crypto.PrivateKey 对象。
//
// 入参:
//   - privkeyPEM: 私钥 PEM 内容。
//
// 出参:
//   - privkey: crypto.PrivateKey 对象，可能是 rsa.PrivateKey、ecdsa.PrivateKey 或 ed25519.PrivateKey。
//   - err: 错误。
func ParsePrivateKeyFromPEM(privkeyPEM string) (privkey crypto.PrivateKey, err error) {
	pemData := []byte(privkeyPEM)
	return certcrypto.ParsePEMPrivateKey(pemData)
}

// 从 PEM 编码的私钥字符串解析并返回一个 ecdsa.PrivateKey 对象。
//
// 入参:
//   - privkeyPEM: 私钥 PEM 内容。
//
// 出参:
//   - privkey: ecdsa.PrivateKey 对象。
//   - err: 错误。
func ParseECPrivateKeyFromPEM(privkeyPEM string) (privkey *ecdsa.PrivateKey, err error) {
	pemData := []byte(privkeyPEM)

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	privkey, err = x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privkey, nil
}

// 从 PEM 编码的私钥字符串解析并返回一个 rsa.PrivateKey 对象。
//
// 入参:
//   - privkeyPEM: 私钥 PEM 内容。
//
// 出参:
//   - privkey: rsa.PrivateKey 对象。
//   - err: 错误。
func ParsePKCS1PrivateKeyFromPEM(privkeyPEM string) (privkey *rsa.PrivateKey, err error) {
	pemData := []byte(privkeyPEM)

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	privkey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privkey, nil
}
