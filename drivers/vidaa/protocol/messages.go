package protocol

import (
	"encoding/json"
	"errors"
)

type VidaaAppConnect struct {
	AppVersion    int    `json:"app_version"`
	ConnectResult int    `json:"connect_result"`
	DeviceType    string `json:"device_type"`
}

func (v VidaaAppConnect) Validate() error {
	if v.AppVersion < 0 {
		return errors.New("bad version")
	}
	if v.DeviceType == "" {
		return errors.New("device type required")
	}
	return nil
}

type AuthenticationCode struct {
	AuthNum string `json:"authNum"`
}

func (a AuthenticationCode) Validate() error {
	if len(a.AuthNum) != 4 {
		return errors.New("auth code must be 4 digits")
	}
	for _, c := range a.AuthNum {
		if c < '0' || c > '9' {
			return errors.New("auth code must be digits")
		}
	}
	return nil
}

type TokenIssuance struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Token        string `json:"token,omitempty"`
}

func ParseTokenIssuance(b []byte) (*TokenIssuance, error) {
	var t TokenIssuance
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}
	if t.AccessToken == "" && t.Token == "" && t.RefreshToken == "" {
		return nil, errors.New("empty token response")
	}
	return &t, nil
}
