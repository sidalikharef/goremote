package registry

import (
	"context"
	"sync"
	"testing"

	"github.com/sidalikharef/goremote/tv"
)

type mockDriver struct {
	name string
}

func (m *mockDriver) Name() string                            { return m.name }
func (m *mockDriver) Key(ctx context.Context, k string) error { return nil }
func (m *mockDriver) Apps(ctx context.Context) ([]tv.App, error) {
	return []tv.App{{ID: "1", Name: "Test"}}, nil
}
func (m *mockDriver) Launch(ctx context.Context, id string) error { return nil }

func TestRegistryConcurrency(t *testing.T) {
	r := New()
	r.RegisterFactory("mock", func(opts map[string]string) (tv.Driver, error) {
		return &mockDriver{name: "mock"}, nil
	})
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			did := "dev" + string(rune('0'+i%10)) + string(rune('a'+i%26))
			if _, err := r.RegisterDevice(did, "d", "mock", "1.2.3.4"); err == nil {
				if d, ok := r.GetDevice(did); ok {
					d.SetDriver(&mockDriver{name: "mock"})
					_ = d.Call(context.Background(), func(drv tv.Driver) error {
						return drv.Key(context.Background(), "Home")
					})
					_ = d.GetState()
				}
			}
			_ = r.ListDevices()
		}(i)
	}
	wg.Wait()
}

func TestRegistryStates(t *testing.T) {
	r := New()
	d, err := r.RegisterDevice("id1", "n1", "roku", "192.168.1.100")
	if err != nil {
		t.Fatal(err)
	}
	if d.GetState() != tv.StateDisconnected {
		t.Fatal("bad state")
	}
	d.SetDriver(&mockDriver{name: "roku"})
	if d.GetState() != tv.StateConnected {
		t.Fatal("bad state")
	}
	d.SetState(tv.StatePairingRequired)
	if d.GetState() != tv.StatePairingRequired {
		t.Fatal("bad state")
	}
}
