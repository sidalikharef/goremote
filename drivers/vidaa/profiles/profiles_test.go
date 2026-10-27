package profiles

import "testing"

func TestGetUnknown(t *testing.T) {
	p := Get(FamilyUnknown)
	if p.WildcardSubs {
		t.Fatal("should not allow wildcards")
	}
	if !p.RequiresMTLS {
		t.Fatal("mTLS expected")
	}
}

func TestGetLegacy(t *testing.T) {
	p := Get(FamilyLegacy)
	if p.PinPairing {
		t.Fatal("legacy should not require pin per docs")
	}
}
