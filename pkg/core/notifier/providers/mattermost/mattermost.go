package mattermost

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/certimate-go/certimate/pkg/core"
)

type NotifierProviderConfig struct {
	// Mattermost 服务地址。
	ServerUrl string `json:"serverUrl"`
	// Mattermost 用户名。
	Username string `json:"username"`
	// Mattermost 密码。
	Password string `json:"password"`
	// Mattermost 频道 ID。
	ChannelId string `json:"channelId"`
}

type NotifierProvider struct {
	config     *NotifierProviderConfig
	logger     *slog.Logger
	httpClient *resty.Client
}

var _ core.Notifier = (*NotifierProvider)(nil)

func NewNotifierProvider(config *NotifierProviderConfig) (*NotifierProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the notifier provider is nil")
	}

	client := resty.New()

	return &NotifierProvider{
		config:     config,
		logger:     slog.Default(),
		httpClient: client,
	}, nil
}

func (n *NotifierProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		n.logger = slog.New(slog.DiscardHandler)
	} else {
		n.logger = logger
	}
}

func (n *NotifierProvider) Notify(ctx context.Context, subject string, message string) (*core.NotifyResult, error) {
	serverUrl := strings.TrimRight(n.config.ServerUrl, "/")

	// REF: https://developers.mattermost.com/api-documentation/#/operations/Login
	loginReq := n.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "certimate").
		SetBody(map[string]any{
			"login_id": n.config.Username,
			"password": n.config.Password,
		})
	loginResp, err := loginReq.Post(fmt.Sprintf("%s/api/v4/users/login", serverUrl))
	if err != nil {
		return nil, fmt.Errorf("mattermost api error: failed to send request: %w", err)
	} else if loginResp.IsError() {
		return nil, fmt.Errorf("mattermost api error: unexpected status code: %d, resp: %s", loginResp.StatusCode(), loginResp.String())
	} else if loginResp.Header().Get("Token") == "" {
		return nil, fmt.Errorf("mattermost api error: received empty login token")
	}

	// REF: https://developers.mattermost.com/api-documentation/#/operations/CreatePost
	postReq := n.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+loginResp.Header().Get("Token")).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "certimate").
		SetBody(map[string]any{
			"channel_id": n.config.ChannelId,
			"props": map[string]interface{}{
				"attachments": []map[string]interface{}{
					{
						"title": subject,
						"text":  message,
					},
				},
			},
		})
	postResp, err := postReq.Post(fmt.Sprintf("%s/api/v4/posts", serverUrl))
	if err != nil {
		return nil, fmt.Errorf("mattermost api error: failed to send request: %w", err)
	} else if postResp.IsError() {
		return nil, fmt.Errorf("mattermost api error: unexpected status code: %d, resp: %s", postResp.StatusCode(), postResp.String())
	}

	return &core.NotifyResult{}, nil
}
