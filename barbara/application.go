package barbara

import (
	"os"

	"github.com/therecipe/qt/core"
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
)

// Application is a type that sets up the Barbara QApplication, connecting event handlers, and
// orchestrating the lifecycle of Barbara's UI.
// TODO(elliot): Application is a bit of a rubbish name.
type Application struct {
	app *widgets.QApplication

	primaryConfig   WindowConfig
	secondaryConfig WindowConfig
}

// NewApplication returns a new instance of Application.
func NewApplication(primaryConfig, secondaryConfig WindowConfig) *Application {
	application := &Application{
		// TODO(elliot): Not exactly testable, is it this?
		app:             widgets.NewQApplication(len(os.Args), os.Args),
		primaryConfig:   primaryConfig,
		secondaryConfig: secondaryConfig,
	}

	application.applyEventHandlers()
	application.applyStylesheet()

	return application
}

// CreateWindows submits an internal event that triggers the creation of Barbara bar windows, and
// starts all modules configured for those windows.
func (a *Application) CreateWindows() {
	a.postEvent(eventCreateWindows)
}

// DestroyWindows submits an internal event that triggers the destruction of all Barbara bar
// windows, and halts all background processes from modules.
func (a *Application) DestroyWindows() {
	a.postEvent(eventDestroyWindows)
}

// RecreateWindows basically does what calling DestroyWindows and then CreateWindows would do.
func (a *Application) RecreateWindows() {
	a.postEvent(eventRecreateWindows)
}

// postEvent provides an easier way to send an event to the underlying QApplication.
func (a *Application) postEvent(eventType core.QEvent__Type) {
	a.app.PostEvent(a.app, core.NewQEvent(eventType), 0)
}

// QApplication returns the already instantiated QApplication instance.
func (a *Application) QApplication() *widgets.QApplication {
	return a.app
}

func (a *Application) applyEventHandlers() {
	var windows []*Window

	a.app.ConnectEvent(func(e *core.QEvent) bool {
		switch e.Type() {
		case eventCreateWindows:
			windows = make([]*Window, 0, 0) // Reset

			// Get primary screen so we know which bar config to load.
			primaryScreen := a.app.PrimaryScreen()

			// Create a bar for each screen.
			screens := a.app.Screens()
			for _, screen := range screens {
				barConfig := a.secondaryConfig
				if primaryScreen != nil && screen.Name() == primaryScreen.Name() {
					barConfig = a.primaryConfig
				}

				leftModules := BuildModules(barConfig.Left)
				rightModules := BuildModules(barConfig.Right)

				window := NewWindow(screen, barConfig.Position)
				window.Render(leftModules, rightModules)

				windows = append(windows, window)
			}
		case eventDestroyWindows:
			for _, window := range windows {
				window.Destroy()
			}
		case eventRecreateWindows:
			for _, window := range windows {
				window.Destroy()
			}

			// Send event to create new windows.
			a.app.PostEvent(a.app, core.NewQEvent(eventCreateWindows), 0)
		}

		return true
	})
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
