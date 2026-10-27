package remote

import (
	"context"
	"net"
	"testing"

	"github.com/sidalikharef/goremote/drivers/vidaa/transport"
)

func startMock(t *testing.T) func() {
	ln, err := net.Listen("tcp", "127.0.0.1:18840")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			c.Close()
		}
	}()
	return func() { ln.Close() }
}

func TestUnauth(t *testing.T) {
	cl, _ := transport.New(transport.Config{Broker: "tcp://127.0.0.1:18841", ClientID: "c1"})
	r := New(cl, "c1", "d1", false)
	if err := r.SendKey(context.Background(), "Up"); err == nil {
		t.Fatal("expected")
	}
}

func TestUnknownKey(t *testing.T) {
	cl, _ := transport.New(transport.Config{Broker: "tcp://127.0.0.1:18842", ClientID: "c1"})
	r := New(cl, "c1", "d1", true)
	if err := r.SendKey(context.Background(), "NOPE"); err == nil {
		t.Fatal("expected")
	}
}

func TestMap(t *testing.T) {
	if v, ok := MapKey("Up"); !ok || v != "KEY_UP" {
		t.Fatal("bad")
	}
	if v, ok := MapKey("Select"); !ok || v != "KEY_OK" {
		t.Fatal("bad")
	}
	if v, ok := MapKey("Back"); !ok || v != "KEY_RETURNS" {
		t.Fatal("bad")
	}
}
