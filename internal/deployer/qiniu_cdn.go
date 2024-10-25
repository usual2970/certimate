package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/qiniu/go-sdk/v7/auth"

	"github.com/usual2970/certimate/internal/domain"
	xhttp "github.com/usual2970/certimate/internal/utils/http"
)

const qiniuGateway = "http://api.qiniu.com"

type QiniuCDNDeployer struct {
	option      *DeployerOption
	info        []string
	credentials *auth.Credentials
}

func NewQiniuCDNDeployer(option *DeployerOption) (*QiniuCDNDeployer, error) {
	access := &domain.QiniuAccess{}
	json.Unmarshal([]byte(option.Access), access)

	return &QiniuCDNDeployer{
		option: option,
		info:   make([]string, 0),

		credentials: auth.New(access.AccessKey, access.SecretKey),
	}, nil
}

func (d *QiniuCDNDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *QiniuCDNDeployer) GetInfo() []string {
	return d.info
}

func (d *QiniuCDNDeployer) Deploy(ctx context.Context) error {
	// 上传证书
	certId, err := d.uploadCert()
	if err != nil {
		return fmt.Errorf("uploadCert failed: %w", err)
	}

	// 获取域名信息
	domainInfo, err := d.getDomainInfo()
	if err != nil {
		return fmt.Errorf("getDomainInfo failed: %w", err)
	}

	// 判断域名是否启用 https
	if domainInfo.Https != nil && domainInfo.Https.CertID != "" {
		// 启用了 https
		// 修改域名证书
		err = d.modifyDomainCert(certId, domainInfo.Https.ForceHttps, domainInfo.Https.Http2Enable)
		if err != nil {
			return fmt.Errorf("modifyDomainCert failed: %w", err)
		}
	} else {
		// 没启用 https
		// 启用 https
		err = d.enableHttps(certId)
		if err != nil {
			return fmt.Errorf("enableHttps failed: %w", err)
		}
	}

	return nil
}

func (d *QiniuCDNDeployer) enableHttps(certId string) error {
	domain := d.option.DeployConfig.GetDomain()
	path := fmt.Sprintf("/domain/%s/sslize", domain)

	body := &qiniuModifyDomainCertReq{
		CertID:      certId,
		ForceHttps:  true,
		Http2Enable: true,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("enable https failed: %w", err)
	}

	_, err = d.req(qiniuGateway+path, http.MethodPut, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("enable https failed: %w", err)
	}

	return nil
}

type qiniuDomainInfo struct {
	Https *qiniuModifyDomainCertReq `json:"https"`
}

func (d *QiniuCDNDeployer) getDomainInfo() (*qiniuDomainInfo, error) {
	domain := d.option.DeployConfig.GetDomain()

	path := fmt.Sprintf("/domain/%s", domain)

	res, err := d.req(qiniuGateway+path, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("req failed: %w", err)
	}

	resp := &qiniuDomainInfo{}
	err = json.Unmarshal(res, resp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	return resp, nil
}

type qiniuUploadCertReq struct {
	Name       string `json:"name"`
	CommonName string `json:"common_name"`
	Pri        string `json:"pri"`
	Ca         string `json:"ca"`
}

type qiniuUploadCertResp struct {
	CertID string `json:"certID"`
}

func (d *QiniuCDNDeployer) uploadCert() (string, error) {
	path := "/sslcert"

	body := &qiniuUploadCertReq{
		Name:       getDeployString(d.option.DeployConfig, "domain"),
		CommonName: getDeployString(d.option.DeployConfig, "domain"),
		Pri:        d.option.Certificate.PrivateKey,
		Ca:         d.option.Certificate.Certificate,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("json.Marshal failed: %w", err)
	}

	res, err := d.req(qiniuGateway+path, http.MethodPost, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("req failed: %w", err)
	}
	resp := &qiniuUploadCertResp{}
	err = json.Unmarshal(res, resp)
	if err != nil {
		return "", fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	return resp.CertID, nil
}

type qiniuModifyDomainCertReq struct {
	CertID      string `json:"certId"`
	ForceHttps  bool   `json:"forceHttps"`
	Http2Enable bool   `json:"http2Enable"`
}

func (d *QiniuCDNDeployer) modifyDomainCert(certId string, forceHttps, http2Enable bool) error {
	domain := d.option.DeployConfig.GetDomain()
	path := fmt.Sprintf("/domain/%s/httpsconf", domain)

	body := &qiniuModifyDomainCertReq{
		CertID:      certId,
		ForceHttps:  forceHttps,
		Http2Enable: http2Enable,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json.Marshal failed: %w", err)
	}

	_, err = d.req(qiniuGateway+path, http.MethodPut, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("req failed: %w", err)
	}

	return nil
}

func (d *QiniuCDNDeployer) req(url, method string, body io.Reader) ([]byte, error) {
	req := xhttp.BuildReq(url, method, body, map[string]string{
		"Content-Type": "application/json",
	})

	if err := d.credentials.AddToken(auth.TokenQBox, req); err != nil {
		return nil, fmt.Errorf("credentials.AddToken failed: %w", err)
	}

	respBody, err := xhttp.ToRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ToRequest failed: %w", err)
	}

	defer respBody.Close()

	res, err := io.ReadAll(respBody)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll failed: %w", err)
	}

	return res, nil
}
