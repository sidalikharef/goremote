package capabilities

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/sidalikharef/goremote/drivers/vidaa/protocol"
	"github.com/sidalikharef/goremote/drivers/vidaa/transport"
)

type Service struct {
	tr       *transport.Client
	clientID string
	auth     bool
	deviceID string
}

func New(tr *transport.Client, clientID, deviceID string, auth bool) *Service {
	return &Service{tr: tr, clientID: clientID, deviceID: deviceID, auth: auth}
}

func (s *Service) ListSources(ctx context.Context) error {
	if s == nil || s.tr == nil {
		return errors.New("nil transport")
	}
	if !s.auth {
		return errors.New("not authenticated")
	}
	topic, err := protocol.TopicTVActions("ui_service", s.clientID, "sourcelist")
	if err != nil {
		return err
	}
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return s.tr.Publish(ctx2, topic, 0, false, "{}")
}

func (s *Service) ChangeSource(ctx context.Context, sourceID string) error {
	if s == nil || s.tr == nil {
		return errors.New("nil transport")
	}
	if !s.auth {
		return errors.New("not authenticated")
	}
	if sourceID == "" {
		return errors.New("invalid source")
	}
	topic, err := protocol.TopicTVActions("ui_service", s.clientID, "changesource")
	if err != nil {
		return err
	}
	payload, _ := json.Marshal(map[string]string{"sourceId": sourceID})
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return s.tr.Publish(ctx2, topic, 0, false, payload)
}

func (s *Service) ListApps(ctx context.Context) error {
	if s == nil || s.tr == nil {
		return errors.New("nil transport")
	}
	if !s.auth {
		return errors.New("not authenticated")
	}
	topic, err := protocol.TopicTVActions("ui_service", s.clientID, "applist")
	if err != nil {
		return err
	}
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return s.tr.Publish(ctx2, topic, 0, false, "{}")
}

func (s *Service) LaunchApp(ctx context.Context, appID string) error {
	if s == nil || s.tr == nil {
		return errors.New("nil transport")
	}
	if !s.auth {
		return errors.New("not authenticated")
	}
	if appID == "" {
		return errors.New("invalid app")
	}
	topic, err := protocol.TopicTVActions("ui_service", s.clientID, "launchapp")
	if err != nil {
		return err
	}
	payload, _ := json.Marshal(map[string]string{"appId": appID})
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return s.tr.Publish(ctx2, topic, 0, false, payload)
}

func (s *Service) GetTVState(ctx context.Context) error {
	if s == nil || s.tr == nil {
		return errors.New("nil transport")
	}
	if !s.auth {
		return errors.New("not authenticated")
	}
	topic, err := protocol.TopicTVActions("ui_service", s.clientID, "gettvstate")
	if err != nil {
		return err
	}
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return s.tr.Publish(ctx2, topic, 0, false, "{}")
}
