package barbara

import (
	"github.com/therecipe/qt/widgets"
)

// Module represents a Barbara bar module, i.e. a combination of UI and functionality that can be
// presented on a Barbara bar.
type Module interface {
	// Render attempts to return a QWidget, which will be placed in one of the bar's layout boxes.
	Render(alignment ModuleAlignment, position WindowPosition) (widgets.QWidget_ITF, error)
	// Destroy frees up all resources for this module, stopping any background processes.
	Destroy() error
}

const (
	// ModuleAlignmentLeft is passed to modules when rendered on the left side of a bar.
	ModuleAlignmentLeft ModuleAlignment = iota
	// ModuleAlignmentRight is passed to modules when rendered on the right side of a bar.
	ModuleAlignmentRight
)

// ModuleAlignment represents the possible alignment of a module in the bar.
type ModuleAlignment int

// ModuleFactory is a type that generalises the process of creating modules. A module factory can
// be instantiated with all dependencies needed for a module to function, and then it can build a
// module instance with some given configuration on-demand when a bar is being rendered.
type ModuleFactory interface {
	// Build returns a new Module instance using the given configuration.
	Build(parent widgets.QWidget_ITF) (Module, error)
}
