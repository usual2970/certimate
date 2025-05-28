package unicloudwebhost

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	unisdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/dcloud/unicloud"
)

type DeployerConfig struct {
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

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *unisdk.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.Username, config.Password)
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

	return &deployer.DeployResult{}, nil
}

func createSdkClient(username, password string) (*unisdk.Client, error) {
	if username == "" {
		return nil, errors.New("invalid unicloud username")
	}

	if password == "" {
		return nil, errors.New("invalid unicloud password")
	}

	client := unisdk.NewClient(username, password)
	return client, nil
}
