package jdcloudssl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	wangsusdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/wangsu/certificate"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
	typeutil "github.com/usual2970/certimate/internal/pkg/utils/type"
)

type UploaderConfig struct {
	// 网宿云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 网宿云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *wangsusdk.Client
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

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
	// REF: https://www.wangsu.com/document/api-doc/22675?productCode=certificatemanagement
	listCertificatesResp, err := u.sdkClient.ListCertificates()
	u.logger.Debug("sdk request 'certificatemanagement.ListCertificates'", slog.Any("response", listCertificatesResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'certificatemanagement.ListCertificates': %w", err)
	}

	if listCertificatesResp.Certificates != nil {
		for _, certificate := range listCertificatesResp.Certificates {
			// 对比证书序列号
			if !strings.EqualFold(certX509.SerialNumber.Text(16), certificate.Serial) {
				continue
			}

			// 再对比证书有效期
			cstzone := time.FixedZone("CST", 8*60*60)
			oldCertNotBefore, _ := time.ParseInLocation(time.DateTime, certificate.ValidityFrom, cstzone)
			oldCertNotAfter, _ := time.ParseInLocation(time.DateTime, certificate.ValidityTo, cstzone)
			if !certX509.NotBefore.Equal(oldCertNotBefore) || !certX509.NotAfter.Equal(oldCertNotAfter) {
				continue
			}

			// 如果以上信息都一致，则视为已存在相同证书，直接返回
			u.logger.Info("ssl certificate already exists")
			return &uploader.UploadResult{
				CertId:   certificate.CertificateId,
				CertName: certificate.Name,
			}, nil
		}
	}

	// 生成新证书名（需符合网宿云命名规则）
	var certId string
	certName := fmt.Sprintf("certimate_%d", time.Now().UnixMilli())

	// 新增证书
	// REF: https://www.wangsu.com/document/api-doc/25199?productCode=certificatemanagement
	createCertificateReq := &wangsusdk.CreateCertificateRequest{
		Name:        typeutil.ToPtr(certName),
		Certificate: typeutil.ToPtr(certPEM),
		PrivateKey:  typeutil.ToPtr(privkeyPEM),
		Comment:     typeutil.ToPtr("upload from certimate"),
	}
	createCertificateResp, err := u.sdkClient.CreateCertificate(createCertificateReq)
	u.logger.Debug("sdk request 'certificatemanagement.CreateCertificate'", slog.Any("request", createCertificateReq), slog.Any("response", createCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'certificatemanagement.CreateCertificate': %w", err)
	}

	// 网宿云证书 URL 中包含证书 ID
	// 格式：
	//    https://open.chinanetcenter.com/api/certificate/100001
	wangsuCertIdMatches := regexp.MustCompile(`/certificate/([0-9]+)`).FindStringSubmatch(createCertificateResp.CertificateUrl)
	if len(wangsuCertIdMatches) > 1 {
		certId = wangsuCertIdMatches[1]
	} else {
		return nil, fmt.Errorf("received empty certificate id")
	}

	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*wangsusdk.Client, error) {
	if accessKeyId == "" {
		return nil, errors.New("invalid wangsu access key id")
	}

	if accessKeySecret == "" {
		return nil, errors.New("invalid wangsu access key secret")
	}

	return wangsusdk.NewClient(accessKeyId, accessKeySecret), nil
}
