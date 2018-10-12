package menu

// Config holds all menu module configuration.
type Config struct {
	// Label is the text used on the menu button.
	Label string `json:"label"`
	// Items is a slice of menu items.
	Items []ItemConfig `json:"items"`
}

// ItemConfig holds configuration for a single menu item.
type ItemConfig struct {
	// Separator is true if this item is a separator.
	Separator bool `json:"separator,omitempty"`
	// Label is the text used for this menu item.
	Label string `json:"label,omitempty"`
	// Icon is an icon name.
	// TODO(elliot): Document further once I know more about how this is going to work.
	Icon string `json:"icon,omitempty"`
	// Exec specifies the command to execute when this menu item is activated.
	Exec string `json:"exec,omitempty"`
}
