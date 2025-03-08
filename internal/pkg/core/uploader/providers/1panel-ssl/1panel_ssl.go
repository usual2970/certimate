package onepanelssl

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	opsdk "github.com/usual2970/certimate/internal/pkg/vendors/1panel-sdk"
)

type UploaderConfig struct {
	// 1Panel 地址。
	ApiUrl string `json:"apiUrl"`
	// 1Panel 接口密钥。
	ApiKey string `json:"apiKey"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	sdkClient *opsdk.Client
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ApiUrl, config.ApiKey)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &UploaderProvider{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *UploaderProvider) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 遍历证书列表，避免重复上传
	if res, err := u.getExistCert(ctx, certPem, privkeyPem); err != nil {
		return nil, err
	} else if res != nil {
		return res, nil
	}

	// 生成新证书名（需符合 1Panel 命名规则）
	certName := fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传证书
	uploadWebsiteSSLReq := &opsdk.UploadWebsiteSSLRequest{
		Type:        "paste",
		Description: certName,
		Certificate: certPem,
		PrivateKey:  privkeyPem,
	}
	uploadWebsiteSSLResp, err := u.sdkClient.UploadWebsiteSSL(uploadWebsiteSSLReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request '1panel.UploadWebsiteSSL'")
	}

	// 遍历证书列表，获取刚刚上传证书 ID
	if res, err := u.getExistCert(ctx, certPem, privkeyPem); err != nil {
		return nil, err
	} else if res == nil {
		return nil, fmt.Errorf("no ssl certificate found, may be upload failed (code: %d, message: %s)", uploadWebsiteSSLResp.GetCode(), uploadWebsiteSSLResp.GetMessage())
	} else {
		return res, nil
	}
}

func (u *UploaderProvider) getExistCert(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	searchWebsiteSSLPageNumber := int32(1)
	searchWebsiteSSLPageSize := int32(100)
	for {
		searchWebsiteSSLReq := &opsdk.SearchWebsiteSSLRequest{
			Page:     searchWebsiteSSLPageNumber,
			PageSize: searchWebsiteSSLPageSize,
		}
		searchWebsiteSSLResp, err := u.sdkClient.SearchWebsiteSSL(searchWebsiteSSLReq)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request '1panel.SearchWebsiteSSL'")
		}

		for _, sslItem := range searchWebsiteSSLResp.Data.Items {
			if strings.TrimSpace(sslItem.PEM) == strings.TrimSpace(certPem) &&
				strings.TrimSpace(sslItem.PrivateKey) == strings.TrimSpace(privkeyPem) {
				// 如果已存在相同证书，直接返回已有的证书信息
				return &uploader.UploadResult{
					CertId:   fmt.Sprintf("%d", sslItem.ID),
					CertName: sslItem.Description,
				}, nil
			}
		}

		if len(searchWebsiteSSLResp.Data.Items) < int(searchWebsiteSSLPageSize) {
			break
		} else {
			searchWebsiteSSLPageNumber++
		}
	}

	return nil, nil
}

func createSdkClient(apiUrl, apiKey string) (*opsdk.Client, error) {
	if _, err := url.Parse(apiUrl); err != nil {
		return nil, errors.New("invalid 1panel api url")
	}

	if apiKey == "" {
		return nil, errors.New("invalid 1panel api key")
	}

	client := opsdk.NewClient(apiUrl, apiKey)
	return client, nil
}
