package huaweicloudelb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	hcelb "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3"
	hcelbmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
	hcelbregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/region"
	hciam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	hciammodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	hciamregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"

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
	sdkClient *hcelb.ElbClient
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
	// REF: https://support.huaweicloud.com/api-elb/ListCertificates.html
	listCertificatesLimit := int32(2000)
	var listCertificatesMarker *string = nil
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		listCertificatesReq := &hcelbmodel.ListCertificatesRequest{
			Limit:  typeutil.ToPtr(listCertificatesLimit),
			Marker: listCertificatesMarker,
			Type:   &[]string{"server"},
		}
		listCertificatesResp, err := u.sdkClient.ListCertificates(listCertificatesReq)
		u.logger.Debug("sdk request 'elb.ListCertificates'", slog.Any("request", listCertificatesReq), slog.Any("response", listCertificatesResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'elb.ListCertificates': %w", err)
		}

		if listCertificatesResp.Certificates != nil {
			for _, certDetail := range *listCertificatesResp.Certificates {
				var isSameCert bool
				if certDetail.Certificate == certPEM {
					isSameCert = true
				} else {
					oldCertX509, err := certutil.ParseCertificateFromPEM(certDetail.Certificate)
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
			listCertificatesMarker = listCertificatesResp.PageInfo.NextMarker
		}
	}

	// 获取项目 ID
	// REF: https://support.huaweicloud.com/api-iam/iam_06_0001.html
	projectId, err := getSdkProjectId(u.config.AccessKeyId, u.config.SecretAccessKey, u.config.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to get SDK project id: %w", err)
	}

	// 生成新证书名（需符合华为云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 创建新证书
	// REF: https://support.huaweicloud.com/api-elb/CreateCertificate.html
	createCertificateReq := &hcelbmodel.CreateCertificateRequest{
		Body: &hcelbmodel.CreateCertificateRequestBody{
			Certificate: &hcelbmodel.CreateCertificateOption{
				ProjectId:   typeutil.ToPtr(projectId),
				Name:        typeutil.ToPtr(certName),
				Certificate: typeutil.ToPtr(certPEM),
				PrivateKey:  typeutil.ToPtr(privkeyPEM),
			},
		},
	}
	createCertificateResp, err := u.sdkClient.CreateCertificate(createCertificateReq)
	u.logger.Debug("sdk request 'elb.CreateCertificate'", slog.Any("request", createCertificateReq), slog.Any("response", createCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'elb.CreateCertificate': %w", err)
	}

	certId = createCertificateResp.Certificate.Id
	certName = createCertificateResp.Certificate.Name
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*hcelb.ElbClient, error) {
	if region == "" {
		region = "cn-north-4" // ELB 服务默认区域：华北四北京
	}

	auth, err := basic.NewCredentialsBuilder().
		WithAk(accessKeyId).
		WithSk(secretAccessKey).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	hcRegion, err := hcelbregion.SafeValueOf(region)
	if err != nil {
		return nil, err
	}

	hcClient, err := hcelb.ElbClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := hcelb.NewElbClient(hcClient)
	return client, nil
}

func getSdkProjectId(accessKeyId, secretAccessKey, region string) (string, error) {
	if region == "" {
		region = "cn-north-4" // IAM 服务默认区域：华北四北京
	}

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
