package barbara

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/therecipe/qt/widgets"
)

var (
	// mfcs is a map from module name to it's ModuleFactory's constructor function. Registered
	// module types will end up in this map.
	mfcs map[string]ModuleFactoryConstructor
	// mfcsMu protects access to the mfcs map.
	mfcsMu = &sync.Mutex{}
)

// BuildModules takes a slice of raw JSON configuration bytes, and uses it build
func BuildModules(config []json.RawMessage) []ModuleFactory {
	var factories []ModuleFactory

	for _, raw := range config {
		var moduleConf ModuleConfig

		err := json.Unmarshal(raw, &moduleConf)
		if err != nil {
			// TODO(elliot): Some context?
			log.Fatal(err)
		}

		factories = append(factories, mfcs[moduleConf.Kind](raw))
	}

	return factories
}

// RegisterModule takes a module name, and a ModuleFactoryConstructor and registers it so that it
// can be used to construct new Module instances to display on a bar.
func RegisterModule(name string, mfc ModuleFactoryConstructor) {
	mfcsMu.Lock()
	defer mfcsMu.Unlock()

	if mfcs == nil {
		mfcs = make(map[string]ModuleFactoryConstructor)
	}

	mfcs[name] = mfc
}

// Module represents a Barbara bar module, i.e. a combination of UI and functionality that can be
// presented on a Barbara bar.
type Module interface {
	// Render attempts to return a QWidget, which will be placed in one of the bar's layout boxes.
	Render() (widgets.QWidget_ITF, error)
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
	// NOTE(elliot): This interface is likely to change over time as more module specific info needs
	// to be given to modules.
	Build(parent widgets.QWidget_ITF) (Module, error)
}

// ModuleFactoryConstructor is a function that is used to construct a new ModuleFactory instance.
type ModuleFactoryConstructor func(raw json.RawMessage) ModuleFactory

// AlignmentAwareModuleFactory extends the ModuleFactory interface to also allow setting a
// ModuleAlignment value for position elements within the module.
type AlignmentAwareModuleFactory interface {
	ModuleFactory
	SetAlignment(alignment ModuleAlignment)
}

// WindowAwareModuleFactory extends the ModuleFactory interface to also allow setting a Window value
// for accessing things like Window position, and Window screen.
type WindowAwareModuleFactory interface {
	ModuleFactory
	SetWindow(window *Window)
}
