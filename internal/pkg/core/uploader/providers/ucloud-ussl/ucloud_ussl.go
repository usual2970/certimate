package ucloudussl

import (
	"context"
	"crypto/md5"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ucloud/ucloud-sdk-go/ucloud"
	ucloudauth "github.com/ucloud/ucloud-sdk-go/ucloud/auth"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	usslsdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/ucloud/ussl"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

type UploaderConfig struct {
	// 优刻得 API 私钥。
	PrivateKey string `json:"privateKey"`
	// 优刻得 API 公钥。
	PublicKey string `json:"publicKey"`
	// 优刻得项目 ID。
	ProjectId string `json:"projectId,omitempty"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *usslsdk.USSLClient
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.PrivateKey, config.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	return &UploaderProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (u *UploaderProvider) WithLogger(logger *slog.Logger) uploader.Uploader {
	if logger == nil {
		u.logger = slog.Default()
	} else {
		u.logger = logger
	}
	return u
}

func (u *UploaderProvider) Upload(ctx context.Context, certPEM string, privkeyPEM string) (res *uploader.UploadResult, err error) {
	// 生成新证书名（需符合优刻得命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 生成优刻得所需的证书参数
	certPEMBase64 := base64.StdEncoding.EncodeToString([]byte(certPEM))
	privkeyPEMBase64 := base64.StdEncoding.EncodeToString([]byte(privkeyPEM))
	certMd5 := md5.Sum([]byte(certPEMBase64 + privkeyPEMBase64))
	certMd5Hex := hex.EncodeToString(certMd5[:])

	// 上传托管证书
	// REF: https://docs.ucloud.cn/api/usslcertificate-api/upload_normal_certificate
	uploadNormalCertificateReq := u.sdkClient.NewUploadNormalCertificateRequest()
	uploadNormalCertificateReq.CertificateName = ucloud.String(certName)
	uploadNormalCertificateReq.SslPublicKey = ucloud.String(certPEMBase64)
	uploadNormalCertificateReq.SslPrivateKey = ucloud.String(privkeyPEMBase64)
	uploadNormalCertificateReq.SslMD5 = ucloud.String(certMd5Hex)
	if u.config.ProjectId != "" {
		uploadNormalCertificateReq.ProjectId = ucloud.String(u.config.ProjectId)
	}
	uploadNormalCertificateResp, err := u.sdkClient.UploadNormalCertificate(uploadNormalCertificateReq)
	u.logger.Debug("sdk request 'ussl.UploadNormalCertificate'", slog.Any("request", uploadNormalCertificateReq), slog.Any("response", uploadNormalCertificateResp))
	if err != nil {
		if uploadNormalCertificateResp != nil && uploadNormalCertificateResp.GetRetCode() == 80035 {
			if res, err := u.getCertIfExists(ctx, certPEM); err != nil {
				return nil, err
			} else if res == nil {
				return nil, errors.New("ucloud ssl: no certificate found")
			} else {
				u.logger.Info("ssl certificate already exists")
				return res, nil
			}
		}

		return nil, fmt.Errorf("failed to execute sdk request 'ussl.UploadNormalCertificate': %w", err)
	}

	certId = fmt.Sprintf("%d", uploadNormalCertificateResp.CertificateID)
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
		ExtendedData: map[string]any{
			"resourceId": uploadNormalCertificateResp.LongResourceID,
		},
	}, nil
}

func (u *UploaderProvider) getCertIfExists(ctx context.Context, certPEM string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 遍历获取用户证书列表，避免重复上传
	// REF: https://docs.ucloud.cn/api/usslcertificate-api/get_certificate_list
	// REF: https://docs.ucloud.cn/api/usslcertificate-api/download_certificate
	getCertificateListPage := int(1)
	getCertificateListLimit := int(1000)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		getCertificateListReq := u.sdkClient.NewGetCertificateListRequest()
		getCertificateListReq.Mode = ucloud.String("trust")
		getCertificateListReq.Domain = ucloud.String(certX509.Subject.CommonName)
		getCertificateListReq.Sort = ucloud.String("2")
		getCertificateListReq.Page = ucloud.Int(getCertificateListPage)
		getCertificateListReq.PageSize = ucloud.Int(getCertificateListLimit)
		if u.config.ProjectId != "" {
			getCertificateListReq.ProjectId = ucloud.String(u.config.ProjectId)
		}
		getCertificateListResp, err := u.sdkClient.GetCertificateList(getCertificateListReq)
		u.logger.Debug("sdk request 'ussl.GetCertificateList'", slog.Any("request", getCertificateListReq), slog.Any("response", getCertificateListResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'ussl.GetCertificateList': %w", err)
		}

		if getCertificateListResp.CertificateList != nil {
			for _, certInfo := range getCertificateListResp.CertificateList {
				// 优刻得未提供可唯一标识证书的字段，只能通过多个字段尝试对比来判断是否为同一证书
				// 先分别对比证书的多域名、品牌、有效期，再对比签名算法

				if len(certX509.DNSNames) == 0 || certInfo.Domains != strings.Join(certX509.DNSNames, ",") {
					continue
				}

				if len(certX509.Issuer.Organization) == 0 || certInfo.Brand != certX509.Issuer.Organization[0] {
					continue
				}

				if int64(certInfo.NotBefore) != certX509.NotBefore.UnixMilli() || int64(certInfo.NotAfter) != certX509.NotAfter.UnixMilli() {
					continue
				}

				getCertificateDetailInfoReq := u.sdkClient.NewGetCertificateDetailInfoRequest()
				getCertificateDetailInfoReq.CertificateID = ucloud.Int(certInfo.CertificateID)
				if u.config.ProjectId != "" {
					getCertificateDetailInfoReq.ProjectId = ucloud.String(u.config.ProjectId)
				}
				getCertificateDetailInfoResp, err := u.sdkClient.GetCertificateDetailInfo(getCertificateDetailInfoReq)
				if err != nil {
					return nil, fmt.Errorf("failed to execute sdk request 'ussl.GetCertificateDetailInfo': %w", err)
				}

				switch certX509.SignatureAlgorithm {
				case x509.SHA256WithRSA:
					if !strings.EqualFold(getCertificateDetailInfoResp.CertificateInfo.Algorithm, "SHA256-RSA") {
						continue
					}
				case x509.SHA384WithRSA:
					if !strings.EqualFold(getCertificateDetailInfoResp.CertificateInfo.Algorithm, "SHA384-RSA") {
						continue
					}
				case x509.SHA512WithRSA:
					if !strings.EqualFold(getCertificateDetailInfoResp.CertificateInfo.Algorithm, "SHA512-RSA") {
						continue
					}
				case x509.SHA256WithRSAPSS:
					if !strings.EqualFold(getCertificateDetailInfoResp.CertificateInfo.Algorithm, "SHA256-RSAPSS") {
						continue
					}
				case x509.SHA384WithRSAPSS:
					if !strings.EqualFold(getCertificateDetailInfoResp.CertificateInfo.Algorithm, "SHA384-RSAPSS") {
						continue
					}
				case x509.SHA512WithRSAPSS:
					if !strings.EqualFold(getCertificateDetailInfoResp.CertificateInfo.Algorithm, "SHA512-RSAPSS") {
						continue
					}
				case x509.ECDSAWithSHA256:
					if !strings.EqualFold(getCertificateDetailInfoResp.CertificateInfo.Algorithm, "ECDSA-SHA256") {
						continue
					}
				case x509.ECDSAWithSHA384:
					if !strings.EqualFold(getCertificateDetailInfoResp.CertificateInfo.Algorithm, "ECDSA-SHA384") {
						continue
					}
				case x509.ECDSAWithSHA512:
					if !strings.EqualFold(getCertificateDetailInfoResp.CertificateInfo.Algorithm, "ECDSA-SHA512") {
						continue
					}
				default:
					// 未知签名算法，跳过
					continue
				}

				return &uploader.UploadResult{
					CertId:   fmt.Sprintf("%d", certInfo.CertificateID),
					CertName: certInfo.Name,
					ExtendedData: map[string]any{
						"resourceId": certInfo.CertificateSN,
					},
				}, nil
			}
		}

		if getCertificateListResp.CertificateList == nil || len(getCertificateListResp.CertificateList) < int(getCertificateListLimit) {
			break
		} else {
			getCertificateListPage++
		}
	}

	return nil, nil
}

func createSdkClient(privateKey, publicKey string) (*usslsdk.USSLClient, error) {
	cfg := ucloud.NewConfig()

	credential := ucloudauth.NewCredential()
	credential.PrivateKey = privateKey
	credential.PublicKey = publicKey

	client := usslsdk.NewClient(&cfg, &credential)
	return client, nil
}
