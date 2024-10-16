package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	xhttp "certimate/internal/utils/http"
)

type webhookAccess struct {
	Url string `json:"url"`
}

type hookData struct {
	Domain      string            `json:"domain"`
	Certificate string            `json:"certificate"`
	PrivateKey  string            `json:"privateKey"`
	Variables   map[string]string `json:"variables"`
}

type webhook struct {
	option *DeployerOption
	infos  []string
}

func NewWebhook(option *DeployerOption) (Deployer, error) {
	return &webhook{
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (a *webhook) GetID() string {
	return fmt.Sprintf("%s-%s", a.option.AceessRecord.GetString("name"), a.option.AceessRecord.Id)
}

func (w *webhook) GetInfo() []string {
	return w.infos
}

func (w *webhook) Deploy(ctx context.Context) error {
	access := &webhookAccess{}
	if err := json.Unmarshal([]byte(w.option.Access), access); err != nil {
		return fmt.Errorf("failed to parse hook access config: %w", err)
	}

	data := &hookData{
		Domain:      w.option.Domain,
		Certificate: w.option.Certificate.Certificate,
		PrivateKey:  w.option.Certificate.PrivateKey,
		Variables:   getDeployVariables(w.option.DeployConfig),
	}

	body, _ := json.Marshal(data)

	resp, err := xhttp.Req(access.Url, http.MethodPost, bytes.NewReader(body), map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return fmt.Errorf("failed to send hook request: %w", err)
	}

	w.infos = append(w.infos, toStr("webhook response", string(resp)))

	return nil
}
