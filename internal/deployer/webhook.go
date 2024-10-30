package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/domain"
	xhttp "github.com/usual2970/certimate/internal/utils/http"
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
		return xerrors.Wrap(err, "failed to get access")
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
		return xerrors.Wrap(err, "failed to send webhook request")
	}

	d.infos = append(d.infos, toStr("Webhook Response", string(resp)))

	return nil
}
