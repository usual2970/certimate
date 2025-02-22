package jdcloudvod

import (
	"context"
	"fmt"
	"strconv"
	"time"

	jdCore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jdVodApi "github.com/jdcloud-api/jdcloud-sdk-go/services/vod/apis"
	jdVodClient "github.com/jdcloud-api/jdcloud-sdk-go/services/vod/client"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
)

type DeployerConfig struct {
	// 京东云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 京东云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 点播加速域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    logger.Logger
	sdkClient *jdVodClient.VodClient
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
	// 查询域名列表
	// REF: https://docs.jdcloud.com/cn/video-on-demand/api/listdomains
	var domainId int
	listDomainsPageNumber := 1
	listDomainsPageSize := 100
	for {
		listDomainsReq := jdVodApi.NewListDomainsRequest()
		listDomainsReq.SetPageNumber(1)
		listDomainsReq.SetPageSize(100)
		listDomainsResp, err := d.sdkClient.ListDomains(listDomainsReq)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'vod.ListDomains'")
		}

		for _, domain := range listDomainsResp.Result.Content {
			if domain.Name == d.config.Domain {
				domainId, _ = strconv.Atoi(domain.Id)
				break
			}
		}

		if len(listDomainsResp.Result.Content) < listDomainsPageSize {
			break
		} else {
			listDomainsPageNumber++
		}
	}
	if domainId == 0 {
		return nil, xerrors.New("domain not found")
	}

	// 查询域名 SSL 配置
	// REF: https://docs.jdcloud.com/cn/video-on-demand/api/gethttpssl
	getHttpSslReq := jdVodApi.NewGetHttpSslRequest(domainId)
	getHttpSslResp, err := d.sdkClient.GetHttpSsl(getHttpSslReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'vod.GetHttpSsl'")
	} else {
		d.logger.Logt("已查询到域名 SSL 配置", getHttpSslResp)
	}

	// 设置域名 SSL 配置
	// REF: https://docs.jdcloud.com/cn/video-on-demand/api/sethttpssl
	setHttpSslReq := jdVodApi.NewSetHttpSslRequest(domainId)
	setHttpSslReq.SetTitle(fmt.Sprintf("certimate-%d", time.Now().UnixMilli()))
	setHttpSslReq.SetSslCert(certPem)
	setHttpSslReq.SetSslKey(privkeyPem)
	setHttpSslReq.SetSource("default")
	setHttpSslReq.SetJumpType(getHttpSslResp.Result.JumpType)
	setHttpSslReq.SetEnabled(true)
	setHttpSslResp, err := d.sdkClient.SetHttpSsl(setHttpSslReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'vod.SetHttpSsl'")
	} else {
		d.logger.Logt("已设置域名 SSL 配置", setHttpSslResp)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*jdVodClient.VodClient, error) {
	clientCredentials := jdCore.NewCredentials(accessKeyId, accessKeySecret)
	client := jdVodClient.NewVodClient(clientCredentials)
	client.SetLogger(jdCore.NewDefaultLogger(jdCore.LogWarn))
	return client, nil
}
