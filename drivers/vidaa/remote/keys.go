package remote

var KeyMap = map[string]string{
	"Up":         "KEY_UP",
	"Down":       "KEY_DOWN",
	"Left":       "KEY_LEFT",
	"Right":      "KEY_RIGHT",
	"Select":     "KEY_OK",
	"Back":       "KEY_RETURNS",
	"Home":       "KEY_HOME",
	"Play":       "KEY_PLAY",
	"Pause":      "KEY_PAUSE",
	"Stop":       "KEY_STOP",
	"Rewind":     "KEY_BACK",
	"Forward":    "KEY_FORWARDS",
	"Power":      "KEY_POWER",
	"VolumeUp":   "KEY_VOLUMEUP",
	"VolumeDown": "KEY_VOLUMEDOWN",
	"Mute":       "KEY_MUTE",
}

func MapKey(k string) (string, bool) {
	v, ok := KeyMap[k]
	return v, ok
}
