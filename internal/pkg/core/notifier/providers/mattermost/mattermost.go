package mattermost

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type NotifierConfig struct {
	// 服务地址。
	ServerUrl string `json:"serverUrl"`
	// 用户名。
	Username string `json:"username"`
	// 密码。
	Password string `json:"password"`
	// 频道 ID。
	ChannelId string `json:"channelId"`
}

type NotifierProvider struct {
	config     *NotifierConfig
	logger     *slog.Logger
	httpClient *resty.Client
}

var _ notifier.Notifier = (*NotifierProvider)(nil)

func NewNotifier(config *NotifierConfig) (*NotifierProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client := resty.New()

	return &NotifierProvider{
		config:     config,
		logger:     slog.Default(),
		httpClient: client,
	}, nil
}

func (n *NotifierProvider) WithLogger(logger *slog.Logger) notifier.Notifier {
	if logger == nil {
		n.logger = slog.Default()
	} else {
		n.logger = logger
	}
	return n
}

func (n *NotifierProvider) Notify(ctx context.Context, subject string, message string) (res *notifier.NotifyResult, err error) {
	serverUrl := strings.TrimRight(n.config.ServerUrl, "/")

	// REF: https://developers.mattermost.com/api-documentation/#/operations/Login
	loginReq := n.httpClient.R().
		SetHeader("Content-Type", "application/json").
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
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+loginResp.Header().Get("Token")).
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

	return &notifier.NotifyResult{}, nil
}
