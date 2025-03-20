package uclouducdn

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	xerrors "github.com/pkg/errors"
	uCdn "github.com/ucloud/ucloud-sdk-go/services/ucdn"
	usdk "github.com/ucloud/ucloud-sdk-go/ucloud"
	uAuth "github.com/ucloud/ucloud-sdk-go/ucloud/auth"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/ucloud-ussl"
)

type DeployerConfig struct {
	// 优刻得 API 私钥。
	PrivateKey string `json:"privateKey"`
	// 优刻得 API 公钥。
	PublicKey string `json:"publicKey"`
	// 优刻得项目 ID。
	ProjectId string `json:"projectId,omitempty"`
	// 加速域名 ID。
	DomainId string `json:"domainId"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *uCdn.UCDNClient
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.PrivateKey, config.PublicKey)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		PrivateKey: config.PrivateKey,
		PublicKey:  config.PublicKey,
		ProjectId:  config.ProjectId,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &DeployerProvider{
		config:      config,
		logger:      slog.Default(),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	d.sslUploader.WithLogger(logger)
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 USSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 获取加速域名配置
	// REF: https://docs.ucloud.cn/api/ucdn-api/get_ucdn_domain_config
	getUcdnDomainConfigReq := d.sdkClient.NewGetUcdnDomainConfigRequest()
	getUcdnDomainConfigReq.DomainId = []string{d.config.DomainId}
	if d.config.ProjectId != "" {
		getUcdnDomainConfigReq.ProjectId = usdk.String(d.config.ProjectId)
	}
	getUcdnDomainConfigResp, err := d.sdkClient.GetUcdnDomainConfig(getUcdnDomainConfigReq)
	d.logger.Debug("sdk request 'ucdn.GetUcdnDomainConfig'", slog.Any("request", getUcdnDomainConfigReq), slog.Any("response", getUcdnDomainConfigResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'ucdn.GetUcdnDomainConfig'")
	} else if len(getUcdnDomainConfigResp.DomainList) == 0 {
		return nil, errors.New("no domain found")
	}

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
	d.logger.Debug("sdk request 'ucdn.UpdateUcdnDomainHttpsConfigV2'", slog.Any("request", updateUcdnDomainHttpsConfigV2Req), slog.Any("response", updateUcdnDomainHttpsConfigV2Resp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'ucdn.UpdateUcdnDomainHttpsConfigV2'")
	}

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
