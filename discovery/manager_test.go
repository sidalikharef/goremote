package discovery

import (
	"testing"
)

func TestDeduplication(t *testing.T) {
	m := NewManager()
	res := []Result{
		{ID: "id1", Address: "192.168.1.10", Name: "TV1"},
		{ID: "id1", Address: "192.168.1.10", Name: "TV1"},
		{ID: "id2", Address: "192.168.1.11", Name: "TV2"},
	}
	out := m.Deduplicate(res)
	if len(out) != 2 {
		t.Fatalf("want 2 got %d", len(out))
	}
}

func TestMalformedResponse(t *testing.T) {
	m := NewManager()
	res := []Result{
		{ID: "", Address: "", Name: "bad"},
		{ID: "id3", Address: "not-an-ip", Name: "bad2"},
	}
	out := m.Deduplicate(res)
	if len(out) != 0 {
		t.Fatalf("want 0 got %d", len(out))
	}
}

func TestManualAdd(t *testing.T) {
	m := NewManager()
	r, err := m.AddManual("m1", "Manual", "roku", "192.168.1.50")
	if err != nil {
		t.Fatal(err)
	}
	if r.Source != "manual" {
		t.Fatal("bad source")
	}
	if _, err := m.AddManual("m2", "Manual", "roku", "bad-ip"); err == nil {
		t.Fatal("should fail")
	}
}

func TestIPChange(t *testing.T) {
	m := NewManager()
	m.AddManual("id1", "TV", "roku", "192.168.1.10")
	if !m.HandleIPChange("id1", "192.168.1.20") {
		t.Fatal("should detect change")
	}
}
