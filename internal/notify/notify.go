package notify

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/repository"
	"github.com/certimate-go/certimate/pkg/core"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

// Deprecated: v0.4.x 将废弃
func SendToAllChannels(subject, message string) error {
	notifiers, err := getEnabledNotifiers()
	if err != nil {
		return err
	}
	if len(notifiers) == 0 {
		return nil
	}

	var eg errgroup.Group
	for _, n := range notifiers {
		if n == nil {
			continue
		}

		eg.Go(func() error {
			_, err := n.Notify(context.Background(), subject, message)
			return err
		})
	}

	err = eg.Wait()
	return err
}

// Deprecated: v0.4.x 将废弃
func SendToChannel(subject, message string, channel string, channelConfig map[string]any) error {
	notifier, err := createNotifierProviderUseGlobalSettings(domain.NotifyChannelType(channel), channelConfig)
	if err != nil {
		return err
	}

	_, err = notifier.Notify(context.Background(), subject, message)
	return err
}

// Deprecated: v0.4.x 将废弃
func getEnabledNotifiers() ([]core.Notifier, error) {
	settingsRepo := repository.NewSettingsRepository()
	settings, err := settingsRepo.GetByName(context.Background(), "notifyChannels")
	if err != nil {
		return nil, fmt.Errorf("find notifyChannels error: %w", err)
	}

	rs := make(map[string]map[string]any)
	if err := json.Unmarshal([]byte(settings.Content), &rs); err != nil {
		return nil, fmt.Errorf("unmarshal notifyChannels error: %w", err)
	}

	notifiers := make([]core.Notifier, 0)
	for k, v := range rs {
		if !xmaps.GetBool(v, "enabled") {
			continue
		}

		notifier, err := createNotifierProviderUseGlobalSettings(domain.NotifyChannelType(k), v)
		if err != nil {
			continue
		}

		notifiers = append(notifiers, notifier)
	}

	return notifiers, nil
}
