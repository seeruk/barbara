package menu

import (
	"log"
	"os/exec"
	"strings"

	"github.com/seeruk/board/barbara"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// Module is a Barbara Module that presents a menu. Menus contain menu items that are able to be
// clicked. You can also include separators.
type Module struct {
	config Config

	button *widgets.QPushButton
	menu   *widgets.QMenu
}

// NewModule returns a new Module instance.
func NewModule(config Config) *Module {
	return &Module{
		config: config,
	}
}

// Render attempts to return a button widget that will open a menu containing some pre-configured
// menu items, ready to be placed onto a bar.
func (m *Module) Render(alignment barbara.Alignment, position barbara.Position) (widgets.QWidget_ITF, error) {
	button, err := m.createButton()
	if err != nil {
		return nil, err
	}

	m.menu = m.createMenu()

	m.button = button
	m.button.ConnectClicked(m.onButtonClicked(alignment, position))

	return m.button, nil
}

// createButton attempts to create a new button with the current user's name/username as it's label.
func (m *Module) createButton() (*widgets.QPushButton, error) {
	button := widgets.NewQPushButton2(m.config.Label, nil)
	button.SetFlat(true)
	button.SetProperty("class", core.NewQVariant14("barbara-button"))

	return button, nil
}

// createMenu creates a menu with pre-configured menu items attached.
func (m *Module) createMenu() *widgets.QMenu {
	menu := widgets.NewQMenu(nil)

	var items []*widgets.QAction
	for _, itemConfig := range m.config.Items {
		items = append(items, m.createMenuItem(itemConfig))
	}

	menu.AddActions(items)

	return menu
}

// createMenuItem creates a new menu item based on the given item configuration.
func (m *Module) createMenuItem(config ItemConfig) *widgets.QAction {
	if config.Separator {
		item := widgets.NewQAction(nil)
		item.SetSeparator(true)

		return item
	}

	item := widgets.NewQAction2(config.Label, nil)
	if config.Icon != "" {
		// TODO(elliot): Need some kind of generic icon-getting function, so I can get icons by
		// name, using the user's theme. Will need for tray anyway.
		icon := gui.NewQIcon()
		icon.SetThemeName("Paper")

		item = widgets.NewQAction3(icon.FromTheme(config.Icon), config.Label, nil)
	}

	// The "triggered" handler is run when the menu item is activated (e.g. clicked).
	item.ConnectTriggered(m.onMenuItemTriggered(config))

	return item
}

// onButtonClicked is the button click handler, used to show the menu.
func (m *Module) onButtonClicked(alignment barbara.Alignment, position barbara.Position) func(bool) {
	// Everything here needs to be handled dynamically, because the bar could move after being
	// rendered - so we recalculate menu position every click.
	return func(_ bool) {
		bsh := m.button.SizeHint()
		msh := m.menu.SizeHint()

		var x int
		if alignment == barbara.AlignmentRight {
			// Place on the right of the button by moving the menu right the whole width of the
			// button, minus the menu's width, lining up the right edge of the menu with the right
			// edge of the button.
			x = bsh.Width() - msh.Width()
		}

		var y int
		if position == barbara.PositionBottom {
			// Place above button, by moving the menu up the menu's height over the button.
			y = -msh.Height()
		} else {
			// Place under button, by moving the menu down the button's height.
			y = bsh.Height()
		}

		// Finally, show the menu.
		m.menu.Popup(m.button.MapToGlobal(core.NewQPoint2(x, y)), nil)
	}
}

// onMenuItemTriggered is the menu item activation handler, used to execute menu item commands.
func (m *Module) onMenuItemTriggered(config ItemConfig) func(bool) {
	args := strings.Split(config.Exec, " ")
	argc := len(args)

	return func(_ bool) {
		if argc == 0 {
			return
		}

		// TODO(elliot): With context?
		// Don't worry, this won't panic...
		cmd := exec.Command(args[0], args[1:]...)
		err := cmd.Run() // TODO(elliot): Use start, log output.
		if err != nil {
			// TODO(elliot): Add context.
			log.Println(err)
		}
	}
}

// Destroy frees up resources. There are no background processes in a menu module.
func (m *Module) Destroy() error {
	if m.button != nil {
		m.button.Destroy(true, true)
	}

	if m.menu != nil {
		m.menu.Destroy(true, true)
	}

	m.button = nil
	m.menu = nil

	return nil
}
