package barbara

import (
	"encoding/json"
	"sync"

	"github.com/therecipe/qt/widgets"
)

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

// ModuleBuilder is a type that generalises the process of creating modules. A module factory can
// be instantiated with all dependencies needed for a module to function, and then it can build a
// module instance with some given configuration on-demand when a bar is being rendered.
type ModuleBuilder interface {
	// Build returns a new Module instance using the given configuration.
	// NOTE(elliot): This interface is likely to change over time as more module specific info needs
	// to be given to modules.
	Build(parent widgets.QWidget_ITF) (Module, error)
}

// AlignmentAwareModuleBuilder extends the ModuleBuilder interface to also allow setting a
// ModuleAlignment value for position elements within the module.
type AlignmentAwareModuleBuilder interface {
	ModuleBuilder
	SetAlignment(alignment ModuleAlignment)
}

// ConfigAwareModuleBuilder extends the ModuleBuilder interface to also allow settings some raw JSON
// formatted configuration which
type ConfigAwareModuleBuilder interface {
	ModuleBuilder
	SetConfig(raw json.RawMessage)
}

// WindowAwareModuleBuilder extends the ModuleBuilder interface to also allow setting a Window value
// for accessing things like Window position, and Window screen.
type WindowAwareModuleBuilder interface {
	ModuleBuilder
	SetWindow(window *Window)
}

// ModuleConfig is the common configuration for a Barbara module.
type ModuleConfig struct {
	// Kind specifies the kind of module that this configuration is for. This allows the correct
	// ModuleBuilder to be used to build the module.
	Kind string `json:"kind"`
}

// ModuleBuilderConstructor ...
type ModuleBuilderConstructor func() ModuleBuilder

// ModuleBuilderFactory is a type that can create ModuleBuilders that are registered with in it.
type ModuleBuilderFactory struct {
	sync.RWMutex
	// mbcs is a map of module name to ModuleBuilderConstructor, allowing new instances of a
	// ModuleBuilder to be constructed.
	mbcs map[string]ModuleBuilderConstructor
}

// ModuleBuilderFactory returns a new ModuleBuilderFactory instance.
func NewModuleBuilderFactory() *ModuleBuilderFactory {
	return &ModuleBuilderFactory{
		mbcs: make(map[string]ModuleBuilderConstructor),
	}
}

// Create attempts to use a ModuleBuilderConstructor to create a new ModuleBuilder instance, and
// return it. If a ModuleBuilderConstructor is not registered by the given name, the second return
// value will be false, and nil will be returned as the ModuleBuilder.
func (f *ModuleBuilderFactory) Create(name string) (ModuleBuilder, bool) {
	f.RLock()
	defer f.RUnlock()

	// Can't just return this, we need to make sure that we can run the ModuleBuilderConstructor.
	mbc, ok := f.mbcs[name]
	if !ok {
		return nil, ok
	}

	return mbc(), ok
}

// RegisterConstructor registers the given ModuleBuilderConstructor with the given name in this
// ModuleBuilderFactory instance, allowing it to be created by Create later.
func (f *ModuleBuilderFactory) RegisterConstructor(name string, mbc ModuleBuilderConstructor) {
	f.Lock()
	defer f.Unlock()

	f.mbcs[name] = mbc
}
