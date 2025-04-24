package notify

import (
	"context"
	"fmt"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/domain/dtos"
)

// Deprecated: v0.4.x 将废弃
const (
	notifyTestTitle = "测试通知"
	notifyTestBody  = "欢迎使用 Certimate ，这是一条测试通知。"
)

// Deprecated: v0.4.x 将废弃
type settingsRepository interface {
	GetByName(ctx context.Context, name string) (*domain.Settings, error)
}

// Deprecated: v0.4.x 将废弃
type NotifyService struct {
	settingsRepo settingsRepository
}

// Deprecated: v0.4.x 将废弃
func NewNotifyService(settingsRepo settingsRepository) *NotifyService {
	return &NotifyService{
		settingsRepo: settingsRepo,
	}
}

// Deprecated: v0.4.x 将废弃
func (n *NotifyService) Test(ctx context.Context, req *dtos.NotifyTestPushReq) error {
	settings, err := n.settingsRepo.GetByName(ctx, "notifyChannels")
	if err != nil {
		return fmt.Errorf("failed to get notify channels settings: %w", err)
	}

	channelConfig, err := settings.GetNotifyChannelConfig(string(req.Channel))
	if err != nil {
		return fmt.Errorf("failed to get notify channel \"%s\" config: %w", req.Channel, err)
	}

	return SendToChannel(notifyTestTitle, notifyTestBody, string(req.Channel), channelConfig)
}
