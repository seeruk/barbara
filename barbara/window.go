package barbara

import (
	"log"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// Window represents an on-screen Barbara "bar" window, as in, something that will manifest and
// manage the state of something that a window manager will manage. Externally, the QMainWindow
// widget will be what is used and "shown".
type Window struct {
	leftLayout   *widgets.QHBoxLayout
	rightLayout  *widgets.QHBoxLayout
	windowLayout *widgets.QHBoxLayout
	screen       *gui.QScreen
	window       *widgets.QMainWindow

	// The built, started modules that are currently in-use in this Window.
	modules []Module
}

// NewWindow creates a new instance of Window.
func NewWindow(screen *gui.QScreen) *Window {
	// Construct the window with all static parameters set.
	window := widgets.NewQMainWindow(nil, core.Qt__Window)
	window.SetWindowTitle("Barbara Bar")
	window.SetAttribute(core.Qt__WA_X11NetWmWindowTypeDock, true) // In X11, this makes the dock.
	// TODO(elliot): Wayland?

	return &Window{
		screen: screen,
		window: window,
	}
}

// createLayout attaches the layout widgets to this Window, providing the Qt containers that Barbara
// modules' Qt widgets can be placed in.
func (w *Window) createLayout() {
	w.leftLayout = widgets.NewQHBoxLayout()
	w.leftLayout.SetAlign(core.Qt__AlignLeft)

	w.rightLayout = widgets.NewQHBoxLayout()
	w.rightLayout.SetAlign(core.Qt__AlignRight)

	// Create the window layout, which will act as the layout for the underlying QMainWindow's
	// central widget's layout.
	w.windowLayout = widgets.NewQHBoxLayout()
	w.windowLayout.SetContentsMargins(7, 7, 7, 7)

	// Add the left and right layout widgets, providing them equal, but positive stretch so they
	// meet in the middle of the window by default.
	w.windowLayout.AddLayout(w.leftLayout, 1)
	w.windowLayout.AddLayout(w.rightLayout, 1)

}

// updateDimensions sets the height of the window based on the window's contents.
func (w *Window) updateDimensions() {
	if w.windowLayout == nil {
		return
	}

	w.window.SetFixedHeight(w.windowLayout.SizeHint().Height())
}

// updatePosition uses the geometry of the screen that this window will be displayed on, and moves
// the bar to the bottom of the screen.
func (w *Window) updatePosition() {
	geo := w.screen.Geometry()

	// Move the window into position.
	// TODO(elliot): Should support being at the top too.
	w.window.Move2(geo.X(), geo.Height()-w.window.Height())
}

// addModuleToLayout adds a module to the specified layout. Adding a module is complex enough that
// given the same functionality is needed for both ends of the bar, this function was necessary.
func (w *Window) addModuleToLayout(layout *widgets.QHBoxLayout, alignment core.Qt__AlignmentFlag, factory ModuleFactory) error {
	align := AlignmentLeft
	if alignment == core.Qt__AlignRight {
		align = AlignmentRight
	}

	module, err := factory.Build()
	if err != nil {
		return err
	}

	widget, err := module.Render(align, PositionBottom) // TODO(elliot): Un-hard-code.
	if err != nil {
		return err
	}

	// Add the widget the layout, add it to our list of modules, so we can destroy it later if
	// we need to re-render our Window(s).
	layout.AddWidget(widget, 0, alignment)

	// Add the module to the instantiated modules list, so that they may be destroyed later.
	w.modules = append(w.modules, module)

	return nil
}

// Render resizes, repositions, and then displays the window for this bar.
func (w *Window) Render(leftFactories, rightFactories []ModuleFactory) {
	// ...

	// Create the layout widgets.
	w.createLayout()

	for _, factory := range leftFactories {
		err := w.addModuleToLayout(w.leftLayout, core.Qt__AlignLeft, factory)
		if err != nil {
			// TODO(elliot): Add context.
			log.Println(err)
		}
	}

	for _, factory := range rightFactories {
		err := w.addModuleToLayout(w.rightLayout, core.Qt__AlignRight, factory)
		if err != nil {
			// TODO(elliot): Add context.
			log.Println(err)
		}
	}

	// Create the layout to the window so that all UI elements attached to the layout will be
	// displayed once the window is shown.
	centralWidget := widgets.NewQWidget(nil, 0)
	centralWidget.SetLayout(w.windowLayout)

	// Attach the central widget to the window.
	w.window.SetCentralWidget(centralWidget)

	// updateDimensions should also be called fairly late, prior to this we shouldn't really have a
	// set height - this means the window will have the correct size once spawned to fit all of it's
	// modules in. Maybe this behaviour will change to be static in the future.
	w.updateDimensions()

	// updatePosition must be called very late, to ensure the position is calculated correctly.
	w.updatePosition()

	// Finally, show the window.
	w.window.Show()
}

// Destroy stops all background processes in modules used in this bar, then destroys this window. In
// turn, all sub-windows are also destroyed.
func (w *Window) Destroy() {
	for _, module := range w.modules {
		err := module.Destroy()
		if err != nil {
			// TODO(elliot): Add context.
			log.Println(err)
		}
	}

	// Destroy everything, including sub-windows, and all widgets attached - freeing up resources.
	w.window.Destroy(true, true)
}
