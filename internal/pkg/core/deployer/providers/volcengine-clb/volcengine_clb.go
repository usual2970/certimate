package volcengineclb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	xerrors "github.com/pkg/errors"
	veClb "github.com/volcengine/volcengine-go-sdk/service/clb"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"
	veSession "github.com/volcengine/volcengine-go-sdk/volcengine/session"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/volcengine-certcenter"
)

type DeployerConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 火山引擎地域。
	Region string `json:"region"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 负载均衡监听器 ID。
	// 部署资源类型为 [RESOURCE_TYPE_LISTENER] 时必填。
	ListenerId string `json:"listenerId,omitempty"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *veClb.CLB
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
		Region:          config.Region,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
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
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	d.sslUploader.WithLogger(logger)
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到证书中心
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case RESOURCE_TYPE_LISTENER:
		if err := d.deployToListener(ctx, upres.CertId); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type: %s", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *DeployerProvider) deployToListener(ctx context.Context, cloudCertId string) error {
	if d.config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}

	// 修改监听器
	// REF: https://www.volcengine.com/docs/6406/71775
	modifyListenerAttributesReq := &veClb.ModifyListenerAttributesInput{
		ListenerId:              ve.String(d.config.ListenerId),
		CertificateSource:       ve.String("cert_center"),
		CertCenterCertificateId: ve.String(cloudCertId),
	}
	modifyListenerAttributesResp, err := d.sdkClient.ModifyListenerAttributes(modifyListenerAttributesReq)
	d.logger.Debug("sdk request 'clb.ModifyListenerAttributes'", slog.Any("request", modifyListenerAttributesReq), slog.Any("response", modifyListenerAttributesResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'clb.ModifyListenerAttributes'")
	}

	return nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*veClb.CLB, error) {
	config := ve.NewConfig().WithRegion(region).WithAkSk(accessKeyId, accessKeySecret)

	session, err := veSession.NewSession(config)
	if err != nil {
		return nil, err
	}

	client := veClb.New(session)
	return client, nil
}
