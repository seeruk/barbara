package barbara

import "encoding/json"

// Config holds all application configuration.
type Config struct {
	Primary   WindowConfig `json:"primary"`
	Secondary WindowConfig `json:"secondary"`
}

// WindowConfig holds the configuration for a single on-screen bar.
type WindowConfig struct {
	Position WindowPosition    `json:"position"`
	Left     []json.RawMessage `json:"left"`
	Right    []json.RawMessage `json:"right"`
}

// ModuleConfig is the common configuration for a Barbara module.
type ModuleConfig struct {
	// Kind specifies the kind of module that this configuration is for. This allows the correct
	// ModuleFactory to be used to build the module.
	Kind string `json:"kind"`
}

// LoadConfig returns Barbara's configuration. It will either default to a directory under the
// user's home directory, or can be overridden via the environment.
func LoadConfig() (Config, error) {
	in := []byte(`{
		"primary": {
			"position": "bottom",
			"right": [
				{
					"kind": "menu",
					"label": "Elliot Wright",
					"items": [
						{ "label": "Log Off", "icon": "system-log-out", "exec": "i3-msg exit" },
						{ "separator": true },
						{ "label": "Reboot", "icon": "system-reboot", "exec": "sudo systemctl reboot" },
						{ "label": "Shutdown", "icon": "system-shutdown", "exec": "sudo systemctl poweroff" }
					]
				},
				{
					"kind": "clock",
					"format": "15:04:05\nMon, 02 Jan"
				}
			]
		}
	}`)

	var config Config

	err := json.Unmarshal(in, &config)

	return config, err
}
