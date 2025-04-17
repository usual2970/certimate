package volcenginecertcenter

import (
	"context"
	"log/slog"

	xerrors "github.com/pkg/errors"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"
	vesession "github.com/volcengine/volcengine-go-sdk/volcengine/session"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	veccsdk "github.com/usual2970/certimate/internal/pkg/vendors/volcengine-sdk/certcenter"
)

type UploaderConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 火山引擎地域。
	Region string `json:"region"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *veccsdk.CertCenter
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &UploaderProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (u *UploaderProvider) WithLogger(logger *slog.Logger) uploader.Uploader {
	if logger == nil {
		u.logger = slog.Default()
	} else {
		u.logger = logger
	}
	return u
}

func (u *UploaderProvider) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 上传证书
	// REF: https://www.volcengine.com/docs/6638/1365580
	importCertificateReq := &veccsdk.ImportCertificateInput{
		CertificateInfo: &veccsdk.ImportCertificateInputCertificateInfo{
			CertificateChain: ve.String(certPem),
			PrivateKey:       ve.String(privkeyPem),
		},
		Repeatable: ve.Bool(false),
	}
	importCertificateResp, err := u.sdkClient.ImportCertificate(importCertificateReq)
	u.logger.Debug("sdk request 'certcenter.ImportCertificate'", slog.Any("request", importCertificateReq), slog.Any("response", importCertificateResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'certcenter.ImportCertificate'")
	}

	var certId string
	if importCertificateResp.InstanceId != nil && *importCertificateResp.InstanceId != "" {
		certId = *importCertificateResp.InstanceId
	}
	if importCertificateResp.RepeatId != nil && *importCertificateResp.RepeatId != "" {
		certId = *importCertificateResp.RepeatId
	}

	if certId == "" {
		return nil, xerrors.New("failed to get certificate id, both `InstanceId` and `RepeatId` are empty")
	}

	return &uploader.UploadResult{
		CertId: certId,
	}, nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*veccsdk.CertCenter, error) {
	if region == "" {
		region = "cn-beijing" // 证书中心默认区域：北京
	}

	config := ve.NewConfig().WithRegion(region).WithAkSk(accessKeyId, accessKeySecret)

	session, err := vesession.NewSession(config)
	if err != nil {
		return nil, err
	}

	client := veccsdk.New(session)
	return client, nil
}
