package modules

import (
	"encoding/json"

	"github.com/therecipe/qt/widgets"
)

// Module ...
type Module interface {
	// Render ...
	Render() (widgets.QWidget_ITF, error)
	// Destroy ...
	Destroy() error
}

// ModuleFactory ...
type ModuleFactory interface {
	Build(config json.RawMessage) Module
}
