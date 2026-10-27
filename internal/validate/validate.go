package validate

import (
	"fmt"
	"net"
	"strings"
)

// IP validates an IP address.
func IP(s string) error {
	if net.ParseIP(s) == nil {
		return fmt.Errorf("invalid ip")
	}
	return nil
}

// KeyName validates a key name against allowlist.
func KeyName(s string) error {
	if s == "" {
		return fmt.Errorf("empty key")
	}
	allowed := map[string]bool{
		"Home": true, "Back": true, "Up": true, "Down": true, "Left": true, "Right": true,
		"Select": true, "Play": true, "Rev": true, "Fwd": true, "Stop": true, "Pause": true,
		"VolumeUp": true, "VolumeDown": true, "VolumeMute": true, "Power": true,
		"PowerOn": true, "PowerOff": true, "Enter": true, "Info": true, "Search": true,
	}
	// Also allow Lit_ prefixed keys
	if strings.HasPrefix(s, "Lit_") {
		if len(s) > 4 && len(s) < 64 {
			return nil
		}
		return fmt.Errorf("invalid lit key")
	}
	if allowed[s] {
		return nil
	}
	return fmt.Errorf("invalid key")
}

// AppID validates app ID.
func AppID(s string) error {
	if s == "" {
		return fmt.Errorf("empty app id")
	}
	if len(s) > 128 {
		return fmt.Errorf("invalid app id")
	}
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' && r != '_' {
			return fmt.Errorf("invalid app id")
		}
	}
	return nil
}
