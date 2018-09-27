package clock

import (
	"github.com/seeruk/barbara/barbara"
	"github.com/therecipe/qt/widgets"
)

// ModuleFactory is a factory that produces new "clock" Module instances.
type ModuleFactory struct{}

// NewModuleFactory returns a new ModuleFactory instance.
func NewModuleFactory() *ModuleFactory {
	return &ModuleFactory{}
}

// Build returns a new "clock" Module instance.s
func (f *ModuleFactory) Build(parent widgets.QWidget_ITF) (barbara.Module, error) {
	return NewModule(parent), nil
}
