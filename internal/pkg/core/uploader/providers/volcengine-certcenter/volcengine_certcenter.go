package volcenginecertcenter

import (
	"context"

	xerrors "github.com/pkg/errors"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"
	veSession "github.com/volcengine/volcengine-go-sdk/volcengine/session"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	vesdkCc "github.com/usual2970/certimate/internal/pkg/vendors/volcengine-sdk/certcenter"
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
	sdkClient *vesdkCc.CertCenter
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
		sdkClient: client,
	}, nil
}

func (u *UploaderProvider) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 上传证书
	// REF: https://www.volcengine.com/docs/6638/1365580
	importCertificateReq := &vesdkCc.ImportCertificateInput{
		CertificateInfo: &vesdkCc.ImportCertificateInputCertificateInfo{
			CertificateChain: ve.String(certPem),
			PrivateKey:       ve.String(privkeyPem),
		},
		Repeatable: ve.Bool(false),
	}
	importCertificateResp, err := u.sdkClient.ImportCertificate(importCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'certcenter.ImportCertificate'")
	}

	var certId string
	if importCertificateResp.InstanceId != nil {
		certId = *importCertificateResp.InstanceId
	}
	if importCertificateResp.RepeatId != nil {
		certId = *importCertificateResp.RepeatId
	}
	return &uploader.UploadResult{
		CertId: certId,
	}, nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*vesdkCc.CertCenter, error) {
	if region == "" {
		region = "cn-beijing" // 证书中心默认区域：北京
	}

	config := ve.NewConfig().WithRegion(region).WithAkSk(accessKeyId, accessKeySecret)

	session, err := veSession.NewSession(config)
	if err != nil {
		return nil, err
	}

	client := vesdkCc.New(session)
	return client, nil
}
