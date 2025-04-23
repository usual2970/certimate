package jdcloudvod

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	jdcore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jdvodapi "github.com/jdcloud-api/jdcloud-sdk-go/services/vod/apis"
	jdvodclient "github.com/jdcloud-api/jdcloud-sdk-go/services/vod/client"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
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
	logger    *slog.Logger
	sdkClient *jdvodclient.VodClient
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
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	// 查询域名列表
	// REF: https://docs.jdcloud.com/cn/video-on-demand/api/listdomains
	var domainId int
	listDomainsPageNumber := 1
	listDomainsPageSize := 100
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		listDomainsReq := jdvodapi.NewListDomainsRequest()
		listDomainsReq.SetPageNumber(1)
		listDomainsReq.SetPageSize(100)
		listDomainsResp, err := d.sdkClient.ListDomains(listDomainsReq)
		d.logger.Debug("sdk request 'vod.ListDomains'", slog.Any("request", listDomainsReq), slog.Any("response", listDomainsResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'vod.ListDomains': %w", err)
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
		return nil, errors.New("domain not found")
	}

	// 查询域名 SSL 配置
	// REF: https://docs.jdcloud.com/cn/video-on-demand/api/gethttpssl
	getHttpSslReq := jdvodapi.NewGetHttpSslRequest(domainId)
	getHttpSslResp, err := d.sdkClient.GetHttpSsl(getHttpSslReq)
	d.logger.Debug("sdk request 'vod.GetHttpSsl'", slog.Any("request", getHttpSslReq), slog.Any("response", getHttpSslResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'vod.GetHttpSsl': %w", err)
	}

	// 设置域名 SSL 配置
	// REF: https://docs.jdcloud.com/cn/video-on-demand/api/sethttpssl
	setHttpSslReq := jdvodapi.NewSetHttpSslRequest(domainId)
	setHttpSslReq.SetTitle(fmt.Sprintf("certimate-%d", time.Now().UnixMilli()))
	setHttpSslReq.SetSslCert(certPEM)
	setHttpSslReq.SetSslKey(privkeyPEM)
	setHttpSslReq.SetSource("default")
	setHttpSslReq.SetJumpType(getHttpSslResp.Result.JumpType)
	setHttpSslReq.SetEnabled(true)
	setHttpSslResp, err := d.sdkClient.SetHttpSsl(setHttpSslReq)
	d.logger.Debug("sdk request 'vod.SetHttpSsl'", slog.Any("request", setHttpSslReq), slog.Any("response", setHttpSslResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'vod.SetHttpSsl': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*jdvodclient.VodClient, error) {
	clientCredentials := jdcore.NewCredentials(accessKeyId, accessKeySecret)
	client := jdvodclient.NewVodClient(clientCredentials)
	client.SetLogger(jdcore.NewDefaultLogger(jdcore.LogWarn))
	return client, nil
}
