package gcorecdn

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	gprovider "github.com/G-Core/gcorelabscdn-go/gcore/provider"
	gresources "github.com/G-Core/gcorelabscdn-go/resources"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/gcore-cdn"
	gcoresdk "github.com/usual2970/certimate/internal/pkg/vendors/gcore-sdk/common"
)

type DeployerConfig struct {
	// Gcore API Token。
	ApiToken string `json:"apiToken"`
	// CDN 资源 ID。
	ResourceId int64 `json:"resourceId"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *gresources.Service
	sslUploader uploader.Uploader
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

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		ApiToken: config.ApiToken,
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
	if d.config.ResourceId == 0 {
		return nil, errors.New("config `resourceId` is required")
	}

	// 上传证书到 CDN
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 获取 CDN 资源详情
	// REF: https://api.gcore.com/docs/cdn#tag/CDN-resources/paths/~1cdn~1resources~1%7Bresource_id%7D/get
	getResourceResp, err := d.sdkClient.Get(context.TODO(), d.config.ResourceId)
	d.logger.Debug("sdk request 'resources.Get'", slog.Any("resourceId", d.config.ResourceId), slog.Any("response", getResourceResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'resources.Get'")
	}

	// 更新 CDN 资源详情
	// REF: https://api.gcore.com/docs/cdn#tag/CDN-resources/operation/change_cdn_resource
	updateResourceCertId, _ := strconv.ParseInt(upres.CertId, 10, 64)
	updateResourceReq := &gresources.UpdateRequest{
		Description:        getResourceResp.Description,
		Active:             getResourceResp.Active,
		OriginGroup:        int(getResourceResp.OriginGroup),
		OriginProtocol:     getResourceResp.OriginProtocol,
		SecondaryHostnames: getResourceResp.SecondaryHostnames,
		SSlEnabled:         true,
		SSLData:            int(updateResourceCertId),
		ProxySSLEnabled:    getResourceResp.ProxySSLEnabled,
		ProxySSLCA:         &getResourceResp.ProxySSLCA,
		ProxySSLData:       &getResourceResp.ProxySSLData,
		Options:            getResourceResp.Options,
	}
	updateResourceResp, err := d.sdkClient.Update(context.TODO(), d.config.ResourceId, updateResourceReq)
	d.logger.Debug("sdk request 'resources.Update'", slog.Int64("resourceId", d.config.ResourceId), slog.Any("request", updateResourceReq), slog.Any("response", updateResourceResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'resources.Update'")
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiToken string) (*gresources.Service, error) {
	if apiToken == "" {
		return nil, errors.New("invalid gcore api token")
	}

	requester := gprovider.NewClient(
		gcoresdk.BASE_URL,
		gprovider.WithSigner(gcoresdk.NewAuthRequestSigner(apiToken)),
	)
	service := gresources.NewService(requester)
	return service, nil
}
