package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
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
	httpClient *httpclient.Client
}

var _ deployer.Deployer = (*WebhookDeployer)(nil)

func New(config *WebhookDeployerConfig) (*WebhookDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *WebhookDeployerConfig, logger logger.Logger) (*WebhookDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client := httpclient.NewClient(httpclient.WithHTTPTimeout(30 * time.Second))

	return &WebhookDeployer{
		config:     config,
		logger:     logger,
		httpClient: client,
	}, nil
}

type webhookData struct {
	SubjectAltNames string            `json:"subjectAltNames"`
	Certificate     string            `json:"certificate"`
	PrivateKey      string            `json:"privateKey"`
	Variables       map[string]string `json:"variables"`
}

func (d *WebhookDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	certX509, err := x509.ParseCertificateFromPEM(certPem)
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

	reqBody, _ := json.Marshal(&webhookData)
	resp, err := d.httpClient.Post(d.config.WebhookUrl, bytes.NewReader(reqBody), map[string][]string{"Content-Type": {"application/json"}})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to send webhook request")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to read response body")
	}

	d.logger.Logt("Webhook Response", string(respBody))

	return &deployer.DeployResult{
		ExtendedData: map[string]any{
			"responseText": string(respBody),
		},
	}, nil
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
