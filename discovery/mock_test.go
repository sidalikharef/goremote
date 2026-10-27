package discovery

import (
	"context"
	"testing"
	"time"
)

type mockSvc struct {
	results []Result
}

func (m *mockSvc) Discover(ctx context.Context, opts Options) ([]Result, error) {
	return m.results, nil
}

func TestMockDiscovery(t *testing.T) {
	svc := &mockSvc{results: []Result{
		{ID: "uuid:roku-1", Address: "192.168.1.100", Name: "roku:ecp", DriverHint: "roku", Source: "ssdp"},
	}}
	res, err := svc.Discover(context.Background(), Options{Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatal("bad")
	}
	m := NewManager()
	out := m.Deduplicate(res)
	if len(out) != 1 {
		t.Fatal("bad dedup")
	}
}
