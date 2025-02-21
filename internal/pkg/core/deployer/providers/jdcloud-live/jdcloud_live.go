package jdcloudlive

import (
	"context"

	jdCore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jdLiveApi "github.com/jdcloud-api/jdcloud-sdk-go/services/live/apis"
	jdLiveClient "github.com/jdcloud-api/jdcloud-sdk-go/services/live/client"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
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
	logger    logger.Logger
	sdkClient *jdLiveClient.LiveClient
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &DeployerProvider{
		config:    config,
		logger:    logger.NewNilLogger(),
		sdkClient: client,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger logger.Logger) *DeployerProvider {
	d.logger = logger
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 设置直播证书
	// REF: https://docs.jdcloud.com/cn/live-video/api/setlivedomaincertificate
	setLiveDomainCertificateReq := jdLiveApi.NewSetLiveDomainCertificateRequest(d.config.Domain, "on")
	setLiveDomainCertificateReq.SetCert(certPem)
	setLiveDomainCertificateReq.SetKey(privkeyPem)
	setLiveDomainCertificateResp, err := d.sdkClient.SetLiveDomainCertificate(setLiveDomainCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'live.SetLiveDomainCertificate'")
	} else {
		d.logger.Logt("已设置直播证书", setLiveDomainCertificateResp)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*jdLiveClient.LiveClient, error) {
	clientCredentials := jdCore.NewCredentials(accessKeyId, accessKeySecret)
	client := jdLiveClient.NewLiveClient(clientCredentials)
	client.SetLogger(jdCore.NewDefaultLogger(jdCore.LogWarn))
	return client, nil
}
