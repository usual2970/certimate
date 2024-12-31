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
	Url string `json:"url"`
	// Webhook 变量字典。
	Variables map[string]string `json:"variables,omitempty"`
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

	// TODO: 自定义回调数据
	reqBody, _ := json.Marshal(&webhookData{
		SubjectAltNames: strings.Join(certX509.DNSNames, ","),
		Certificate:     certPem,
		PrivateKey:      privkeyPem,
		Variables:       d.config.Variables,
	})
	resp, err := d.httpClient.Post(d.config.Url, bytes.NewReader(reqBody), map[string][]string{"Content-Type": {"application/json"}})
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
