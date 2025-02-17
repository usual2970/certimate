package gcorecdn

import (
	"context"
	"errors"
	"strconv"

	gprovider "github.com/G-Core/gcorelabscdn-go/gcore/provider"
	gresources "github.com/G-Core/gcorelabscdn-go/resources"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/gcore-cdn"
	gcoresdk "github.com/usual2970/certimate/internal/pkg/vendors/gcore-sdk/common"
)

type GcoreCDNDeployerConfig struct {
	// Gcore API Token。
	ApiToken string `json:"apiToken"`
	// CDN 资源 ID。
	ResourceId int64 `json:"resourceId"`
}

type GcoreCDNDeployer struct {
	config      *GcoreCDNDeployerConfig
	logger      logger.Logger
	sdkClient   *gresources.Service
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*GcoreCDNDeployer)(nil)

func New(config *GcoreCDNDeployerConfig) (*GcoreCDNDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *GcoreCDNDeployerConfig, logger logger.Logger) (*GcoreCDNDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.ApiToken)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.New(&uploadersp.GcoreCDNUploaderConfig{
		ApiToken: config.ApiToken,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &GcoreCDNDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *GcoreCDNDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	if d.config.ResourceId == 0 {
		return nil, errors.New("config `resourceId` is required")
	}

	// 上传证书到 CDN
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Logt("certificate file uploaded", upres)
	}

	// 获取 CDN 资源详情
	// REF: https://api.gcore.com/docs/cdn#tag/CDN-resources/paths/~1cdn~1resources~1%7Bresource_id%7D/get
	getResourceResp, err := d.sdkClient.Get(context.TODO(), d.config.ResourceId)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'resources.Get'")
	} else {
		d.logger.Logt("已获取 CDN 资源详情", getResourceResp)
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
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'resources.Update'")
	} else {
		d.logger.Logt("已更新 CDN 资源详情", updateResourceResp)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiToken string) (*gresources.Service, error) {
	requester := gprovider.NewClient(
		gcoresdk.BASE_URL,
		gprovider.WithSigner(gcoresdk.NewAuthRequestSigner(apiToken)),
	)
	service := gresources.NewService(requester)
	return service, nil
}
