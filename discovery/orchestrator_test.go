package discovery

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockP struct {
	name string
	res  []Result
	err  error
}

func (m *mockP) Name() string { return m.name }
func (m *mockP) Discover(ctx context.Context, opts Options) ([]Result, error) {
	return m.res, m.err
}

func TestOrchestrator_MultiProvider(t *testing.T) {
	mgr := NewManager()
	p1 := &mockP{name: "p1", res: []Result{{ID: "r1", Name: "Roku", Address: "192.168.1.10", DriverHint: "roku"}}}
	p2 := &mockP{name: "p2", res: []Result{{ID: "v1", Name: "Hisense", Address: "192.168.1.20", DriverHint: "vidaa"}}}
	o := NewOrchestrator(mgr, p1, p2)
	res, errs := o.Discover(context.Background(), Options{Timeout: time.Second})
	if len(res) != 2 {
		t.Fatalf("want 2 got %d", len(res))
	}
	if len(errs) != 0 {
		t.Fatalf("unexpected errs")
	}
}

func TestOrchestrator_PartialFailure(t *testing.T) {
	mgr := NewManager()
	p1 := &mockP{name: "p1", res: []Result{{ID: "r1", Name: "Roku", Address: "192.168.1.11", DriverHint: "roku"}}}
	p2 := &mockP{name: "p2", res: nil, err: errors.New("boom")}
	o := NewOrchestrator(mgr, p1, p2)
	res, errs := o.Discover(context.Background(), Options{Timeout: time.Second})
	if len(res) != 1 {
		t.Fatalf("want 1 got %d", len(res))
	}
	if len(errs) != 1 {
		t.Fatalf("want 1 err")
	}
}
