package menu

import (
	"encoding/json"

	"github.com/seeruk/board/barbara"
)

// ModuleFactory is a factory that produces new "menu" Module instances.
type ModuleFactory struct {
	config json.RawMessage
}

// NewModuleFactory returns a new ModuleFactory instance.
func NewModuleFactory(config json.RawMessage) *ModuleFactory {
	return &ModuleFactory{
		config: config,
	}
}

// Build returns a new "menu" Module instance.s
func (f *ModuleFactory) Build() (barbara.Module, error) {
	var config Config

	err := json.Unmarshal(f.config, &config)
	if err != nil {
		// TODO(elliot): Add context.
		return nil, err
	}

	return NewModule(config), nil
}
