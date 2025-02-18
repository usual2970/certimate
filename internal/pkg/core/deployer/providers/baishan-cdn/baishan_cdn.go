package baishancdn

import (
	"context"
	"errors"
	"fmt"
	"time"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	bssdk "github.com/usual2970/certimate/internal/pkg/vendors/baishan-sdk"
)

type BaishanCDNDeployerConfig struct {
	// 白山云 API Token。
	ApiToken string `json:"apiToken"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type BaishanCDNDeployer struct {
	config    *BaishanCDNDeployerConfig
	logger    logger.Logger
	sdkClient *bssdk.Client
}

var _ deployer.Deployer = (*BaishanCDNDeployer)(nil)

func New(config *BaishanCDNDeployerConfig) (*BaishanCDNDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *BaishanCDNDeployerConfig, logger logger.Logger) (*BaishanCDNDeployer, error) {
	if config == nil {
		panic("config is nil")
	}

	if logger == nil {
		panic("logger is nil")
	}

	client, err := createSdkClient(config.ApiToken)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &BaishanCDNDeployer{
		logger:    logger,
		config:    config,
		sdkClient: client,
	}, nil
}

func (d *BaishanCDNDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 查询域名配置
	// REF: https://portal.baishancloud.com/track/document/api/1/1065
	getDomainConfigReq := &bssdk.GetDomainConfigRequest{
		Domains: d.config.Domain,
		Config:  "https",
	}
	getDomainConfigResp, err := d.sdkClient.GetDomainConfig(getDomainConfigReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'baishan.GetDomainConfig'")
	} else if len(getDomainConfigResp.Data) == 0 {
		return nil, errors.New("domain config not found")
	} else {
		d.logger.Logt("已查询到域名配置", getDomainConfigResp)
	}

	// 新增证书
	// REF: https://portal.baishancloud.com/track/document/downloadPdf/1441
	createCertificateReq := &bssdk.CreateCertificateRequest{
		Certificate: certPem,
		Key:         privkeyPem,
		Name:        fmt.Sprintf("certimate_%d", time.Now().UnixMilli()),
	}
	createCertificateResp, err := d.sdkClient.CreateCertificate(createCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'baishan.CreateCertificate'")
	} else {
		d.logger.Logt("已新增证书", createCertificateResp)
	}

	// 设置域名配置
	// REF: https://portal.baishancloud.com/track/document/api/1/1045
	setDomainConfigReq := &bssdk.SetDomainConfigRequest{
		Domains: d.config.Domain,
		Config: &bssdk.DomainConfig{
			Https: &bssdk.DomainConfigHttps{
				CertId:      createCertificateResp.Data.CertId,
				ForceHttps:  getDomainConfigResp.Data[0].Config.Https.ForceHttps,
				EnableHttp2: getDomainConfigResp.Data[0].Config.Https.EnableHttp2,
				EnableOcsp:  getDomainConfigResp.Data[0].Config.Https.EnableOcsp,
			},
		},
	}
	setDomainConfigResp, err := d.sdkClient.SetDomainConfig(setDomainConfigReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'baishan.SetDomainConfig'")
	} else {
		d.logger.Logt("已设置域名配置", setDomainConfigResp)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiToken string) (*bssdk.Client, error) {
	if apiToken == "" {
		return nil, errors.New("invalid baishan api token")
	}

	client := bssdk.NewClient(apiToken)
	return client, nil
}
