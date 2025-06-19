package cert

import (
	"bytes"
	"encoding/pem"
	"errors"
	"time"

	"github.com/pavlo-v-chernykh/keystore-go/v4"
	"software.sslmate.com/src/go-pkcs12"
)

// 将 PEM 编码的证书字符串转换为 PFX 格式。
//
// 入参:
//   - certPEM: 证书 PEM 内容。
//   - privkeyPEM: 私钥 PEM 内容。
//   - pfxPassword: PFX 导出密码。
//
// 出参:
//   - data: PFX 格式的证书数据。
//   - err: 错误。
func TransformCertificateFromPEMToPFX(certPEM string, privkeyPEM string, pfxPassword string) ([]byte, error) {
	cert, err := ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	privkey, err := ParsePrivateKeyFromPEM(privkeyPEM)
	if err != nil {
		return nil, err
	}

	pfxData, err := pkcs12.LegacyRC2.Encode(privkey, cert, nil, pfxPassword)
	if err != nil {
		return nil, err
	}

	return pfxData, nil
}

// 将 PEM 编码的证书字符串转换为 JKS 格式。
//
// 入参:
//   - certPEM: 证书 PEM 内容。
//   - privkeyPEM: 私钥 PEM 内容。
//   - jksAlias: JKS 别名。
//   - jksKeypass: JKS 密钥密码。
//   - jksStorepass: JKS 存储密码。
//
// 出参:
//   - data: JKS 格式的证书数据。
//   - err: 错误。
func TransformCertificateFromPEMToJKS(certPEM string, privkeyPEM string, jksAlias string, jksKeypass string, jksStorepass string) ([]byte, error) {
	certBlock, _ := pem.Decode([]byte(certPEM))
	if certBlock == nil {
		return nil, errors.New("failed to decode certificate PEM")
	}

	privkeyBlock, _ := pem.Decode([]byte(privkeyPEM))
	if privkeyBlock == nil {
		return nil, errors.New("failed to decode private key PEM")
	}

	ks := keystore.New()
	entry := keystore.PrivateKeyEntry{
		CreationTime: time.Now(),
		PrivateKey:   privkeyBlock.Bytes,
		CertificateChain: []keystore.Certificate{
			{
				Type:    "X509",
				Content: certBlock.Bytes,
			},
		},
	}

	if err := ks.SetPrivateKeyEntry(jksAlias, entry, []byte(jksKeypass)); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := ks.Store(&buf, []byte(jksStorepass)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
