package baiducloudcert

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/certutil"
	bdsdk "github.com/usual2970/certimate/internal/pkg/vendors/baiducloud-sdk/cert"
)

type UploaderConfig struct {
	// 百度智能云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 百度智能云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *bdsdk.Client
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.SecretAccessKey)
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
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 遍历证书列表，避免重复上传
	// REF: https://cloud.baidu.com/doc/Reference/s/Gjwvz27xu#35-%E6%9F%A5%E7%9C%8B%E8%AF%81%E4%B9%A6%E5%88%97%E8%A1%A8%E8%AF%A6%E6%83%85
	listCertDetail, err := u.sdkClient.ListCertDetail()
	u.logger.Debug("sdk request 'cert.ListCertDetail'", slog.Any("response", listCertDetail))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cert.ListCertDetail'")
	} else {
		for _, certDetail := range listCertDetail.Certs {
			// 先对比证书通用名称
			if !strings.EqualFold(certX509.Subject.CommonName, certDetail.CertCommonName) {
				continue
			}

			// 再对比证书有效期
			oldCertNotBefore, _ := time.Parse("2006-01-02T15:04:05Z", certDetail.CertStartTime)
			oldCertNotAfter, _ := time.Parse("2006-01-02T15:04:05Z", certDetail.CertStopTime)
			if !certX509.NotBefore.Equal(oldCertNotBefore) || !certX509.NotAfter.Equal(oldCertNotAfter) {
				continue
			}

			// 再对比证书多域名
			if certDetail.CertDNSNames != strings.Join(certX509.DNSNames, ",") {
				continue
			}

			// 最后对比证书内容
			getCertDetailResp, err := u.sdkClient.GetCertRawData(certDetail.CertId)
			u.logger.Debug("sdk request 'cert.GetCertRawData'", slog.Any("certId", certDetail.CertId), slog.Any("response", getCertDetailResp))
			if err != nil {
				return nil, xerrors.Wrap(err, "failed to execute sdk request 'cert.GetCertRawData'")
			} else {
				oldCertX509, err := certutil.ParseCertificateFromPEM(getCertDetailResp.CertServerData)
				if err != nil {
					continue
				}
				if !certutil.EqualCertificate(certX509, oldCertX509) {
					continue
				}
			}

			// 如果以上信息都一致，则视为已存在相同证书，直接返回
			u.logger.Info("ssl certificate already exists")
			return &uploader.UploadResult{
				CertId:   certDetail.CertId,
				CertName: certDetail.CertName,
			}, nil
		}
	}

	// 创建证书
	// REF: https://cloud.baidu.com/doc/Reference/s/Gjwvz27xu#31-%E5%88%9B%E5%BB%BA%E8%AF%81%E4%B9%A6
	createCertReq := &bdsdk.CreateCertArgs{}
	createCertReq.CertName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())
	createCertReq.CertServerData = certPem
	createCertReq.CertPrivateData = privkeyPem
	createCertResp, err := u.sdkClient.CreateCert(createCertReq)
	u.logger.Debug("sdk request 'cert.CreateCert'", slog.Any("request", createCertReq), slog.Any("response", createCertResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cert.CreateCert'")
	}

	return &uploader.UploadResult{
		CertId:   createCertResp.CertId,
		CertName: createCertResp.CertName,
	}, nil
}

func createSdkClient(accessKeyId, secretAccessKey string) (*bdsdk.Client, error) {
	client, err := bdsdk.NewClient(accessKeyId, secretAccessKey, "")
	if err != nil {
		return nil, err
	}

	return client, nil
}
