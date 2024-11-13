package volcenginelive

import (
	"context"
	"errors"
	"fmt"
	xerrors "github.com/pkg/errors"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/cast"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
	live "github.com/volcengine/volc-sdk-golang/service/live/v20230101"
	"strings"
	"time"
)

type VolcengineLiveUploaderConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	Region          string `json:"region"`
}

type VolcengineLiveUploader struct {
	config    *VolcengineLiveUploaderConfig
	sdkClient *live.Live
}

var _ uploader.Uploader = (*VolcengineLiveUploader)(nil)

func New(config *VolcengineLiveUploaderConfig) (*VolcengineLiveUploader, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	client := live.NewInstance()
	client.SetAccessKey(config.AccessKeyId)
	client.SetSecretKey(config.AccessKeySecret)

	return &VolcengineLiveUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *VolcengineLiveUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := x509.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}
	apiCtx := context.Background()
	// 查询证书列表，避免重复上传
	// REF: https://www.volcengine.com/docs/6469/1186278#%E6%9F%A5%E8%AF%A2%E8%AF%81%E4%B9%A6%E5%88%97%E8%A1%A8
	describeServerCertificatesReq := &live.ListCertV2Body{}
	describeServerCertificatesResp, err := u.sdkClient.ListCertV2(apiCtx, describeServerCertificatesReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'live.ListCertV2'")
	}

	if describeServerCertificatesResp.Result.CertList != nil {
		for _, item := range describeServerCertificatesResp.Result.CertList {

			sslDetailReq := &live.DescribeCertDetailSecretV2Body{
				ChainID: cast.StringPtr(item.ChainID),
			}
			// 查询证书详细信息
			// REF: https://www.volcengine.com/docs/6469/1186278#%E6%9F%A5%E7%9C%8B%E8%AF%81%E4%B9%A6%E8%AF%A6%E6%83%85
			certDetail, detailErr := u.sdkClient.DescribeCertDetailSecretV2(apiCtx, sslDetailReq)
			if detailErr != nil {
				continue
			}
			var isSameCert bool
			certificate := strings.Join(certDetail.Result.SSL.Chain, "\n\n")
			if certificate == certPem {
				isSameCert = true
			} else {
				cert, err := x509.ParseCertificateFromPEM(certificate)
				if err != nil {
					continue
				}

				isSameCert = x509.EqualCertificate(cert, certX509)
			}
			// 如果已存在相同证书，直接返回已有的证书信息
			if isSameCert {
				return &uploader.UploadResult{
					CertId:   item.ChainID,
					CertName: item.CertName,
				}, nil
			}
		}
	}

	// 生成新证书名（需符合火山引擎命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())
	// 上传新证书
	// REF: https://www.volcengine.com/docs/6469/1186278#%E6%B7%BB%E5%8A%A0%E8%AF%81%E4%B9%A6
	createCertReq := &live.CreateCertBody{
		CertName:    &certName,
		UseWay:      `https`,
		ProjectName: cast.StringPtr(`default`),
		Rsa: live.CreateCertBodyRsa{
			Prikey: privkeyPem,
			Pubkey: certPem,
		},
	}
	createCertResp, err := u.sdkClient.CreateCert(apiCtx, createCertReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'live.CreateCert'")
	}

	certId = *createCertResp.Result.ChainID
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}
