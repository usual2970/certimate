package x509

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
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
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert, nil
}

// 比较两个 x509.Certificate 对象，判断它们是否是同一张证书。
// 注意，这不是精确比较，而只是基于证书序列号和数字签名的快速判断，但对于权威 CA 签发的证书来说不会存在误判。
//
// 入参:
//   - a: 待比较的第一个 x509.Certificate 对象。
//   - b: 待比较的第二个 x509.Certificate 对象。
//
// 出参:
//   - 是否相同。
func EqualCertificate(a, b *x509.Certificate) bool {
	return string(a.Signature) == string(b.Signature) &&
		a.SignatureAlgorithm == b.SignatureAlgorithm &&
		a.SerialNumber.String() == b.SerialNumber.String() &&
		a.Issuer.SerialNumber == b.Issuer.SerialNumber &&
		a.Subject.SerialNumber == b.Subject.SerialNumber
}

// 将 ECDSA 私钥转换为 PEM 格式的字符串。
func PrivateKeyToPEM(privateKey *ecdsa.PrivateKey) (string, error) {
	data, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal EC private key: %w", err)
	}

	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: data,
	}

	return string(pem.EncodeToMemory(block)), nil
}

// 从 PEM 编码的私钥字符串解析并返回一个 ECDSA 私钥对象。
func ParsePrivateKeyFromPEM(privateKeyPem string) (*ecdsa.PrivateKey, error) {
	pemData := []byte(privateKeyPem)

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}
