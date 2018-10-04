package clock

import (
	"github.com/seeruk/barbara/barbara"
	"github.com/therecipe/qt/widgets"
)

// ModuleBuilder is a builder that produces new "clock" Module instances.
type ModuleBuilder struct {
	// TODO(elliot): Config...
}

// NewModuleBuilder returns a new ModuleBuilder instance.
func NewModuleBuilder() barbara.ModuleBuilder {
	return &ModuleBuilder{}
}

// Build returns a new "clock" Module instance.
func (f *ModuleBuilder) Build(parent widgets.QWidget_ITF) (barbara.Module, error) {
	return NewModule(parent), nil
}
