package gcorecdn

import (
	"context"
	"errors"
	"fmt"
	"time"

	gprovider "github.com/G-Core/gcorelabscdn-go/gcore/provider"
	gsslcerts "github.com/G-Core/gcorelabscdn-go/sslcerts"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	gcoresdk "github.com/usual2970/certimate/internal/pkg/vendors/gcore-sdk/common"
)

type GcoreCDNUploaderConfig struct {
	// Gcore API Token。
	ApiToken string `json:"apiToken"`
}

type GcoreCDNUploader struct {
	config    *GcoreCDNUploaderConfig
	sdkClient *gsslcerts.Service
}

var _ uploader.Uploader = (*GcoreCDNUploader)(nil)

func New(config *GcoreCDNUploaderConfig) (*GcoreCDNUploader, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ApiToken)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &GcoreCDNUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *GcoreCDNUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 生成新证书名（需符合 Gcore 命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate_%d", time.Now().UnixMilli())

	// 新增证书
	// REF: https://api.gcore.com/docs/cdn#tag/CA-certificates/operation/ca_certificates-add
	createCertificateReq := &gsslcerts.CreateRequest{
		Name:           certName,
		Cert:           certPem,
		PrivateKey:     privkeyPem,
		Automated:      false,
		ValidateRootCA: false,
	}
	createCertificateResp, err := u.sdkClient.Create(context.TODO(), createCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'sslcerts.Create'")
	}

	certId = fmt.Sprintf("%d", createCertificateResp.ID)
	certName = createCertificateResp.Name
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createSdkClient(apiToken string) (*gsslcerts.Service, error) {
	if apiToken == "" {
		return nil, errors.New("invalid gcore api token")
	}

	requester := gprovider.NewClient(
		gcoresdk.BASE_URL,
		gprovider.WithSigner(gcoresdk.NewAuthRequestSigner(apiToken)),
	)
	service := gsslcerts.NewService(requester)
	return service, nil
}
