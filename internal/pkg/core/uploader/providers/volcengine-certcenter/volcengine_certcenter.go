package volcenginecertcenter

import (
	"context"
	"errors"

	xerrors "github.com/pkg/errors"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"
	veSession "github.com/volcengine/volcengine-go-sdk/volcengine/session"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	veCertCenter "github.com/usual2970/certimate/internal/pkg/vendors/volcengine-sdk/certcenter"
)

type VolcEngineCertCenterUploaderConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 火山引擎区域。
	Region string `json:"region"`
}

type VolcEngineCertCenterUploader struct {
	config    *VolcEngineCertCenterUploaderConfig
	sdkClient *veCertCenter.CertCenter
}

var _ uploader.Uploader = (*VolcEngineCertCenterUploader)(nil)

func New(config *VolcEngineCertCenterUploaderConfig) (*VolcEngineCertCenterUploader, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client: %w")
	}

	return &VolcEngineCertCenterUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *VolcEngineCertCenterUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 上传证书
	// REF: https://www.volcengine.com/docs/6638/1365580
	importCertificateReq := &veCertCenter.ImportCertificateInput{
		CertificateInfo: &veCertCenter.ImportCertificateInputCertificateInfo{
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

func createSdkClient(accessKeyId, accessKeySecret, region string) (*veCertCenter.CertCenter, error) {
	if region == "" {
		region = "cn-beijing"
	}

	config := ve.NewConfig().WithRegion(region).WithAkSk(accessKeyId, accessKeySecret)

	session, err := veSession.NewSession(config)
	if err != nil {
		return nil, err
	}

	client := veCertCenter.New(session)
	return client, nil
}
