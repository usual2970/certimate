package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"certimate/internal/domain"
	xhttp "certimate/internal/utils/http"
)

type WebhookDeployer struct {
	option *DeployerOption
	infos  []string
}

func NewWebhookDeployer(option *DeployerOption) (Deployer, error) {
	return &WebhookDeployer{
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (d *WebhookDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *WebhookDeployer) GetInfo() []string {
	return d.infos
}

type webhookData struct {
	Domain      string            `json:"domain"`
	Certificate string            `json:"certificate"`
	PrivateKey  string            `json:"privateKey"`
	Variables   map[string]string `json:"variables"`
}

func (d *WebhookDeployer) Deploy(ctx context.Context) error {
	access := &domain.WebhookAccess{}
	if err := json.Unmarshal([]byte(d.option.Access), access); err != nil {
		return fmt.Errorf("failed to parse hook access config: %w", err)
	}

	data := &webhookData{
		Domain:      d.option.Domain,
		Certificate: d.option.Certificate.Certificate,
		PrivateKey:  d.option.Certificate.PrivateKey,
		Variables:   getDeployVariables(d.option.DeployConfig),
	}

	body, _ := json.Marshal(data)

	resp, err := xhttp.Req(access.Url, http.MethodPost, bytes.NewReader(body), map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return fmt.Errorf("failed to send hook request: %w", err)
	}

	d.infos = append(d.infos, toStr("webhook response", string(resp)))

	return nil
}
