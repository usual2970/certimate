package huaweicloudscm

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	hcscm "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/scm/v3"
	hcscmmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/scm/v3/model"
	hcscmregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/scm/v3/region"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
	typeutil "github.com/usual2970/certimate/internal/pkg/utils/type"
)

type UploaderConfig struct {
	// 华为云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 华为云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 华为云区域。
	Region string `json:"region"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *hcscm.ScmClient
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.SecretAccessKey, config.Region)
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
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 遍历查询已有证书，避免重复上传
	// REF: https://support.huaweicloud.com/api-ccm/ListCertificates.html
	// REF: https://support.huaweicloud.com/api-ccm/ExportCertificate_0.html
	listCertificatesLimit := int32(50)
	listCertificatesOffset := int32(0)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		listCertificatesReq := &hcscmmodel.ListCertificatesRequest{
			Limit:   typeutil.ToPtr(listCertificatesLimit),
			Offset:  typeutil.ToPtr(listCertificatesOffset),
			SortDir: typeutil.ToPtr("DESC"),
			SortKey: typeutil.ToPtr("certExpiredTime"),
		}
		listCertificatesResp, err := u.sdkClient.ListCertificates(listCertificatesReq)
		u.logger.Debug("sdk request 'scm.ListCertificates'", slog.Any("request", listCertificatesReq), slog.Any("response", listCertificatesResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'scm.ListCertificates': %w", err)
		}

		if listCertificatesResp.Certificates != nil {
			for _, certDetail := range *listCertificatesResp.Certificates {
				exportCertificateReq := &hcscmmodel.ExportCertificateRequest{
					CertificateId: certDetail.Id,
				}
				exportCertificateResp, err := u.sdkClient.ExportCertificate(exportCertificateReq)
				u.logger.Debug("sdk request 'scm.ExportCertificate'", slog.Any("request", exportCertificateReq), slog.Any("response", exportCertificateResp))
				if err != nil {
					if exportCertificateResp != nil && exportCertificateResp.HttpStatusCode == 404 {
						continue
					}
					return nil, fmt.Errorf("failed to execute sdk request 'scm.ExportCertificate': %w", err)
				}

				var isSameCert bool
				if *exportCertificateResp.Certificate == certPEM {
					isSameCert = true
				} else {
					oldCertX509, err := certutil.ParseCertificateFromPEM(*exportCertificateResp.Certificate)
					if err != nil {
						continue
					}

					isSameCert = certutil.EqualCertificate(certX509, oldCertX509)
				}

				// 如果已存在相同证书，直接返回
				if isSameCert {
					u.logger.Info("ssl certificate already exists")
					return &uploader.UploadResult{
						CertId:   certDetail.Id,
						CertName: certDetail.Name,
					}, nil
				}
			}
		}

		if listCertificatesResp.Certificates == nil || len(*listCertificatesResp.Certificates) < int(listCertificatesLimit) {
			break
		} else {
			listCertificatesOffset += listCertificatesLimit
		}
	}

	// 生成新证书名（需符合华为云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://support.huaweicloud.com/api-ccm/ImportCertificate.html
	importCertificateReq := &hcscmmodel.ImportCertificateRequest{
		Body: &hcscmmodel.ImportCertificateRequestBody{
			Name:        certName,
			Certificate: certPEM,
			PrivateKey:  privkeyPEM,
		},
	}
	importCertificateResp, err := u.sdkClient.ImportCertificate(importCertificateReq)
	u.logger.Debug("sdk request 'scm.ImportCertificate'", slog.Any("request", importCertificateReq), slog.Any("response", importCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'scm.ImportCertificate': %w", err)
	}

	certId = *importCertificateResp.CertificateId
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*hcscm.ScmClient, error) {
	if region == "" {
		region = "cn-north-4" // SCM 服务默认区域：华北四北京
	}

	auth, err := basic.NewCredentialsBuilder().
		WithAk(accessKeyId).
		WithSk(secretAccessKey).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	hcRegion, err := hcscmregion.SafeValueOf(region)
	if err != nil {
		return nil, err
	}

	hcClient, err := hcscm.ScmClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := hcscm.NewScmClient(hcClient)
	return client, nil
}
