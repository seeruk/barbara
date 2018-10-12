package barbara

import (
	"encoding/json"
	"log"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var (
	// eventCreateWindows is the event used to render and start all Barbara UI and module activity.
	eventCreateWindows = core.QEvent__Type(2000)
	// eventDestroyWindows is the event used to completely destroy Barbara's UI, and as a result, it
	// will also halt background module activity.
	eventDestroyWindows = core.QEvent__Type(2001)
	// eventRecreateWindows is the event used to trigger a full UI re-render, effectively restarting
	// Barbara in-place. This will also restart modules.
	eventRecreateWindows = core.QEvent__Type(2002)
	// eventExit is used to signal to the QApplication event loop that it should exit.
	eventExit = core.QEvent__Type(2003)
)

// Application is a type that sets up the Barbara QApplication, connecting event handlers, and
// orchestrating the lifecycle of Barbara's UI.
// TODO(elliot): Application is a bit of a rubbish name.
type Application struct {
	app     *widgets.QApplication
	windows []*Window

	moduleFactory   *ModuleFactory
	primaryConfig   WindowConfig
	secondaryConfig WindowConfig
}

// NewApplication returns a new instance of Application.
func NewApplication(
	moduleFactory *ModuleFactory,
	primaryConfig, secondaryConfig WindowConfig,
) *Application {
	application := &Application{
		// TODO(elliot): Not exactly testable, is it this?
		app:             widgets.NewQApplication(len(os.Args), os.Args),
		moduleFactory:   moduleFactory,
		primaryConfig:   primaryConfig,
		secondaryConfig: secondaryConfig,
	}

	application.applyEventHandlers()
	application.applyStylesheet()

	return application
}

// CreateWindows provides a thread-safe mechanism for running the code that handles creating all
// Barbara bars. Internally it uses Qt's event system to ensure that the event is handled on the
// main thread.
func (a *Application) CreateWindows() {
	a.postEvent(eventCreateWindows)
}

// DestroyWindows provides a thread-safe mechanism for running the code that handles destroying all
// Barbara bars. Internally it uses Qt's event system to ensure that the event is handled on the
// main thread.
func (a *Application) DestroyWindows() {
	a.postEvent(eventDestroyWindows)
}

// RecreateWindows basically does what calling DestroyWindows and then CreateWindows would do.
func (a *Application) RecreateWindows() {
	a.postEvent(eventRecreateWindows)
}

// Exit provides a thread-safe mechanism for signalling for the QApplication to exit gracefully. It
// also destroys all windows, stopping all modules.
func (a *Application) Exit() {
	a.DestroyWindows()
	a.postEvent(eventExit)
}

// postEvent provides an easier way to send an event to the underlying QApplication.
func (a *Application) postEvent(eventType core.QEvent__Type) {
	a.app.PostEvent(a.app, core.NewQEvent(eventType), 0)
}

// QApplication returns the already instantiated QApplication instance.
func (a *Application) QApplication() *widgets.QApplication {
	return a.app
}

// applyEventHandlers connects the QApplication with the "user-defined" events from this package.
func (a *Application) applyEventHandlers() {
	// We use Qt's event handling internally here so that the event handling code is executed in the
	// main thread. This keeps Qt happy with our concurrent code.
	a.app.ConnectEvent(func(e *core.QEvent) bool {
		switch e.Type() {
		case eventCreateWindows:
			a.onCreateWindowsEvent()
		case eventDestroyWindows:
			a.onDestroyWindowsEvent()
		case eventRecreateWindows:
			a.onRecreateWindowsEvent()
		case eventExit:
			a.onExit()
		}

		return true
	})
}

// onCreateWindowsEvent is an internal event handler run via Qt when a Qt user event with the type
// defined in eventCreateWindows is received.
func (a *Application) onCreateWindowsEvent() {
	// Get primary screen so we know which bar config to load, and all screens to iterate over.
	primaryScreen := a.app.PrimaryScreen()
	screens := a.app.Screens()

	// Create a bar for each screen.
	a.windows = make([]*Window, 0, len(screens)) // Reset
	for _, screen := range screens {
		a.windows = append(a.windows, a.createWindow(primaryScreen, screen))
	}
}

// createWindow ...
func (a *Application) createWindow(primaryScreen, screen *gui.QScreen) *Window {
	config := a.secondaryConfig
	if primaryScreen != nil && screen.Name() == primaryScreen.Name() {
		config = a.primaryConfig
	}

	window := NewWindow(config, screen)

	leftModules := a.createModules(ModuleAlignmentLeft, config.Left, window)
	rightModules := a.createModules(ModuleAlignmentRight, config.Right, window)

	window.Render(leftModules, rightModules)

	return window
}

// createModules ...
func (a *Application) createModules(alignment ModuleAlignment, rawConfigs []json.RawMessage, window *Window) []Module {
	modules := make([]Module, 0, len(rawConfigs))

	for _, rawConfig := range rawConfigs {
		var moduleConfig ModuleConfig

		// First, determine the Module kind.
		err := json.Unmarshal(rawConfig, &moduleConfig)
		if err != nil {
			// TODO(elliot): Better logging.
			log.Println("failed to unmarshal module configuration")
			continue
		}

		mctx := ModuleContext{
			Alignment: alignment,
			Config:    rawConfig,
			Window:    window,
		}

		module, ok := a.moduleFactory.Create(moduleConfig.Kind, mctx)
		if !ok {
			log.Printf("failed to create module with kind %q", moduleConfig.Kind)
			continue
		}

		modules = append(modules, module)
	}

	return modules
}

// onDestroyWindowsEvent is an internal event handler run via Qt when a Qt user event with the type
// defined in eventDestroyWindows is received.
func (a *Application) onDestroyWindowsEvent() {
	for _, window := range a.windows {
		window.Destroy()
	}
}

// onRecreateWindowsEvent is an internal event handler run via Qt when a Qt user event with the type
// defined in eventRecreateWindows is received.
func (a *Application) onRecreateWindowsEvent() {
	for _, window := range a.windows {
		window.Destroy()
	}

	// Send event to create new windows.
	a.app.PostEvent(a.app, core.NewQEvent(eventCreateWindows), 0)
}

// onExit is an internal event handler run via Qt when a Qt user event with the type defined in
// eventQuit is received.
func (a *Application) onExit() {
	a.app.Exit(0)
}

// applyStylesheet applies the global stylesheet for the application.
// TODO(elliot): Templating, configuration.
func (a *Application) applyStylesheet() {
	a.app.SetStyleSheet(`
		QMainWindow {
			background: #1a1a1a;
			margin: 0;
			padding: 0px;
		}

		QLabel {
			color: #e5e5e5;
			font-family: "Fira Sans";
			font-size: 13px;
			padding: 0 0 0 7px;
			text-align: center;
		}

		.barbara-button {
			background-color: #1a1a1a;
			color: #e5e5e5;
			font-family: "Fira Sans";
			font-size: 13px;
			padding: 7px;
		}

		.barbara-button:flat {
			border: 1px solid #333;
			border-radius: 3px;
		}
	`)
}
