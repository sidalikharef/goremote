package auth

import (
	"path/filepath"
	"testing"
)

func TestEnrollLoginLogout(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "auth.json")
	s := New(p)
	if s.HasUsers() {
		t.Fatal("should have no users")
	}
	token, err := s.GenerateBootstrapToken()
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Enroll("admin", "longenoughpass", token); err != nil {
		t.Fatal(err)
	}
	if s.HasUsers() == false {
		t.Fatal("should have users")
	}
	if !s.Authenticate("admin", "longenoughpass") {
		t.Fatal("auth failed")
	}
	if s.Authenticate("admin", "wrong") {
		t.Fatal("bad auth")
	}
	tok, err := s.CreateSession()
	if err != nil {
		t.Fatal(err)
	}
	if !s.ValidateSession(tok) {
		t.Fatal("bad session")
	}
	s.RevokeSession(tok)
	if s.ValidateSession(tok) {
		t.Fatal("revoked")
	}
}

func TestSecondEnrollRejected(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "auth.json")
	s := New(p)
	token, _ := s.GenerateBootstrapToken()
	if err := s.Enroll("admin", "longenoughpass", token); err != nil {
		t.Fatal(err)
	}
	if err := s.Enroll("admin2", "longenoughpass", "bad"); err == nil {
		t.Fatal("should fail")
	}
}
