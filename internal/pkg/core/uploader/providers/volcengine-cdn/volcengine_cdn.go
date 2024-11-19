package volcenginecdn

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	xerrors "github.com/pkg/errors"
	veCdn "github.com/volcengine/volc-sdk-golang/service/cdn"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/cast"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
)

type VolcEngineCDNUploaderConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type VolcEngineCDNUploader struct {
	config    *VolcEngineCDNUploaderConfig
	sdkClient *veCdn.CDN
}

var _ uploader.Uploader = (*VolcEngineCDNUploader)(nil)

func New(config *VolcEngineCDNUploaderConfig) (*VolcEngineCDNUploader, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	client := veCdn.NewInstance()
	client.Client.SetAccessKey(config.AccessKeyId)
	client.Client.SetSecretKey(config.AccessKeySecret)

	return &VolcEngineCDNUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *VolcEngineCDNUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := x509.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 查询证书列表，避免重复上传
	// REF: https://www.volcengine.com/docs/6454/125709
	listCertInfoPageNum := int64(1)
	listCertInfoPageSize := int64(100)
	listCertInfoTotal := 0
	listCertInfoReq := &veCdn.ListCertInfoRequest{
		PageNum:  cast.Int64Ptr(listCertInfoPageNum),
		PageSize: cast.Int64Ptr(listCertInfoPageSize),
		Source:   "volc_cert_center",
	}
	for {
		listCertInfoResp, err := u.sdkClient.ListCertInfo(listCertInfoReq)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.ListCertInfo'")
		}

		if listCertInfoResp.Result.CertInfo != nil {
			for _, certDetail := range listCertInfoResp.Result.CertInfo {
				fingerprintSha1 := sha1.Sum(certX509.Raw)
				fingerprintSha256 := sha256.Sum256(certX509.Raw)
				isSameCert := strings.EqualFold(hex.EncodeToString(fingerprintSha1[:]), certDetail.CertFingerprint.Sha1) &&
					strings.EqualFold(hex.EncodeToString(fingerprintSha256[:]), certDetail.CertFingerprint.Sha256)
				// 如果已存在相同证书，直接返回已有的证书信息
				if isSameCert {
					return &uploader.UploadResult{
						CertId:   certDetail.CertId,
						CertName: certDetail.Desc,
					}, nil
				}
			}
		}

		listCertInfoLen := len(listCertInfoResp.Result.CertInfo)
		if listCertInfoLen < int(listCertInfoPageSize) || int(listCertInfoResp.Result.Total) <= listCertInfoTotal+listCertInfoLen {
			break
		} else {
			listCertInfoPageNum++
			listCertInfoTotal += listCertInfoLen
		}
	}

	// 生成新证书名（需符合火山引擎命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://www.volcengine.com/docs/6454/1245763
	addCertificateReq := &veCdn.AddCertificateRequest{
		Certificate: certPem,
		PrivateKey:  privkeyPem,
		Source:      cast.StringPtr("volc_cert_center"),
		Desc:        cast.StringPtr(certName),
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
