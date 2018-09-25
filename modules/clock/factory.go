package clock

import (
	"encoding/json"

	"github.com/seeruk/board/barbara"
)

// ModuleFactory is a factory that produces new "clock" Module instances.
type ModuleFactory struct{}

// NewModuleFactory returns a new ModuleFactory instance.
func NewModuleFactory() *ModuleFactory {
	return &ModuleFactory{}
}

// Build returns a new "clock" Module instance.s
func (f *ModuleFactory) Build(config json.RawMessage) barbara.Module {
	return NewModule()
}
