package notify

import (
	"context"
	"fmt"

	"certimate/internal/domain"
)

const (
	notifyTestTitle = "测试通知"
	notifyTestBody  = "欢迎使用 Certimate ，这是一条测试通知。"
)

type SettingRepository interface {
	GetByName(ctx context.Context, name string) (*domain.Setting, error)
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
		return fmt.Errorf("get notify channels setting failed: %w", err)
	}

	conf, err := setting.GetChannelContent(req.Channel)
	if err != nil {
		return fmt.Errorf("get notify channel %s config failed: %w", req.Channel, err)
	}

	return SendTest(&sendTestParam{
		Title:   notifyTestTitle,
		Content: notifyTestBody,
		Channel: req.Channel,
		Conf:    conf,
	})
}
