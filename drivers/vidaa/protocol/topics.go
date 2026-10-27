package protocol

import (
	"fmt"
	"regexp"
	"strings"
)

var clientIDRe = regexp.MustCompile(`^[a-zA-Z0-9:\-._$]{1,64}$`)

func ValidateClientID(id string) error {
	if id == "" {
		return fmt.Errorf("empty client id")
	}
	if !clientIDRe.MatchString(id) {
		return fmt.Errorf("invalid client id")
	}
	if strings.Contains(id, "..") {
		return fmt.Errorf("invalid client id")
	}
	return nil
}

func TopicTVActions(service, clientID, action string) (string, error) {
	if err := ValidateClientID(clientID); err != nil {
		return "", err
	}
	if service == "" || action == "" {
		return "", fmt.Errorf("missing parts")
	}
	// basic allowlist
	switch service {
	case "remote_service", "ui_service", "platform_service":
	default:
		return "", fmt.Errorf("invalid service")
	}
	return fmt.Sprintf("/remoteapp/tv/%s/%s/actions/%s", service, clientID, action), nil
}

func TopicMobileData(clientID, service, dataType string) (string, error) {
	if err := ValidateClientID(clientID); err != nil {
		return "", err
	}
	if service == "" || dataType == "" {
		return "", fmt.Errorf("missing parts")
	}
	switch service {
	case "ui_service", "platform_service", "remote_service":
	default:
		return "", fmt.Errorf("invalid service")
	}
	return fmt.Sprintf("/remoteapp/mobile/%s/%s/data/%s", clientID, service, dataType), nil
}
