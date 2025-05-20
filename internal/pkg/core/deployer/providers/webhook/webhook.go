package webhook

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

type DeployerConfig struct {
	// Webhook URL。
	WebhookUrl string `json:"webhookUrl"`
	// Webhook 回调数据（application/json 或 application/x-www-form-urlencoded 格式）。
	WebhookData string `json:"webhookData,omitempty"`
	// 请求谓词。
	// 零值时默认为 "POST"。
	Method string `json:"method,omitempty"`
	// 请求标头。
	Headers map[string]string `json:"headers,omitempty"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
}

type DeployerProvider struct {
	config     *DeployerConfig
	logger     *slog.Logger
	httpClient *resty.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second)
	if config.AllowInsecureConnections {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return &DeployerProvider{
		config:     config,
		logger:     slog.Default(),
		httpClient: client,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse x509: %w", err)
	}

	// 提取服务器证书和中间证书
	serverCertPEM, intermediaCertPEM, err := certutil.ExtractCertificatesFromPEM(certPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to extract certs: %w", err)
	}

	// 处理 Webhook URL
	webhookUrl, err := url.Parse(d.config.WebhookUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse webhook url: %w", err)
	} else if webhookUrl.Scheme != "http" && webhookUrl.Scheme != "https" {
		return nil, fmt.Errorf("unsupported webhook url scheme '%s'", webhookUrl.Scheme)
	} else {
		webhookUrl.Path = strings.ReplaceAll(webhookUrl.Path, "${DOMAIN}", url.PathEscape(certX509.Subject.CommonName))
	}

	// 处理 Webhook 请求谓词
	webhookMethod := strings.ToUpper(d.config.Method)
	if webhookMethod == "" {
		webhookMethod = http.MethodPost
	} else if webhookMethod != http.MethodGet &&
		webhookMethod != http.MethodPost &&
		webhookMethod != http.MethodPut &&
		webhookMethod != http.MethodPatch &&
		webhookMethod != http.MethodDelete {
		return nil, fmt.Errorf("unsupported webhook request method '%s'", webhookMethod)
	}

	// 处理 Webhook 请求标头
	webhookHeaders := make(http.Header)
	for k, v := range d.config.Headers {
		webhookHeaders.Set(k, v)
	}

	// 处理 Webhook 请求内容类型
	const CONTENT_TYPE_JSON = "application/json"
	const CONTENT_TYPE_FORM = "application/x-www-form-urlencoded"
	const CONTENT_TYPE_MULTIPART = "multipart/form-data"
	webhookContentType := webhookHeaders.Get("Content-Type")
	if webhookContentType == "" {
		webhookContentType = CONTENT_TYPE_JSON
		webhookHeaders.Set("Content-Type", CONTENT_TYPE_JSON)
	} else if strings.HasPrefix(webhookContentType, CONTENT_TYPE_JSON) &&
		strings.HasPrefix(webhookContentType, CONTENT_TYPE_FORM) &&
		strings.HasPrefix(webhookContentType, CONTENT_TYPE_MULTIPART) {
		return nil, fmt.Errorf("unsupported webhook content type '%s'", webhookContentType)
	}

	// 处理 Webhook 请求数据
	var webhookData interface{}
	if d.config.WebhookData == "" {
		webhookData = map[string]string{
			"name":    strings.Join(certX509.DNSNames, ";"),
			"cert":    certPEM,
			"privkey": privkeyPEM,
		}
	} else {
		err = json.Unmarshal([]byte(d.config.WebhookData), &webhookData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal webhook data: %w", err)
		}

		replaceJsonValueRecursively(webhookData, "${DOMAIN}", certX509.Subject.CommonName)
		replaceJsonValueRecursively(webhookData, "${DOMAINS}", strings.Join(certX509.DNSNames, ";"))
		replaceJsonValueRecursively(webhookData, "${CERTIFICATE}", certPEM)
		replaceJsonValueRecursively(webhookData, "${SERVER_CERTIFICATE}", serverCertPEM)
		replaceJsonValueRecursively(webhookData, "${INTERMEDIA_CERTIFICATE}", intermediaCertPEM)
		replaceJsonValueRecursively(webhookData, "${PRIVATE_KEY}", privkeyPEM)

		if webhookMethod == http.MethodGet || webhookContentType == CONTENT_TYPE_FORM || webhookContentType == CONTENT_TYPE_MULTIPART {
			temp := make(map[string]string)
			jsonb, err := json.Marshal(webhookData)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal webhook data: %w", err)
			} else if err := json.Unmarshal(jsonb, &temp); err != nil {
				return nil, fmt.Errorf("failed to unmarshal webhook data: %w", err)
			} else {
				webhookData = temp
			}
		}
	}

	// 生成请求
	// 其中 GET 请求需转换为查询参数
	req := d.httpClient.R().SetHeaderMultiValues(webhookHeaders)
	req.URL = webhookUrl.String()
	req.Method = webhookMethod
	if webhookMethod == http.MethodGet {
		req.SetQueryParams(webhookData.(map[string]string))
	} else {
		switch webhookContentType {
		case CONTENT_TYPE_JSON:
			req.SetBody(webhookData)
		case CONTENT_TYPE_FORM:
			req.SetFormData(webhookData.(map[string]string))
		case CONTENT_TYPE_MULTIPART:
			req.SetMultipartFormData(webhookData.(map[string]string))
		}
	}

	// 发送请求
	resp, err := req.Send()
	if err != nil {
		return nil, fmt.Errorf("failed to send webhook request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("unexpected webhook response status code: %d", resp.StatusCode())
	}

	d.logger.Debug("webhook responded", slog.Any("response", resp.String()))

	return &deployer.DeployResult{}, nil
}

func replaceJsonValueRecursively(data interface{}, oldStr, newStr string) interface{} {
	switch v := data.(type) {
	case map[string]any:
		for k, val := range v {
			v[k] = replaceJsonValueRecursively(val, oldStr, newStr)
		}
	case []any:
		for i, val := range v {
			v[i] = replaceJsonValueRecursively(val, oldStr, newStr)
		}
	case string:
		return strings.ReplaceAll(v, oldStr, newStr)
	}
	return data
}
