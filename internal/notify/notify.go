package notify

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/pkg/core/notifier"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

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

func SendToChannel(subject, message string, channel string, channelConfig map[string]any) error {
	notifier, err := createNotifier(channel, channelConfig)
	if err != nil {
		return err
	}

	_, err = notifier.Notify(context.Background(), subject, message)
	return err
}

func getEnabledNotifiers() ([]notifier.Notifier, error) {
	settings, err := app.GetApp().Dao().FindFirstRecordByFilter("settings", "name='notifyChannels'")
	if err != nil {
		return nil, fmt.Errorf("find notifyChannels error: %w", err)
	}

	rs := make(map[string]map[string]any)
	if err := settings.UnmarshalJSONField("content", &rs); err != nil {
		return nil, fmt.Errorf("unmarshal notifyChannels error: %w", err)
	}

	notifiers := make([]notifier.Notifier, 0)
	for k, v := range rs {
		if !maps.GetValueAsBool(v, "enabled") {
			continue
		}

		notifier, err := createNotifier(k, v)
		if err != nil {
			continue
		}

		notifiers = append(notifiers, notifier)
	}

	return notifiers, nil
}
