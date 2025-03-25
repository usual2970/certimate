package baishancdn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	bssdk "github.com/usual2970/certimate/internal/pkg/vendors/baishan-sdk"
)

type DeployerConfig struct {
	// 白山云 API Token。
	ApiToken string `json:"apiToken"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *bssdk.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ApiToken)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &DeployerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 查询域名配置
	// REF: https://portal.baishancloud.com/track/document/api/1/1065
	getDomainConfigReq := &bssdk.GetDomainConfigRequest{
		Domains: d.config.Domain,
		Config:  []string{"https"},
	}
	getDomainConfigResp, err := d.sdkClient.GetDomainConfig(getDomainConfigReq)
	d.logger.Debug("sdk request 'baishan.GetDomainConfig'", slog.Any("request", getDomainConfigReq), slog.Any("response", getDomainConfigResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'baishan.GetDomainConfig'")
	} else if len(getDomainConfigResp.Data) == 0 {
		return nil, errors.New("domain config not found")
	}

	// 新增证书
	// REF: https://portal.baishancloud.com/track/document/downloadPdf/1441
	certificateId := ""
	createCertificateReq := &bssdk.CreateCertificateRequest{
		Certificate: certPem,
		Key:         privkeyPem,
		Name:        fmt.Sprintf("certimate_%d", time.Now().UnixMilli()),
	}
	createCertificateResp, err := d.sdkClient.CreateCertificate(createCertificateReq)
	d.logger.Debug("sdk request 'baishan.CreateCertificate'", slog.Any("request", createCertificateReq), slog.Any("response", createCertificateResp))
	if err != nil {
		if createCertificateResp != nil {
			if createCertificateResp.GetCode() == 400699 && strings.Contains(createCertificateResp.GetMessage(), "this certificate is exists") {
				// 证书已存在，忽略新增证书接口错误
				re := regexp.MustCompile(`\d+`)
				certificateId = re.FindString(createCertificateResp.GetMessage())
			}
		}

		if certificateId == "" {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'baishan.CreateCertificate'")
		}
	} else {
		certificateId = createCertificateResp.Data.CertId.String()
	}

	// 设置域名配置
	// REF: https://portal.baishancloud.com/track/document/api/1/1045
	setDomainConfigReq := &bssdk.SetDomainConfigRequest{
		Domains: d.config.Domain,
		Config: &bssdk.DomainConfig{
			Https: &bssdk.DomainConfigHttps{
				CertId:      json.Number(certificateId),
				ForceHttps:  getDomainConfigResp.Data[0].Config.Https.ForceHttps,
				EnableHttp2: getDomainConfigResp.Data[0].Config.Https.EnableHttp2,
				EnableOcsp:  getDomainConfigResp.Data[0].Config.Https.EnableOcsp,
			},
		},
	}
	setDomainConfigResp, err := d.sdkClient.SetDomainConfig(setDomainConfigReq)
	d.logger.Debug("sdk request 'baishan.SetDomainConfig'", slog.Any("request", setDomainConfigReq), slog.Any("response", setDomainConfigResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'baishan.SetDomainConfig'")
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
