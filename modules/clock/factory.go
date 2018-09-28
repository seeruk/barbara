package clock

import (
	"encoding/json"

	"github.com/seeruk/barbara/barbara"
	"github.com/therecipe/qt/widgets"
)

func init() {
	barbara.RegisterModule("clock", NewModuleFactory)
}

// ModuleFactory is a factory that produces new "clock" Module instances.
type ModuleFactory struct{}

// NewModuleFactory returns a new ModuleFactory instance.
func NewModuleFactory(_ json.RawMessage) barbara.ModuleFactory {
	return &ModuleFactory{}
}

// Build returns a new "clock" Module instance.s
func (f *ModuleFactory) Build(parent widgets.QWidget_ITF) (barbara.Module, error) {
	return NewModule(parent), nil
}
