package webhook

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

type DeployerConfig struct {
	// Webhook URL。
	WebhookUrl string `json:"webhookUrl"`
	// Webhook 回调数据（JSON 格式）。
	WebhookData string `json:"webhookData,omitempty"`
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
	certX509, err := certutil.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse x509: %w", err)
	}

	var webhookData interface{}
	err = json.Unmarshal([]byte(d.config.WebhookData), &webhookData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall webhook data: %w", err)
	}

	replaceJsonValueRecursively(webhookData, "${DOMAIN}", certX509.Subject.CommonName)
	replaceJsonValueRecursively(webhookData, "${DOMAINS}", strings.Join(certX509.DNSNames, ";"))
	replaceJsonValueRecursively(webhookData, "${SUBJECT_ALT_NAMES}", strings.Join(certX509.DNSNames, ";"))
	replaceJsonValueRecursively(webhookData, "${CERTIFICATE}", certPEM)
	replaceJsonValueRecursively(webhookData, "${PRIVATE_KEY}", privkeyPEM)

	resp, err := d.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(webhookData).
		Post(d.config.WebhookUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to send webhook request: %w", err)
	} else if resp.StatusCode() != 200 {
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
