package aliyuncasdeploy

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	alicas "github.com/alibabacloud-go/cas-20200407/v3/client"
	aliopen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/aliyun-cas"
)

type DeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云资源组 ID。
	ResourceGroupId string `json:"resourceGroupId,omitempty"`
	// 阿里云地域。
	Region string `json:"region"`
	// 阿里云云产品资源 ID 数组。
	ResourceIds []string `json:"resourceIds"`
	// 阿里云云联系人 ID 数组。
	// 零值时使用账号下第一个联系人。
	ContactIds []string `json:"contactIds"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *alicas.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	uploader, err := createSslUploader(config.AccessKeyId, config.AccessKeySecret, config.ResourceGroupId, config.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create ssl uploader: %w", err)
	}

	return &DeployerProvider{
		config:      config,
		logger:      slog.Default(),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
	d.sslUploader.WithLogger(logger)
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	if len(d.config.ResourceIds) == 0 {
		return nil, errors.New("config `resourceIds` is required")
	}

	// 上传证书到 CAS
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	contactIds := d.config.ContactIds
	if len(contactIds) == 0 {
		// 获取联系人列表
		// REF: https://help.aliyun.com/zh/ssl-certificate/developer-reference/api-cas-2020-04-07-listcontact
		listContactReq := &alicas.ListContactRequest{
			ShowSize:    tea.Int32(1),
			CurrentPage: tea.Int32(1),
		}
		listContactResp, err := d.sdkClient.ListContact(listContactReq)
		d.logger.Debug("sdk request 'cas.ListContact'", slog.Any("request", listContactReq), slog.Any("response", listContactResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cas.ListContact': %w", err)
		}

		if len(listContactResp.Body.ContactList) > 0 {
			contactIds = []string{fmt.Sprintf("%d", listContactResp.Body.ContactList[0].ContactId)}
		}
	}

	// 创建部署任务
	// REF: https://help.aliyun.com/zh/ssl-certificate/developer-reference/api-cas-2020-04-07-createdeploymentjob
	createDeploymentJobReq := &alicas.CreateDeploymentJobRequest{
		Name:        tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		JobType:     tea.String("user"),
		CertIds:     tea.String(upres.CertId),
		ResourceIds: tea.String(strings.Join(d.config.ResourceIds, ",")),
		ContactIds:  tea.String(strings.Join(contactIds, ",")),
	}
	createDeploymentJobResp, err := d.sdkClient.CreateDeploymentJob(createDeploymentJobReq)
	d.logger.Debug("sdk request 'cas.CreateDeploymentJob'", slog.Any("request", createDeploymentJobReq), slog.Any("response", createDeploymentJobResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cas.CreateDeploymentJob': %w", err)
	}

	// 循环获取部署任务详情，等待任务状态变更
	// REF: https://help.aliyun.com/zh/ssl-certificate/developer-reference/api-cas-2020-04-07-describedeploymentjob
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		describeDeploymentJobReq := &alicas.DescribeDeploymentJobRequest{
			JobId: createDeploymentJobResp.Body.JobId,
		}
		describeDeploymentJobResp, err := d.sdkClient.DescribeDeploymentJob(describeDeploymentJobReq)
		d.logger.Debug("sdk request 'cas.DescribeDeploymentJob'", slog.Any("request", describeDeploymentJobReq), slog.Any("response", describeDeploymentJobResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cas.DescribeDeploymentJob': %w", err)
		}

		if describeDeploymentJobResp.Body.Status == nil || *describeDeploymentJobResp.Body.Status == "editing" {
			return nil, errors.New("unexpected deployment job status")
		}

		if *describeDeploymentJobResp.Body.Status == "success" || *describeDeploymentJobResp.Body.Status == "error" {
			break
		}

		d.logger.Info("waiting for deployment job completion ...")
		time.Sleep(time.Second * 5)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*alicas.Client, error) {
	// 接入点一览 https://api.aliyun.com/product/cas
	var endpoint string
	switch region {
	case "", "cn-hangzhou":
		endpoint = "cas.aliyuncs.com"
	default:
		endpoint = fmt.Sprintf("cas.%s.aliyuncs.com", region)
	}

	config := &aliopen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}

	client, err := alicas.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createSslUploader(accessKeyId, accessKeySecret, resourceGroupId, region string) (uploader.Uploader, error) {
	casRegion := region
	if casRegion != "" {
		// 阿里云 CAS 服务接入点是独立于其他服务的
		// 国内版固定接入点：华东一杭州
		// 国际版固定接入点：亚太东南一新加坡
		if !strings.HasPrefix(casRegion, "cn-") {
			casRegion = "ap-southeast-1"
		} else {
			casRegion = "cn-hangzhou"
		}
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		ResourceGroupId: resourceGroupId,
		Region:          casRegion,
	})
	return uploader, err
}
