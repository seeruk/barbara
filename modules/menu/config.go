package menu

// Config holds all menu module configuration.
type Config struct {
	Label string       `json:"label"`
	Items []ItemConfig `json:"items"`
}

// ItemConfig holds configuration for a single menu item.
type ItemConfig struct {
	Separator bool   `json:"separator,omitempty"`
	Label     string `json:"label,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Exec      string `json:"exec,omitempty"`
}
