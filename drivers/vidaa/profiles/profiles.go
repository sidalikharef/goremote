package profiles

type FirmwareFamily string

const (
	FamilyUnknown   FirmwareFamily = "unknown"
	FamilyLegacy    FirmwareFamily = "legacy"
	FamilyRemotenow FirmwareFamily = "remotenow"
	FamilyMiddle    FirmwareFamily = "middle"
	FamilyModern    FirmwareFamily = "modern"
)

type Profile struct {
	Family            FirmwareFamily
	MQTTVERSION       byte
	Port              int
	RequiresMTLS      bool
	WildcardSubs      bool
	PinPairing        bool
	CredentialStyle   string
	Notes             string
	Sources           []string
}

func Get(f FirmwareFamily) Profile {
	switch f {
	case FamilyLegacy:
		return Profile{
			Family:          FamilyLegacy,
			MQTTVERSION:     3,
			Port:            36669,
			RequiresMTLS:    true,
			WildcardSubs:    false,
			PinPairing:      false,
			CredentialStyle: "static",
			Notes:           "Legacy instant setup (documented)",
			Sources:         []string{"stevene1919/hisense_vidaa/docs/protocol_architecture.md"},
		}
	case FamilyRemotenow:
		return Profile{
			Family:          FamilyRemotenow,
			MQTTVERSION:     3,
			Port:            36669,
			RequiresMTLS:    true,
			WildcardSubs:    false,
			PinPairing:      true,
			CredentialStyle: "dynamic",
			Sources:         []string{"Lucent-/vidaa-control/docs/PROTOCOL.md", "stevene1919/hisense_vidaa"},
		}
	case FamilyMiddle:
		return Profile{
			Family:          FamilyMiddle,
			MQTTVERSION:     3,
			Port:            36669,
			RequiresMTLS:    true,
			WildcardSubs:    false,
			PinPairing:      true,
			CredentialStyle: "dynamic",
		}
	case FamilyModern:
		return Profile{
			Family:          FamilyModern,
			MQTTVERSION:     3,
			Port:            36669,
			RequiresMTLS:    true,
			WildcardSubs:    false,
			PinPairing:      true,
			CredentialStyle: "dynamic",
		}
	default:
		return Profile{
			Family:       FamilyUnknown,
			MQTTVERSION:  3,
			Port:         36669,
			RequiresMTLS: true,
			WildcardSubs: false,
			PinPairing:   false,
			Notes:        "Unverified - do not assume compatibility",
		}
	}
}
