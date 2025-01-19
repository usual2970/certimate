package notify

import (
	"context"
	"fmt"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/domain/dtos"
)

const (
	notifyTestTitle = "测试通知"
	notifyTestBody  = "欢迎使用 Certimate ，这是一条测试通知。"
)

type settingsRepository interface {
	GetByName(ctx context.Context, name string) (*domain.Settings, error)
}

type NotifyService struct {
	settingsRepo settingsRepository
}

func NewNotifyService(settingsRepo settingsRepository) *NotifyService {
	return &NotifyService{
		settingsRepo: settingsRepo,
	}
}

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
