package clock

import (
	"encoding/json"

	"github.com/seeruk/board/modules"
)

// ModuleFactory ...
type ModuleFactory struct{}

// NewModuleFactory ...
func NewModuleFactory() *ModuleFactory {
	return &ModuleFactory{}
}

// Build ...
func (f *ModuleFactory) Build(config json.RawMessage) modules.Module {
	return NewModule()
}
