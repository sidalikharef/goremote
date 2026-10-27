package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sidalikharef/goremote/tv"
)

func TestConfigLoadSaveAtomic(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.json")
	c := &Config{
		Version: 1,
		Devices: []Device{
			{ID: "dev1", Name: "TV1", Driver: "roku", Address: "192.168.1.10", State: tv.StateDisconnected},
		},
		Active: "dev1",
	}
	if err := Save(p, c); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatal(err)
	}
	// check permissions
	info, _ := os.Stat(p)
	if info.Mode().Perm() != 0600 {
		t.Fatalf("bad perms: %v", info.Mode().Perm())
	}
	c2, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(c2.Devices) != 1 || c2.Devices[0].ID != "dev1" {
		t.Fatal("bad load")
	}
}
