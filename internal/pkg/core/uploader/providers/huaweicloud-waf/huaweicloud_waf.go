package huaweicloudwaf

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	hciam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	hciammodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	hciamregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
	hcwaf "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/waf/v1"
	hcwafmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/waf/v1/model"
	hcwafregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/waf/v1/region"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
	typeutil "github.com/usual2970/certimate/internal/pkg/utils/type"
)

type UploaderConfig struct {
	// 华为云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 华为云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 华为云企业项目 ID。
	EnterpriseProjectId string `json:"enterpriseProjectId,omitempty"`
	// 华为云区域。
	Region string `json:"region"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *hcwaf.WafClient
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
		u.logger = slog.New(slog.DiscardHandler)
	} else {
		u.logger = logger
	}
	return u
}

func (u *UploaderProvider) Upload(ctx context.Context, certPEM string, privkeyPEM string) (*uploader.UploadResult, error) {
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 遍历查询已有证书，避免重复上传
	// REF: https://support.huaweicloud.com/api-waf/ListCertificates.html
	// REF: https://support.huaweicloud.com/api-waf/ShowCertificate.html
	listCertificatesPage := int32(1)
	listCertificatesPageSize := int32(100)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		listCertificatesReq := &hcwafmodel.ListCertificatesRequest{
			EnterpriseProjectId: typeutil.ToPtrOrZeroNil(u.config.EnterpriseProjectId),
			Page:                typeutil.ToPtr(listCertificatesPage),
			Pagesize:            typeutil.ToPtr(listCertificatesPageSize),
		}
		listCertificatesResp, err := u.sdkClient.ListCertificates(listCertificatesReq)
		u.logger.Debug("sdk request 'waf.ShowCertificate'", slog.Any("request", listCertificatesReq), slog.Any("response", listCertificatesResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'waf.ListCertificates': %w", err)
		}

		if listCertificatesResp.Items != nil {
			for _, certItem := range *listCertificatesResp.Items {
				showCertificateReq := &hcwafmodel.ShowCertificateRequest{
					EnterpriseProjectId: typeutil.ToPtrOrZeroNil(u.config.EnterpriseProjectId),
					CertificateId:       certItem.Id,
				}
				showCertificateResp, err := u.sdkClient.ShowCertificate(showCertificateReq)
				u.logger.Debug("sdk request 'waf.ShowCertificate'", slog.Any("request", showCertificateReq), slog.Any("response", showCertificateResp))
				if err != nil {
					return nil, fmt.Errorf("failed to execute sdk request 'waf.ShowCertificate': %w", err)
				}

				var isSameCert bool
				if *showCertificateResp.Content == certPEM {
					isSameCert = true
				} else {
					oldCertX509, err := certutil.ParseCertificateFromPEM(*showCertificateResp.Content)
					if err != nil {
						continue
					}

					isSameCert = certutil.EqualCertificate(certX509, oldCertX509)
				}

				// 如果已存在相同证书，直接返回
				if isSameCert {
					u.logger.Info("ssl certificate already exists")
					return &uploader.UploadResult{
						CertId:   certItem.Id,
						CertName: certItem.Name,
					}, nil
				}
			}
		}

		if listCertificatesResp.Items == nil || len(*listCertificatesResp.Items) < int(listCertificatesPageSize) {
			break
		} else {
			listCertificatesPage++
		}
	}

	// 生成新证书名（需符合华为云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 创建证书
	// REF: https://support.huaweicloud.com/api-waf/CreateCertificate.html
	createCertificateReq := &hcwafmodel.CreateCertificateRequest{
		EnterpriseProjectId: typeutil.ToPtrOrZeroNil(u.config.EnterpriseProjectId),
		Body: &hcwafmodel.CreateCertificateRequestBody{
			Name:    certName,
			Content: certPEM,
			Key:     privkeyPEM,
		},
	}
	createCertificateResp, err := u.sdkClient.CreateCertificate(createCertificateReq)
	u.logger.Debug("sdk request 'waf.CreateCertificate'", slog.Any("request", createCertificateReq), slog.Any("response", createCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'waf.CreateCertificate': %w", err)
	}

	certId = *createCertificateResp.Id
	certName = *createCertificateResp.Name
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*hcwaf.WafClient, error) {
	projectId, err := getSdkProjectId(accessKeyId, secretAccessKey, region)
	if err != nil {
		return nil, err
	}

	auth, err := basic.NewCredentialsBuilder().
		WithAk(accessKeyId).
		WithSk(secretAccessKey).
		WithProjectId(projectId).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	hcRegion, err := hcwafregion.SafeValueOf(region)
	if err != nil {
		return nil, err
	}

	hcClient, err := hcwaf.WafClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := hcwaf.NewWafClient(hcClient)
	return client, nil
}

func getSdkProjectId(accessKeyId, secretAccessKey, region string) (string, error) {
	auth, err := global.NewCredentialsBuilder().
		WithAk(accessKeyId).
		WithSk(secretAccessKey).
		SafeBuild()
	if err != nil {
		return "", err
	}

	hcRegion, err := hciamregion.SafeValueOf(region)
	if err != nil {
		return "", err
	}

	hcClient, err := hciam.IamClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return "", err
	}

	client := hciam.NewIamClient(hcClient)

	request := &hciammodel.KeystoneListProjectsRequest{
		Name: &region,
	}
	response, err := client.KeystoneListProjects(request)
	if err != nil {
		return "", err
	} else if response.Projects == nil || len(*response.Projects) == 0 {
		return "", errors.New("no project found")
	}

	return (*response.Projects)[0].Id, nil
}
