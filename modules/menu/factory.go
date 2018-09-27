package menu

import (
	"encoding/json"

	"github.com/seeruk/barbara/barbara"
	"github.com/therecipe/qt/widgets"
)

// ModuleFactory is a factory that produces new "menu" Module instances.
type ModuleFactory struct {
	config json.RawMessage

	alignment barbara.ModuleAlignment
	window    *barbara.Window
}

// NewModuleFactory returns a new ModuleFactory instance.
func NewModuleFactory(config json.RawMessage) *ModuleFactory {
	return &ModuleFactory{
		config: config,
	}
}

// SetAlignment sets alignment in this factory.
func (f *ModuleFactory) SetAlignment(alignment barbara.ModuleAlignment) {
	f.alignment = alignment
}

// SetWindow sets window in this factory.
func (f *ModuleFactory) SetWindow(window *barbara.Window) {
	f.window = window
}

// Build returns a new "menu" Module instance.
func (f *ModuleFactory) Build(parent widgets.QWidget_ITF) (barbara.Module, error) {
	var config Config

	err := json.Unmarshal(f.config, &config)
	if err != nil {
		// TODO(elliot): Add context.
		return nil, err
	}

	return NewModule(config, f.alignment, f.window.Position(), parent), nil
}
