package capabilities

import (
	"context"
	"testing"

	"github.com/sidalikharef/goremote/drivers/vidaa/transport"
)

func TestVolumeAuth(t *testing.T) {
	cl, _ := transport.New(transport.Config{Broker: "tcp://127.0.0.1:18843", ClientID: "c1"})
	v := NewVolume(cl, "c1", "d1", false)
	if _, err := v.GetVolume(context.Background()); err == nil {
		t.Fatal("expected")
	}
	if err := v.ChangeVolume(context.Background(), 50); err == nil {
		t.Fatal("expected")
	}
	if err := v.ChangeVolume(context.Background(), 200); err == nil {
		t.Fatal("expected")
	}
}

func TestSourcesAppsAuth(t *testing.T) {
	cl, _ := transport.New(transport.Config{Broker: "tcp://127.0.0.1:18844", ClientID: "c1"})
	s := New(cl, "c1", "d1", false)
	if err := s.ListSources(context.Background()); err == nil {
		t.Fatal("expected")
	}
	if err := s.ChangeSource(context.Background(), ""); err == nil {
		t.Fatal("expected")
	}
	if err := s.LaunchApp(context.Background(), ""); err == nil {
		t.Fatal("expected")
	}
}
