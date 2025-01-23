package uclouducdn

import (
	"context"
	"errors"
	"strconv"

	xerrors "github.com/pkg/errors"
	uCdn "github.com/ucloud/ucloud-sdk-go/services/ucdn"
	usdk "github.com/ucloud/ucloud-sdk-go/ucloud"
	uAuth "github.com/ucloud/ucloud-sdk-go/ucloud/auth"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploaderp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/ucloud-ussl"
)

type UCloudUCDNDeployerConfig struct {
	// 优刻得 API 私钥。
	PrivateKey string `json:"privateKey"`
	// 优刻得 API 公钥。
	PublicKey string `json:"publicKey"`
	// 优刻得项目 ID。
	ProjectId string `json:"projectId,omitempty"`
	// 加速域名 ID。
	DomainId string `json:"domainId"`
}

type UCloudUCDNDeployer struct {
	config      *UCloudUCDNDeployerConfig
	logger      logger.Logger
	sdkClient   *uCdn.UCDNClient
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*UCloudUCDNDeployer)(nil)

func New(config *UCloudUCDNDeployerConfig) (*UCloudUCDNDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *UCloudUCDNDeployerConfig, logger logger.Logger) (*UCloudUCDNDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.PrivateKey, config.PublicKey)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploaderp.New(&uploaderp.UCloudUSSLUploaderConfig{
		PrivateKey: config.PrivateKey,
		PublicKey:  config.PublicKey,
		ProjectId:  config.ProjectId,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &UCloudUCDNDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *UCloudUCDNDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 USSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 获取加速域名配置
	// REF: https://docs.ucloud.cn/api/ucdn-api/get_ucdn_domain_config
	getUcdnDomainConfigReq := d.sdkClient.NewGetUcdnDomainConfigRequest()
	getUcdnDomainConfigReq.DomainId = []string{d.config.DomainId}
	if d.config.ProjectId != "" {
		getUcdnDomainConfigReq.ProjectId = usdk.String(d.config.ProjectId)
	}
	getUcdnDomainConfigResp, err := d.sdkClient.GetUcdnDomainConfig(getUcdnDomainConfigReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'ucdn.GetUcdnDomainConfig'")
	} else if len(getUcdnDomainConfigResp.DomainList) == 0 {
		return nil, errors.New("no domain found")
	}

	d.logger.Logt("已查询到加速域名配置", getUcdnDomainConfigResp)

	// 更新 HTTPS 加速配置
	// REF: https://docs.ucloud.cn/api/ucdn-api/update_ucdn_domain_https_config_v2
	certId, _ := strconv.Atoi(upres.CertId)
	updateUcdnDomainHttpsConfigV2Req := d.sdkClient.NewUpdateUcdnDomainHttpsConfigV2Request()
	updateUcdnDomainHttpsConfigV2Req.DomainId = usdk.String(d.config.DomainId)
	updateUcdnDomainHttpsConfigV2Req.HttpsStatusCn = usdk.String(getUcdnDomainConfigResp.DomainList[0].HttpsStatusCn)
	updateUcdnDomainHttpsConfigV2Req.HttpsStatusAbroad = usdk.String(getUcdnDomainConfigResp.DomainList[0].HttpsStatusAbroad)
	updateUcdnDomainHttpsConfigV2Req.HttpsStatusAbroad = usdk.String(getUcdnDomainConfigResp.DomainList[0].HttpsStatusAbroad)
	updateUcdnDomainHttpsConfigV2Req.CertId = usdk.Int(certId)
	updateUcdnDomainHttpsConfigV2Req.CertName = usdk.String(upres.CertName)
	updateUcdnDomainHttpsConfigV2Req.CertType = usdk.String("ussl")
	if d.config.ProjectId != "" {
		updateUcdnDomainHttpsConfigV2Req.ProjectId = usdk.String(d.config.ProjectId)
	}
	updateUcdnDomainHttpsConfigV2Resp, err := d.sdkClient.UpdateUcdnDomainHttpsConfigV2(updateUcdnDomainHttpsConfigV2Req)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'ucdn.UpdateUcdnDomainHttpsConfigV2'")
	}

	d.logger.Logt("已更新 HTTPS 加速配置", updateUcdnDomainHttpsConfigV2Resp)

	return &deployer.DeployResult{}, nil
}

func createSdkClient(privateKey, publicKey string) (*uCdn.UCDNClient, error) {
	cfg := usdk.NewConfig()

	credential := uAuth.NewCredential()
	credential.PrivateKey = privateKey
	credential.PublicKey = publicKey

	client := uCdn.NewClient(&cfg, &credential)
	return client, nil
}
