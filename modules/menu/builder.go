package menu

import (
	"encoding/json"

	"github.com/seeruk/barbara/barbara"
	"github.com/therecipe/qt/widgets"
)

// ModuleBuilder is a builder that produces new "menu" Module instances, housing shared dependencies
// to be passed in once the Module is created, etc.
type ModuleBuilder struct {
	config json.RawMessage

	alignment barbara.ModuleAlignment
	window    *barbara.Window
}

// NewModuleBuilder returns a new ModuleBuilder instance.
func NewModuleBuilder() barbara.ModuleBuilder {
	return &ModuleBuilder{}
}

// SetAlignment sets alignment in this builder.
func (b *ModuleBuilder) SetAlignment(alignment barbara.ModuleAlignment) {
	b.alignment = alignment
}

// SetConfig sets config in this builder.
func (b *ModuleBuilder) SetConfig(config json.RawMessage) {
	b.config = config
}

// SetWindow sets window in this builder.
func (b *ModuleBuilder) SetWindow(window *barbara.Window) {
	b.window = window
}

// Build returns a new "menu" Module instance.
func (b *ModuleBuilder) Build(parent widgets.QWidget_ITF) (barbara.Module, error) {
	var config Config

	err := json.Unmarshal(b.config, &config)
	if err != nil {
		// TODO(elliot): Add context.
		return nil, err
	}

	return NewModule(config, b.alignment, b.window.Position(), parent), nil
}
