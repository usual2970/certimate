package bytepluscdn

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	bpCdn "github.com/byteplus-sdk/byteplus-sdk-golang/service/cdn"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/certs"
)

type ByteplusCDNUploaderConfig struct {
	// BytePlus AccessKey。
	AccessKey string `json:"accessKey"`
	// BytePlus SecretKey。
	SecretKey string `json:"secretKey"`
}

type ByteplusCDNUploader struct {
	config    *ByteplusCDNUploaderConfig
	sdkClient *bpCdn.CDN
}

var _ uploader.Uploader = (*ByteplusCDNUploader)(nil)

func New(config *ByteplusCDNUploaderConfig) (*ByteplusCDNUploader, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	client := bpCdn.NewInstance()
	client.Client.SetAccessKey(config.AccessKey)
	client.Client.SetSecretKey(config.SecretKey)

	return &ByteplusCDNUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *ByteplusCDNUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := certs.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 查询证书列表，避免重复上传
	// REF: https://docs.byteplus.com/en/docs/byteplus-cdn/reference-listcertinfo
	listCertInfoPageNum := int64(1)
	listCertInfoPageSize := int64(100)
	listCertInfoTotal := 0
	listCertInfoReq := &bpCdn.ListCertInfoRequest{
		PageNum:  bpCdn.GetInt64Ptr(listCertInfoPageNum),
		PageSize: bpCdn.GetInt64Ptr(listCertInfoPageSize),
		Source:   bpCdn.GetStrPtr("cert_center"),
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

	// 生成新证书名（需符合 BytePlus 命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://docs.byteplus.com/en/docs/byteplus-cdn/reference-addcertificate
	addCertificateReq := &bpCdn.AddCertificateRequest{
		Certificate: certPem,
		PrivateKey:  privkeyPem,
		Source:      bpCdn.GetStrPtr("cert_center"),
		Desc:        bpCdn.GetStrPtr(certName),
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
