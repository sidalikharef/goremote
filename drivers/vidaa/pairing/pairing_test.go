package pairing

import (
	"context"
	"testing"
	"time"
)

func TestNewSession(t *testing.T) {
	s, err := NewSession("dev1", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if s.ID == "" {
		t.Fatal("bad")
	}
	if _, ok := Get(s.ID); !ok {
		t.Fatal("not found")
	}
}

func TestSubmitPIN(t *testing.T) {
	s, _ := NewSession("dev1", time.Minute)
	if err := s.SubmitPIN(context.Background(), "1234"); err != nil {
		t.Fatal(err)
	}
	if err := s.SubmitPIN(context.Background(), "123"); err == nil {
		t.Fatal("expected")
	}
	if err := s.SubmitPIN(context.Background(), "12a4"); err == nil {
		t.Fatal("expected")
	}
}

func TestCancel(t *testing.T) {
	s, _ := NewSession("dev1", time.Minute)
	s.Cancel()
	if s.State != "cancelled" {
		t.Fatal("bad")
	}
}
