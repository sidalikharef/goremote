package tv

import "context"

// App represents an installed application on a TV.
type App struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// State represents the connection/state of a device.
type State string

const (
	StateDisconnected    State = "disconnected"
	StateConnecting      State = "connecting"
	StateConnected       State = "connected"
	StatePairingRequired State = "pairing_required"
)

// Driver is the minimal contract all TV drivers must implement.
type Driver interface {
	// Name returns the driver name (e.g., "roku").
	Name() string
	// Key sends a keypress.
	Key(ctx context.Context, key string) error
	// Apps lists installed apps.
	Apps(ctx context.Context) ([]App, error)
	// Launch launches an app by ID.
	Launch(ctx context.Context, appID string) error
}

// Optional capability interfaces. Drivers may implement any subset.

// Navigator provides directional/navigation controls.
type Navigator interface {
	Up(ctx context.Context) error
	Down(ctx context.Context) error
	Left(ctx context.Context) error
	Right(ctx context.Context) error
	Select(ctx context.Context) error
	Back(ctx context.Context) error
	Home(ctx context.Context) error
}

// VolumeControl provides volume/media controls.
type VolumeControl interface {
	VolumeUp(ctx context.Context) error
	VolumeDown(ctx context.Context) error
	Mute(ctx context.Context) error
}

// Power provides power controls.
type Power interface {
	PowerOn(ctx context.Context) error
	PowerOff(ctx context.Context) error
}

// MediaControl provides playback controls.
type MediaControl interface {
	Play(ctx context.Context) error
	Pause(ctx context.Context) error
	Stop(ctx context.Context) error
	Rewind(ctx context.Context) error
	FastForward(ctx context.Context) error
}

// TextInput provides text entry.
type TextInput interface {
	SendText(ctx context.Context, text string) error
}

// InputSourceSelector provides input/source switching.
type InputSourceSelector interface {
	ListInputs(ctx context.Context) ([]string, error)
	SetInput(ctx context.Context, input string) error
}

// Pairer provides pairing functionality.
type Pairer interface {
	StartPairing(ctx context.Context) (pairingID string, instructions string, err error)
	FinishPairing(ctx context.Context, pairingID string, code string) error
}

// Capabilities describes what a driver supports.
type Capabilities struct {
	Navigation bool `json:"navigation"`
	Volume     bool `json:"volume"`
	Power      bool `json:"power"`
	Media      bool `json:"media"`
	Apps       bool `json:"apps"`
	TextInput  bool `json:"text_input"`
	Inputs     bool `json:"inputs"`
	Pairing    bool `json:"pairing"`
}

// Capable allows drivers to report capabilities.
type Capable interface {
	Capabilities() Capabilities
}

// Discoverer allows drivers to discover devices.
type Discoverer interface {
	Discover(ctx context.Context) ([]DiscoveredDevice, error)
}

// DiscoveredDevice represents a discovered TV.
type DiscoveredDevice struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Driver   string `json:"driver"`
	Address  string `json:"address"`
	Model    string `json:"model,omitempty"`
	Firmware string `json:"firmware,omitempty"`
}
