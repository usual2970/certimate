package volcenginelive

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	velive "github.com/volcengine/volc-sdk-golang/service/live/v20230101"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

type UploaderConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *velive.Live
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client := velive.NewInstance()
	client.SetAccessKey(config.AccessKeyId)
	client.SetSecretKey(config.AccessKeySecret)

	return &UploaderProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (u *UploaderProvider) WithLogger(logger *slog.Logger) uploader.Uploader {
	if logger == nil {
		u.logger = slog.New(slog.DiscardHandler)
	} else {
		u.logger = logger
	}
	return u
}

func (u *UploaderProvider) Upload(ctx context.Context, certPEM string, privkeyPEM string) (*uploader.UploadResult, error) {
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 查询证书列表，避免重复上传
	// REF: https://www.volcengine.com/docs/6469/1186278#%E6%9F%A5%E8%AF%A2%E8%AF%81%E4%B9%A6%E5%88%97%E8%A1%A8
	listCertReq := &velive.ListCertV2Body{}
	listCertResp, err := u.sdkClient.ListCertV2(ctx, listCertReq)
	u.logger.Debug("sdk request 'live.ListCertV2'", slog.Any("request", listCertReq), slog.Any("response", listCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'live.ListCertV2': %w", err)
	}
	if listCertResp.Result.CertList != nil {
		for _, certDetail := range listCertResp.Result.CertList {
			// 查询证书详细信息
			// REF: https://www.volcengine.com/docs/6469/1186278#%E6%9F%A5%E7%9C%8B%E8%AF%81%E4%B9%A6%E8%AF%A6%E6%83%85
			describeCertDetailSecretReq := &velive.DescribeCertDetailSecretV2Body{
				ChainID: ve.String(certDetail.ChainID),
			}
			describeCertDetailSecretResp, err := u.sdkClient.DescribeCertDetailSecretV2(ctx, describeCertDetailSecretReq)
			u.logger.Debug("sdk request 'live.DescribeCertDetailSecretV2'", slog.Any("request", describeCertDetailSecretReq), slog.Any("response", describeCertDetailSecretResp))
			if err != nil {
				continue
			}

			var isSameCert bool
			certificate := strings.Join(describeCertDetailSecretResp.Result.SSL.Chain, "\n\n")
			if certificate == certPEM {
				isSameCert = true
			} else {
				oldCertX509, err := certutil.ParseCertificateFromPEM(certificate)
				if err != nil {
					continue
				}

				isSameCert = certutil.EqualCertificate(certX509, oldCertX509)
			}

			// 如果已存在相同证书，直接返回
			if isSameCert {
				u.logger.Info("ssl certificate already exists")
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
	createCertReq := &velive.CreateCertBody{
		CertName:    ve.String(certName),
		UseWay:      "https",
		ProjectName: ve.String("default"),
		Rsa: velive.CreateCertBodyRsa{
			Prikey: privkeyPEM,
			Pubkey: certPEM,
		},
	}
	createCertResp, err := u.sdkClient.CreateCert(ctx, createCertReq)
	u.logger.Debug("sdk request 'live.CreateCert'", slog.Any("request", createCertReq), slog.Any("response", createCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'live.CreateCert': %w", err)
	}

	certId = *createCertResp.Result.ChainID
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}
