package jdcloudlive

import (
	"context"
	"fmt"
	"log/slog"

	jdcore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jdliveapi "github.com/jdcloud-api/jdcloud-sdk-go/services/live/apis"
	jdliveclient "github.com/jdcloud-api/jdcloud-sdk-go/services/live/client"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
)

type DeployerConfig struct {
	// 京东云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 京东云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 直播播放域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *jdliveclient.LiveClient
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	return &DeployerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	// 设置直播证书
	// REF: https://docs.jdcloud.com/cn/live-video/api/setlivedomaincertificate
	setLiveDomainCertificateReq := jdliveapi.NewSetLiveDomainCertificateRequest(d.config.Domain, "on")
	setLiveDomainCertificateReq.SetCert(certPEM)
	setLiveDomainCertificateReq.SetKey(privkeyPEM)
	setLiveDomainCertificateResp, err := d.sdkClient.SetLiveDomainCertificate(setLiveDomainCertificateReq)
	d.logger.Debug("sdk request 'live.SetLiveDomainCertificate'", slog.Any("request", setLiveDomainCertificateReq), slog.Any("response", setLiveDomainCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'live.SetLiveDomainCertificate': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*jdliveclient.LiveClient, error) {
	clientCredentials := jdcore.NewCredentials(accessKeyId, accessKeySecret)
	client := jdliveclient.NewLiveClient(clientCredentials)
	client.SetLogger(jdcore.NewDefaultLogger(jdcore.LogWarn))
	return client, nil
}
