package protocol

import "testing"

func TestTopicTVActions(t *testing.T) {
	topic, err := TopicTVActions("ui_service", "AA:BB:CC:DD:EE:FF", "vidaa_app_connect")
	if err != nil {
		t.Fatal(err)
	}
	if topic != "/remoteapp/tv/ui_service/AA:BB:CC:DD:EE:FF/actions/vidaa_app_connect" {
		t.Fatal(topic)
	}
}

func TestTopicInvalid(t *testing.T) {
	if _, err := TopicTVActions("bad", "id", "x"); err == nil {
		t.Fatal("expected")
	}
	if _, err := TopicTVActions("ui_service", "../x", "x"); err == nil {
		t.Fatal("expected")
	}
}

func TestAuthCodeValidate(t *testing.T) {
	if err := (AuthenticationCode{AuthNum: "1234"}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (AuthenticationCode{AuthNum: "123"}).Validate(); err == nil {
		t.Fatal("expected")
	}
}

func TestParseToken(t *testing.T) {
	tok, err := ParseTokenIssuance([]byte(`{"access_token":"abc"}`))
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken != "abc" {
		t.Fatal("bad")
	}
	if _, err := ParseTokenIssuance([]byte(`{}`)); err == nil {
		t.Fatal("expected")
	}
}

func TestFSM(t *testing.T) {
	f := NewFSM()
	if !f.Transition(StateConnecting) {
		t.Fatal("bad")
	}
	if !f.Transition(StateTLSEstablished) {
		t.Fatal("bad")
	}
	if !f.Transition(StateAuthRequired) {
		t.Fatal("bad")
	}
	if !f.Transition(StatePairingPending) {
		t.Fatal("bad")
	}
	if !f.Transition(StateAuthenticated) {
		t.Fatal("bad")
	}
}
