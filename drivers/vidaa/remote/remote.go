package remote

import (
	"context"
	"errors"
	"time"

	"github.com/sidalikharef/goremote/drivers/vidaa/protocol"
	"github.com/sidalikharef/goremote/drivers/vidaa/transport"
)

type Client struct {
	tr       *transport.Client
	clientID string
	auth     bool
	deviceID string
}

func New(tr *transport.Client, clientID, deviceID string, auth bool) *Client {
	return &Client{tr: tr, clientID: clientID, deviceID: deviceID, auth: auth}
}

func (c *Client) SendKey(ctx context.Context, key string) error {
	if c == nil || c.tr == nil {
		return errors.New("nil transport")
	}
	if !c.auth {
		return errors.New("not authenticated")
	}
	if c.deviceID == "" {
		return errors.New("device required")
	}
	vidaaKey, ok := MapKey(key)
	if !ok {
		return errors.New("unknown key")
	}
	topic, err := protocol.TopicTVActions("remote_service", c.clientID, "sendkey")
	if err != nil {
		return err
	}
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := c.tr.Publish(ctx2, topic, 0, false, vidaaKey); err != nil {
		return err
	}
	return nil
}
