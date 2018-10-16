package battery

// Config holds all battery module configuration.
type Config struct {
	// PowerSupply is the name of the power supply found in /sys/class/power_supply, e.g. BAT0. The
	// power supply must actually be a battery.
	PowerSupply string `json:"power_supply"`
	// TODO(elliot): Specifying defaults would be useful for things like refresh interval...
}
