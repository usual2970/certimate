package volcenginecdn

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	xerrors "github.com/pkg/errors"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
	"github.com/volcengine/volc-sdk-golang/service/cdn"
)

type VolcengineCDNUploaderConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type VolcengineCDNUploader struct {
	config    *VolcengineCDNUploaderConfig
	sdkClient *cdn.CDN
}

var _ uploader.Uploader = (*VolcengineCDNUploader)(nil)

func New(config *VolcengineCDNUploaderConfig) (*VolcengineCDNUploader, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	instance := cdn.NewInstance()
	client := instance.Client
	client.SetAccessKey(config.AccessKeyId)
	client.SetSecretKey(config.AccessKeySecret)

	return &VolcengineCDNUploader{
		config:    config,
		sdkClient: instance,
	}, nil
}

func (u *VolcengineCDNUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := x509.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}
	// 查询证书列表，避免重复上传
	// REF: https://www.volcengine.com/docs/6454/125709
	pageNum := int64(1)
	pageSize := int64(100)
	certSource := "volc_cert_center"
	listCertInfoReq := &cdn.ListCertInfoRequest{
		PageNum:  &pageNum,
		PageSize: &pageSize,
		Source:   certSource,
	}
	searchTotal := 0
	for {
		listCertInfoResp, err := u.sdkClient.ListCertInfo(listCertInfoReq)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.ListCertInfo'")
		}

		if listCertInfoResp.Result.CertInfo != nil {
			for _, certDetail := range listCertInfoResp.Result.CertInfo {
				hash := sha256.Sum256(certX509.Raw)
				isSameCert := strings.EqualFold(hex.EncodeToString(hash[:]), certDetail.CertFingerprint.Sha256)
				// 如果已存在相同证书，直接返回已有的证书信息
				if isSameCert {
					return &uploader.UploadResult{
						CertId:   certDetail.CertId,
						CertName: certDetail.Desc,
					}, nil
				}
			}
		}

		searchTotal += len(listCertInfoResp.Result.CertInfo)
		if int(listCertInfoResp.Result.Total) > searchTotal {
			pageNum++
		} else {
			break
		}

	}
	// 生成新证书名（需符合火山引擎命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())
	// 上传新证书
	// REF: https://www.volcengine.com/docs/6454/1245763
	addCertificateReq := &cdn.AddCertificateRequest{
		Certificate: certPem,
		PrivateKey:  privkeyPem,
		Source:      &certSource,
		Desc:        &certName,
	}
	addCertificateResp, err := u.sdkClient.AddCertificate(addCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.AddCertificate'")
	}

	certId = addCertificateResp.Result.CertId
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}
