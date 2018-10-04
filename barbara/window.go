package barbara

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// Window represents an on-screen Barbara "bar" window, as in, something that will manifest and
// manage the state of something that a window manager will manage. Externally, the QMainWindow
// widget will be what is used and "shown".
type Window struct {
	modules  []Module
	position WindowPosition

	screen       *gui.QScreen
	leftLayout   *widgets.QHBoxLayout
	rightLayout  *widgets.QHBoxLayout
	windowLayout *widgets.QHBoxLayout
	window       *widgets.QMainWindow
}

// NewWindow creates a new instance of Window.
func NewWindow(screen *gui.QScreen, position WindowPosition) *Window {
	// Construct the window with all static parameters set.
	window := widgets.NewQMainWindow(nil, core.Qt__Window)
	window.SetWindowTitle("Barbara Bar")
	window.SetAttribute(core.Qt__WA_X11NetWmWindowTypeDock, true) // In X11, this makes the dock.
	// TODO(elliot): Wayland?

	return &Window{
		position: position,
		screen:   screen,
		window:   window,
	}
}

// createLayout attaches the layout widgets to this Window, providing the Qt containers that Barbara
// modules' Qt widgets can be placed in.
func (w *Window) createLayout(parent widgets.QWidget_ITF) {
	// Create the window layout, which will act as the layout for the underlying QMainWindow's
	// central widget's layout.
	w.windowLayout = widgets.NewQHBoxLayout2(parent)
	w.windowLayout.SetContentsMargins(7, 7, 7, 7)

	w.leftLayout = widgets.NewQHBoxLayout2(w.windowLayout.Widget())
	w.leftLayout.SetAlign(core.Qt__AlignLeft)

	w.rightLayout = widgets.NewQHBoxLayout2(w.windowLayout.Widget())
	w.rightLayout.SetAlign(core.Qt__AlignRight)

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

	switch w.position {
	case WindowPositionTop:
		w.window.Move2(geo.X(), geo.Y())
	default:
		// Default is bottom.
		w.window.Move2(geo.X(), geo.Y()+geo.Height()-w.window.Height())
	}
}

// addModuleToLayout adds a module to the specified layout. Adding a module is complex enough that
// given the same functionality is needed for both ends of the bar, this function was necessary.
func (w *Window) addModuleToLayout(layout *widgets.QHBoxLayout, alignment core.Qt__AlignmentFlag, factory ModuleFactory) error {
	align := ModuleAlignmentLeft
	if alignment == core.Qt__AlignRight {
		align = ModuleAlignmentRight
	}

	// If the module factory needs to be made aware of the alignment of the module, then set it.
	if f, ok := factory.(AlignmentAwareModuleFactory); ok {
		f.SetAlignment(align)
	}

	// If the module factory needs to be made aware of the Window itself, then set it.
	if f, ok := factory.(WindowAwareModuleFactory); ok {
		f.SetWindow(w)
	}

	module, err := factory.Build(layout.Widget())
	if err != nil {
		return err
	}

	widget, err := module.Render()
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

	// Create the layout to the window so that all UI elements attached to the layout will be
	// displayed once the window is shown.
	centralWidget := widgets.NewQWidget(w.window, 0)

	// Create the layout widgets.
	w.createLayout(centralWidget)

	// Assign the newly created layout to the central widget.
	centralWidget.SetLayout(w.windowLayout)

	// Add all of the configured Barbara modules to their corresponding layout boxes.
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

// Position returns the position of this Window.
func (w *Window) Position() WindowPosition {
	return w.position
}

// Screen returns the QScreen that this Window is placed on.
func (w *Window) Screen() *gui.QScreen {
	return w.screen
}

// WindowConfig holds the configuration for a single on-screen bar.
type WindowConfig struct {
	Position WindowPosition    `json:"position"`
	Left     []json.RawMessage `json:"left"`
	Right    []json.RawMessage `json:"right"`
}

const (
	// WindowPositionTop is passed to modules when rendered at the top of the screen.
	WindowPositionTop WindowPosition = iota
	// WindowPositionBottom is passed to modules when rendered at the bottom of the screen.
	WindowPositionBottom
)

// WindowPosition represents the possible positions of a Barbara bar on the screen.
type WindowPosition int

// UnmarshalJSON allows a JSON string to be unmarshalled into a WindowPosition.
func (p *WindowPosition) UnmarshalJSON(raw []byte) error {
	var str string

	err := json.Unmarshal(raw, &str)
	if err != nil {
		return err
	}

	switch strings.ToLower(str) {
	case "top":
		*p = WindowPositionTop
	case "bottom":
		*p = WindowPositionBottom
	default:
		return fmt.Errorf("invalid position %q", str)
	}

	return nil
}
