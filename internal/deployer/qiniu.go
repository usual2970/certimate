package deployer

import (
	"bytes"
	"certimate/internal/domain"
	xhttp "certimate/internal/utils/http"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/qiniu/go-sdk/v7/auth"
)

const qiniuGateway = "http://api.qiniu.com"

type qiuniu struct {
	option      *DeployerOption
	info        []string
	credentials *auth.Credentials
}

func NewQiNiu(option *DeployerOption) (*qiuniu, error) {
	access := &domain.QiniuAccess{}
	json.Unmarshal([]byte(option.Access), access)

	return &qiuniu{
		option: option,
		info:   make([]string, 0),

		credentials: auth.New(access.AccessKey, access.SecretKey),
	}, nil
}

func (a *qiuniu) GetID() string {
	return fmt.Sprintf("%s-%s", a.option.AceessRecord.GetString("name"), a.option.AceessRecord.Id)
}

func (q *qiuniu) GetInfo() []string {
	return q.info
}

func (q *qiuniu) Deploy(ctx context.Context) error {

	// 上传证书
	certId, err := q.uploadCert()
	if err != nil {
		return fmt.Errorf("uploadCert failed: %w", err)
	}

	// 获取域名信息
	domainInfo, err := q.getDomainInfo()
	if err != nil {
		return fmt.Errorf("getDomainInfo failed: %w", err)
	}

	// 判断域名是否启用 https

	if domainInfo.Https != nil && domainInfo.Https.CertID != "" {
		// 启用了 https
		// 修改域名证书
		err = q.modifyDomainCert(certId)
		if err != nil {
			return fmt.Errorf("modifyDomainCert failed: %w", err)
		}
	} else {
		// 没启用 https
		// 启用 https

		err = q.enableHttps(certId)
		if err != nil {
			return fmt.Errorf("enableHttps failed: %w", err)
		}
	}

	return nil
}

func (q *qiuniu) enableHttps(certId string) error {
	path := fmt.Sprintf("/domain/%s/sslize", getDeployString(q.option.DeployConfig, "domain"))

	body := &modifyDomainCertReq{
		CertID:      certId,
		ForceHttps:  true,
		Http2Enable: true,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("enable https failed: %w", err)
	}

	_, err = q.req(qiniuGateway+path, http.MethodPut, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("enable https failed: %w", err)
	}

	return nil
}

type domainInfo struct {
	Https *modifyDomainCertReq `json:"https"`
}

func (q *qiuniu) getDomainInfo() (*domainInfo, error) {
	path := fmt.Sprintf("/domain/%s", getDeployString(q.option.DeployConfig, "domain"))

	res, err := q.req(qiniuGateway+path, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("req failed: %w", err)
	}

	resp := &domainInfo{}
	err = json.Unmarshal(res, resp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	return resp, nil
}

type uploadCertReq struct {
	Name       string `json:"name"`
	CommonName string `json:"common_name"`
	Pri        string `json:"pri"`
	Ca         string `json:"ca"`
}

type uploadCertResp struct {
	CertID string `json:"certID"`
}

func (q *qiuniu) uploadCert() (string, error) {
	path := "/sslcert"

	body := &uploadCertReq{
		Name:       getDeployString(q.option.DeployConfig, "domain"),
		CommonName: getDeployString(q.option.DeployConfig, "domain"),
		Pri:        q.option.Certificate.PrivateKey,
		Ca:         q.option.Certificate.Certificate,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("json.Marshal failed: %w", err)
	}

	res, err := q.req(qiniuGateway+path, http.MethodPost, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("req failed: %w", err)
	}
	resp := &uploadCertResp{}
	err = json.Unmarshal(res, resp)
	if err != nil {
		return "", fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	return resp.CertID, nil
}

type modifyDomainCertReq struct {
	CertID      string `json:"certId"`
	ForceHttps  bool   `json:"forceHttps"`
	Http2Enable bool   `json:"http2Enable"`
}

func (q *qiuniu) modifyDomainCert(certId string) error {
	path := fmt.Sprintf("/domain/%s/httpsconf", getDeployString(q.option.DeployConfig, "domain"))

	body := &modifyDomainCertReq{
		CertID:      certId,
		ForceHttps:  true,
		Http2Enable: true,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json.Marshal failed: %w", err)
	}

	_, err = q.req(qiniuGateway+path, http.MethodPut, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("req failed: %w", err)
	}

	return nil
}

func (q *qiuniu) req(url, method string, body io.Reader) ([]byte, error) {
	req := xhttp.BuildReq(url, method, body, map[string]string{
		"Content-Type": "application/json",
	})

	if err := q.credentials.AddToken(auth.TokenQBox, req); err != nil {
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
