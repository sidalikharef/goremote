package vidaa

import (
	"errors"

	"github.com/sidalikharef/goremote/drivers/vidaa/driver"
	"github.com/sidalikharef/goremote/drivers/vidaa/profiles"
	"github.com/sidalikharef/goremote/drivers/vidaa/transport"
	"github.com/sidalikharef/goremote/tv"
)

var experimentalEnabled = false

func EnableExperimental() { experimentalEnabled = true }
func DisableExperimental() { experimentalEnabled = false }

func Factory(opts map[string]string) (tv.Driver, error) {
	if !experimentalEnabled {
		return nil, errors.New("vidaa experimental disabled")
	}
	broker := opts["broker"]
	clientID := opts["clientID"]
	deviceID := opts["device_id"]
	if deviceID == "" {
		deviceID = opts["id"]
	}
	if broker == "" || clientID == "" {
		return nil, errors.New("broker and clientid required")
	}
	tr, err := transport.New(transport.Config{Broker: broker, ClientID: clientID})
	if err != nil {
		return nil, err
	}
	_ = profiles.Get(profiles.FamilyUnknown)
	d := driver.New(tr, clientID, deviceID)
	return d, nil
}
