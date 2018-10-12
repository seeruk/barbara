package clock

// Config holds all clock module configuration.
type Config struct {
	// Format is a Go time format string. See the Go documentation for more information.
	Format string `json:"format"`
}
