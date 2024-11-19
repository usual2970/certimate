package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
	xhttp "github.com/usual2970/certimate/internal/utils/http"
)

type WebhookDeployerConfig struct {
	// Webhook URL。
	Url string `json:"url"`
	// Webhook 变量字典。
	Variables map[string]string `json:"variables,omitempty"`
}

type WebhookDeployer struct {
	config *WebhookDeployerConfig
	logger deployer.Logger
}

var _ deployer.Deployer = (*WebhookDeployer)(nil)

func New(config *WebhookDeployerConfig) (*WebhookDeployer, error) {
	return NewWithLogger(config, deployer.NewNilLogger())
}

func NewWithLogger(config *WebhookDeployerConfig, logger deployer.Logger) (*WebhookDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	return &WebhookDeployer{
		config: config,
		logger: logger,
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

	data := &webhookData{
		SubjectAltNames: strings.Join(certX509.DNSNames, ","),
		Certificate:     certPem,
		PrivateKey:      privkeyPem,
		Variables:       d.config.Variables,
	}
	body, _ := json.Marshal(data)
	resp, err := xhttp.Req(d.config.Url, http.MethodPost, bytes.NewReader(body), map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to send webhook request")
	}

	d.logger.Appendt("Webhook Response", string(resp))

	return &deployer.DeployResult{
		DeploymentData: map[string]any{
			"responseText": string(resp),
		},
	}, nil
}
