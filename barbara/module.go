package barbara

import (
	"encoding/json"

	"github.com/therecipe/qt/widgets"
)

// Module represents a Barbara bar module, i.e. a combination of UI and functionality that can be
// presented on a Barbara bar.
type Module interface {
	// Render attempts to return a QWidget, which will be placed in one of the bar's layout boxes.
	Render() (widgets.QLayout_ITF, error)
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

// ModuleConfig is the common configuration for a Barbara module.
type ModuleConfig struct {
	// Kind specifies the kind of module that this configuration is for, allowing the correct Module
	// to be constructed based on the kind.
	Kind string `json:"kind"`
}

// ModuleConstructorFunc is a function used to construct new Module instances.
type ModuleConstructorFunc func(mctx ModuleContext) (Module, error)

// ModuleContext is used to inform a module about it's environment on the bar, e.g. it's alignment,
// the position of the bar itself, and the Module's configuration.
type ModuleContext struct {
	// Alignment is the intended alignment of the Module on a Barbara bar (i.e. left, right).
	Alignment ModuleAlignment
	// Config is the raw configuration bytes. The Module will have to decode it's configuration.
	Config json.RawMessage
	// Window is the Barbara bar's window representation, allowing the module to get info about the
	// window itself, such as it's position on the screen it's on.
	Window *Window
}

// ModuleFactory is a type that ModuleConstructorFunc functions can be registered in to create new
// instances of modules on-demand.
type ModuleFactory struct {
	mcfs map[string]ModuleConstructorFunc
}

// NewModuleFactory constructs a new ModuleFactory instance.
func NewModuleFactory() *ModuleFactory {
	return &ModuleFactory{
		mcfs: make(map[string]ModuleConstructorFunc),
	}
}

// Create attempts to created a new instance of a Module. The Module must be registered first. If an
// unknown Module is requested, (nil, false) will be returned.
func (f *ModuleFactory) Create(name string, mctx ModuleContext) (Module, bool) {
	if mcf, ok := f.mcfs[name]; ok {
		module, err := mcf(mctx)
		if err != nil {
			// TODO(elliot): Do something?
			return nil, false
		}

		return module, true
	}

	return nil, false
}

// RegisterConstructor registers the given ModuleConstructorFunc with the given name in this
// ModuleFactory instance, allowing a new Module instance to be created later.
func (f *ModuleFactory) RegisterConstructor(name string, mcf ModuleConstructorFunc) {
	f.mcfs[name] = mcf
}
