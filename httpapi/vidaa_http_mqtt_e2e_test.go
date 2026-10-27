package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	vidaa "github.com/sidalikharef/goremote/drivers/vidaa"
	"github.com/sidalikharef/goremote/drivers/registry"
	"github.com/sidalikharef/goremote/internal/auth"
)

func TestVidaaHTTPToMQTTEndToEnd(t *testing.T) {
	vidaa.EnableExperimental()
	defer vidaa.DisableExperimental()

	// start mock broker
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:18845").SetClientID("broker-test")
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		t.Skip("no mock broker available")
	}
	defer c.Disconnect(250)
	received := make(chan string, 1)
	topic := "/remoteapp/tv/remote_service/VIDAACLI/actions/sendkey"
	if token := c.Subscribe(topic, 0, func(_ mqtt.Client, msg mqtt.Message) {
		received <- string(msg.Payload())
	}); token.Wait() && token.Error() != nil {
		t.Fatal(token.Error())
	}

	r := registry.New()
	r.RegisterFactory("vidaa", vidaa.Factory)
	as := auth.New(t.TempDir() + "/auth.json")
	s := NewWithAuth(r, as)
	h := SecurityHeaders(s.Routes())
	ts := httptest.NewServer(h)
	defer ts.Close()

	// register device
	body, _ := json.Marshal(map[string]any{"id": "v1", "name": "Vidaa", "address": "127.0.0.1", "enable": true})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/devices/vidaa/register", bytes.NewReader(body))
	resp, _ := http.DefaultClient.Do(req)
	if resp != nil {
		resp.Body.Close()
	}

	// try keypress - expect fail without auth/pair
	req2, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/keypress/Up?device=v1", nil)
	resp2, _ := http.DefaultClient.Do(req2)
	if resp2 != nil {
		resp2.Body.Close()
		if resp2.StatusCode == http.StatusOK {
			t.Fatal("expected failure without auth")
		}
	}

	select {
	case <-received:
		t.Fatal("unexpected mqtt message")
	case <-time.After(200 * time.Millisecond):
		// ok
	}
}
