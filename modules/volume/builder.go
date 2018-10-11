package volume

import (
	"github.com/seeruk/barbara/barbara"
	"github.com/therecipe/qt/widgets"
)

// ModuleBuilder is a builder type that builds a new "volume" Module instance, housing shared
// dependencies to be passed in once the Module is created, etc.
type ModuleBuilder struct {
	window *barbara.Window
}

// NewModuleBuilder returns a new ModuleBuilder instance.
func NewModuleBuilder() *ModuleBuilder {
	return &ModuleBuilder{}
}

// SetWindow sets window in this builder.
func (b *ModuleBuilder) SetWindow(window *barbara.Window) {
	b.window = window
}

// Builder returns a new "volume" Module instance.
func (b *ModuleBuilder) Build(parent widgets.QWidget_ITF) (barbara.Module, error) {
	return nil, nil
}
