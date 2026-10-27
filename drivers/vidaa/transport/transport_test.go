package transport

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func TestNewRequiresBroker(t *testing.T) {
	_, err := New(Config{ClientID: "x"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewRequiresClientID(t *testing.T) {
	_, err := New(Config{Broker: "tcp://localhost:1883"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewTLSConfigRejectsInsecure(t *testing.T) {
	_, err := NewTLSConfig(nil, tls.Certificate{}, true)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPublishSubscribeMock(t *testing.T) {
	// Use a local TCP listener with paho mock? Or just test API
	// Create broker
	broker := "tcp://127.0.0.1:18830"
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("test")
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		t.Skip("no broker")
	}
	defer c.Disconnect(250)
	ch := make(chan []byte, 1)
	if token := c.Subscribe("test/topic", 0, func(client mqtt.Client, msg mqtt.Message) {
		ch <- msg.Payload()
	}); token.Wait() && token.Error() != nil {
		t.Fatal(token.Error())
	}
	if token := c.Publish("test/topic", 0, false, "hello"); token.Wait() && token.Error() != nil {
		t.Fatal(token.Error())
	}
	select {
	case b := <-ch:
		if string(b) != "hello" {
			t.Fatalf("bad: %s", string(b))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout")
	}
}

func TestClientAPI(t *testing.T) {
	cl, err := New(Config{Broker: "tcp://127.0.0.1:18831", ClientID: "c1", Timeout: 2 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = ctx
	_ = cl.IsConnected()
	cl.Disconnect(100)
}
