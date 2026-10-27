package credentials

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStore(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "creds.json")
	st := NewStore(p)
	c := Creds{DeviceID: "d1", AccessToken: "tok"}
	if err := st.Save(c); err != nil {
		t.Fatal(err)
	}
	info, _ := os.Stat(p)
	if info.Mode().Perm() != 0600 {
		t.Fatalf("bad perms %v", info.Mode().Perm())
	}
	l, err := st.Load("d1")
	if err != nil {
		t.Fatal(err)
	}
	if l.AccessToken != "tok" {
		t.Fatal("bad")
	}
	if _, err := st.Load("d2"); err == nil {
		t.Fatal("expected")
	}
	if err := st.Delete(); err != nil {
		t.Fatal(err)
	}
}
