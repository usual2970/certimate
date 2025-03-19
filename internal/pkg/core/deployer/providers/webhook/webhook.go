package webhook

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/utils/certutil"
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

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	certX509, err := certutil.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to parse x509")
	}

	var webhookData interface{}
	err = json.Unmarshal([]byte(d.config.WebhookData), &webhookData)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to unmarshall webhook data")
	}

	replaceJsonValueRecursively(webhookData, "${DOMAIN}", certX509.Subject.CommonName)
	replaceJsonValueRecursively(webhookData, "${DOMAINS}", strings.Join(certX509.DNSNames, ";"))
	replaceJsonValueRecursively(webhookData, "${SUBJECT_ALT_NAMES}", strings.Join(certX509.DNSNames, ";"))
	replaceJsonValueRecursively(webhookData, "${CERTIFICATE}", certPem)
	replaceJsonValueRecursively(webhookData, "${PRIVATE_KEY}", privkeyPem)

	resp, err := d.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(webhookData).
		Post(d.config.WebhookUrl)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to send webhook request")
	} else if resp.StatusCode() != 200 {
		return nil, xerrors.Errorf("unexpected webhook response status code: %d", resp.StatusCode())
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
