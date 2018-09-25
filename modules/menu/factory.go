package menu

import (
	"encoding/json"

	"github.com/seeruk/board/barbara"
)

// ModuleFactory is a factory that produces new "menu" Module instances.
type ModuleFactory struct{}

// NewModuleFactory returns a new ModuleFactory instance.
func NewModuleFactory() *ModuleFactory {
	return &ModuleFactory{}
}

// Build returns a new "menu" Module instance.s
func (f *ModuleFactory) Build(config json.RawMessage) barbara.Module {
	return NewModule(Config{
		Items: []ItemConfig{
			{
				Label: "&Log Off",
				Icon:  "system-log-out",
				Exec:  "i3-msg exit",
			},
			{
				Separator: true,
			},
			{
				Label: "&Reboot",
				Icon:  "system-reboot",
				Exec:  "sudo systemctl reboot",
			},
			{
				Label: "&Shutdown",
				Icon:  "system-shutdown",
				Exec:  "sudo systemctl poweroff",
			},
		},
	})
}
