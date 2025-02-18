package volcenginelive

import (
	"context"
	"fmt"
	"strings"
	"time"

	xerrors "github.com/pkg/errors"
	veLive "github.com/volcengine/volc-sdk-golang/service/live/v20230101"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/certs"
)

type VolcEngineLiveUploaderConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
}

type VolcEngineLiveUploader struct {
	config    *VolcEngineLiveUploaderConfig
	sdkClient *veLive.Live
}

var _ uploader.Uploader = (*VolcEngineLiveUploader)(nil)

func New(config *VolcEngineLiveUploaderConfig) (*VolcEngineLiveUploader, error) {
	if config == nil {
		panic("config is nil")
	}

	client := veLive.NewInstance()
	client.SetAccessKey(config.AccessKeyId)
	client.SetSecretKey(config.AccessKeySecret)

	return &VolcEngineLiveUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *VolcEngineLiveUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := certs.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 查询证书列表，避免重复上传
	// REF: https://www.volcengine.com/docs/6469/1186278#%E6%9F%A5%E8%AF%A2%E8%AF%81%E4%B9%A6%E5%88%97%E8%A1%A8
	listCertReq := &veLive.ListCertV2Body{}
	listCertResp, err := u.sdkClient.ListCertV2(ctx, listCertReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'live.ListCertV2'")
	}
	if listCertResp.Result.CertList != nil {
		for _, certDetail := range listCertResp.Result.CertList {
			// 查询证书详细信息
			// REF: https://www.volcengine.com/docs/6469/1186278#%E6%9F%A5%E7%9C%8B%E8%AF%81%E4%B9%A6%E8%AF%A6%E6%83%85
			describeCertDetailSecretReq := &veLive.DescribeCertDetailSecretV2Body{
				ChainID: ve.String(certDetail.ChainID),
			}
			describeCertDetailSecretResp, err := u.sdkClient.DescribeCertDetailSecretV2(ctx, describeCertDetailSecretReq)
			if err != nil {
				continue
			}

			var isSameCert bool
			certificate := strings.Join(describeCertDetailSecretResp.Result.SSL.Chain, "\n\n")
			if certificate == certPem {
				isSameCert = true
			} else {
				oldCertX509, err := certs.ParseCertificateFromPEM(certificate)
				if err != nil {
					continue
				}

				isSameCert = certs.EqualCertificate(certX509, oldCertX509)
			}

			// 如果已存在相同证书，直接返回已有的证书信息
			if isSameCert {
				return &uploader.UploadResult{
					CertId:   certDetail.ChainID,
					CertName: certDetail.CertName,
				}, nil
			}
		}
	}

	// 生成新证书名（需符合火山引擎命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://www.volcengine.com/docs/6469/1186278#%E6%B7%BB%E5%8A%A0%E8%AF%81%E4%B9%A6
	createCertReq := &veLive.CreateCertBody{
		CertName:    ve.String(certName),
		UseWay:      "https",
		ProjectName: ve.String("default"),
		Rsa: veLive.CreateCertBodyRsa{
			Prikey: privkeyPem,
			Pubkey: certPem,
		},
	}
	createCertResp, err := u.sdkClient.CreateCert(ctx, createCertReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'live.CreateCert'")
	}

	certId = *createCertResp.Result.ChainID
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}
