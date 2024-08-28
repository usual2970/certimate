package deployer

import (
	"bytes"
	xhttp "certimate/internal/utils/http"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type webhookAccess struct {
	Url string `json:"url"`
}

type hookData struct {
	Domain      string `json:"domain"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
}

type webhook struct {
	option *DeployerOption
}

func NewWebhook(option *DeployerOption) (Deployer, error) {

	return &webhook{
		option: option,
	}, nil
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
	}

	body, _ := json.Marshal(data)

	_, err := xhttp.Req(access.Url, http.MethodPost, bytes.NewReader(body), map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return fmt.Errorf("failed to send hook request: %w", err)
	}

	return nil
}
