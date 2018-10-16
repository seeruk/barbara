package battery

// powerSupplyPath is the base path on a Linux system where battery information can be found.
const powerSupplyPath = "/sys/class/power_supply"

// Info contains all information that we need in the battery module. Some of this information is
// shown on the menu popup.
type Info struct {
	Capacity         float64 // capacity
	ChargeFull       float64 // charge_full
	ChargeFullDesign float64 // charge_full_design
	ChargeNow        float64 // charge_now
	CurrentNow       float64 // current_now
	Manufacturer     string  // manufacturer
	ModelName        string  // model_name
	Status           string  // status
	Technology       string  // technology
}
