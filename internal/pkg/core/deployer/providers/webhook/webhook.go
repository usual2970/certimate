package webhook

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/utils/certs"
)

type WebhookDeployerConfig struct {
	// Webhook URL。
	WebhookUrl string `json:"webhookUrl"`
	// Webhook 回调数据（JSON 格式）。
	WebhookData string `json:"webhookData,omitempty"`
}

type WebhookDeployer struct {
	config     *WebhookDeployerConfig
	logger     logger.Logger
	httpClient *resty.Client
}

var _ deployer.Deployer = (*WebhookDeployer)(nil)

func New(config *WebhookDeployerConfig) (*WebhookDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *WebhookDeployerConfig, logger logger.Logger) (*WebhookDeployer, error) {
	if config == nil {
		panic("config is nil")
	}

	if logger == nil {
		panic("logger is nil")
	}

	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second)

	return &WebhookDeployer{
		config:     config,
		logger:     logger,
		httpClient: client,
	}, nil
}

func (d *WebhookDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	certX509, err := certs.ParseCertificateFromPEM(certPem)
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

	d.logger.Logt("Webhook request sent", resp.String())

	return &deployer.DeployResult{}, nil
}

func replaceJsonValueRecursively(data interface{}, oldStr, newStr string) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		for k, val := range v {
			v[k] = replaceJsonValueRecursively(val, oldStr, newStr)
		}
	case []interface{}:
		for i, val := range v {
			v[i] = replaceJsonValueRecursively(val, oldStr, newStr)
		}
	case string:
		return strings.ReplaceAll(v, oldStr, newStr)
	}
	return data
}
