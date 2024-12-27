package notify

import (
	"context"
	"fmt"

	"github.com/usual2970/certimate/internal/domain"
)

const (
	notifyTestTitle = "测试通知"
	notifyTestBody  = "欢迎使用 Certimate ，这是一条测试通知。"
)

type SettingRepository interface {
	GetByName(ctx context.Context, name string) (*domain.Settings, error)
}

type NotifyService struct {
	settingRepo SettingRepository
}

func NewNotifyService(settingRepo SettingRepository) *NotifyService {
	return &NotifyService{
		settingRepo: settingRepo,
	}
}

func (n *NotifyService) Test(ctx context.Context, req *domain.NotifyTestPushReq) error {
	setting, err := n.settingRepo.GetByName(ctx, "notifyChannels")
	if err != nil {
		return fmt.Errorf("failed to get notify channels settings: %w", err)
	}

	channelConfig, err := setting.GetChannelContent(req.Channel)
	if err != nil {
		return fmt.Errorf("failed to get notify channel \"%s\" config: %w", req.Channel, err)
	}

	return SendToChannel(notifyTestTitle, notifyTestBody, req.Channel, channelConfig)
}
