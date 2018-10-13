package battery

// Config holds all battery module configuration.
type Config struct {
	// PowerSupply is the name of the power supply found in /sys/class/power_supply, e.g. BAT0. The
	// power supply must actually be a battery.
	PowerSupply string `json:"power_supply"`
}

// Notes:
// * For me: `cd /sys/class/power_supply/BAT0`
// * Time remaining to empty: `calc (cat charge_now)" / "(cat current_now)`
// * Time remaining to full: `calc "("(cat charge_full)" - "(cat charge_now)") / "(cat current_now)`
// * Current %age: `cat capacity`
// * From hours (decimal) into minutes:  0.6 * (h * 100)
// * Show icon indicating charge, maybe use Paper-Mono-Dark theme?
