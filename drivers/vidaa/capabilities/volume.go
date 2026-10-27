package capabilities

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/sidalikharef/goremote/drivers/vidaa/protocol"
	"github.com/sidalikharef/goremote/drivers/vidaa/transport"
)

type VolumeService struct {
	tr       *transport.Client
	clientID string
	auth     bool
	deviceID string
}

func NewVolume(tr *transport.Client, clientID, deviceID string, auth bool) *VolumeService {
	return &VolumeService{tr: tr, clientID: clientID, deviceID: deviceID, auth: auth}
}

func (v *VolumeService) GetVolume(ctx context.Context) (int, error) {
	if v == nil || v.tr == nil {
		return 0, errors.New("nil transport")
	}
	if !v.auth {
		return 0, errors.New("not authenticated")
	}
	topic, err := protocol.TopicTVActions("platform_service", v.clientID, "getvolume")
	if err != nil {
		return 0, err
	}
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := v.tr.Publish(ctx2, topic, 0, false, "{}"); err != nil {
		return 0, err
	}
	return 0, errors.New("no response")
}

func (v *VolumeService) ChangeVolume(ctx context.Context, vol int) error {
	if v == nil || v.tr == nil {
		return errors.New("nil transport")
	}
	if !v.auth {
		return errors.New("not authenticated")
	}
	if vol < 0 || vol > 100 {
		return errors.New("invalid volume")
	}
	topic, err := protocol.TopicTVActions("platform_service", v.clientID, "changevolume")
	if err != nil {
		return err
	}
	payload, _ := json.Marshal(map[string]int{"volume": vol})
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return v.tr.Publish(ctx2, topic, 0, false, payload)
}
