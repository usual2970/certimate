package jdcloudlive

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/pkg/core"
	jdcore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jdliveapi "github.com/jdcloud-api/jdcloud-sdk-go/services/live/apis"
	jdliveclient "github.com/jdcloud-api/jdcloud-sdk-go/services/live/client"
)

type SSLDeployerProviderConfig struct {
	// 京东云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 京东云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 直播播放域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *jdliveclient.LiveClient
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	return &SSLDeployerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	if d.config.Domain == "" {
		return nil, fmt.Errorf("config `domain` is required")
	}

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

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(accessKeyId, accessKeySecret string) (*jdliveclient.LiveClient, error) {
	clientCredentials := jdcore.NewCredentials(accessKeyId, accessKeySecret)
	client := jdliveclient.NewLiveClient(clientCredentials)
	client.SetLogger(jdcore.NewDefaultLogger(jdcore.LogWarn))
	return client, nil
}
