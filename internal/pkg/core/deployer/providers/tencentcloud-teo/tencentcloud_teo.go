package tencentcloudeteo

import (
	"context"
	"errors"

	xerrors "github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	tcTeo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	providerSsl "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/tencentcloud-ssl"
)

type TencentCloudTEODeployerConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 站点 ID。
	ZoneId string `json:"zoneId"`
	// 加速域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type TencentCloudTEODeployer struct {
	config      *TencentCloudTEODeployerConfig
	logger      deployer.Logger
	sdkClients  *wSdkClients
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*TencentCloudTEODeployer)(nil)

type wSdkClients struct {
	ssl *tcSsl.Client
	teo *tcTeo.Client
}

func New(config *TencentCloudTEODeployerConfig) (*TencentCloudTEODeployer, error) {
	return NewWithLogger(config, deployer.NewNilLogger())
}

func NewWithLogger(config *TencentCloudTEODeployerConfig, logger deployer.Logger) (*TencentCloudTEODeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	clients, err := createSdkClients(config.SecretId, config.SecretKey)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk clients")
	}

	uploader, err := providerSsl.New(&providerSsl.TencentCloudSSLUploaderConfig{
		SecretId:  config.SecretId,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &TencentCloudTEODeployer{
		logger:      logger,
		config:      config,
		sdkClients:  clients,
		sslUploader: uploader,
	}, nil
}

func (d *TencentCloudTEODeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	if d.config.ZoneId == "" {
		return nil, xerrors.New("config `zoneId` is required")
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Appendt("certificate file uploaded", upres)

	// 配置域名证书
	// REF: https://cloud.tencent.com/document/product/1552/80764
	modifyHostsCertificateReq := tcTeo.NewModifyHostsCertificateRequest()
	modifyHostsCertificateReq.ZoneId = common.StringPtr(d.config.ZoneId)
	modifyHostsCertificateReq.Mode = common.StringPtr("sslcert")
	modifyHostsCertificateReq.Hosts = common.StringPtrs([]string{d.config.Domain})
	modifyHostsCertificateReq.ServerCertInfo = []*tcTeo.ServerCertInfo{{CertId: common.StringPtr(upres.CertId)}}
	modifyHostsCertificateResp, err := d.sdkClients.teo.ModifyHostsCertificate(modifyHostsCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'teo.ModifyHostsCertificate'")
	}

	d.logger.Appendt("已配置域名证书", modifyHostsCertificateResp.Response)

	return &deployer.DeployResult{}, nil
}

func createSdkClients(secretId, secretKey string) (*wSdkClients, error) {
	credential := common.NewCredential(secretId, secretKey)

	sslClient, err := tcSsl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	teoClient, err := tcTeo.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &wSdkClients{
		ssl: sslClient,
		teo: teoClient,
	}, nil
}
