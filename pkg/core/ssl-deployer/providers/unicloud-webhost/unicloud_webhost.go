package unicloudwebhost

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/certimate-go/certimate/pkg/core"
	unisdk "github.com/certimate-go/certimate/pkg/sdk3rd/dcloud/unicloud"
)

type SSLDeployerProviderConfig struct {
	// uniCloud 控制台账号。
	Username string `json:"username"`
	// uniCloud 控制台密码。
	Password string `json:"password"`
	// 服务空间提供商。
	// 可取值 "aliyun"、"tencent"。
	SpaceProvider string `json:"spaceProvider"`
	// 服务空间 ID。
	SpaceId string `json:"spaceId"`
	// 托管网站域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *unisdk.Client
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.Username, config.Password)
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
	if d.config.SpaceProvider == "" {
		return nil, errors.New("config `spaceProvider` is required")
	}
	if d.config.SpaceId == "" {
		return nil, errors.New("config `spaceId` is required")
	}
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 变更网站证书
	createDomainWithCertReq := &unisdk.CreateDomainWithCertRequest{
		Provider: d.config.SpaceProvider,
		SpaceId:  d.config.SpaceId,
		Domain:   d.config.Domain,
		Cert:     url.QueryEscape(certPEM),
		Key:      url.QueryEscape(privkeyPEM),
	}
	createDomainWithCertResp, err := d.sdkClient.CreateDomainWithCert(createDomainWithCertReq)
	d.logger.Debug("sdk request 'unicloud.host.CreateDomainWithCert'", slog.Any("request", createDomainWithCertReq), slog.Any("response", createDomainWithCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'unicloud.host.CreateDomainWithCert': %w", err)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(username, password string) (*unisdk.Client, error) {
	return unisdk.NewClient(username, password)
}
