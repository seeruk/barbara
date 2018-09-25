package main

import (
	"os"

	"github.com/seeruk/board/bar"
	"github.com/seeruk/board/modules"
	"github.com/seeruk/board/modules/clock"
	"github.com/therecipe/qt/widgets"
)

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	// TODO(elliot): Apply global styles here, elsewhere should use class names.

	// Create a bar for each screen.
	// TODO(elliot): React to screens being connected / disconnected and recreate all bars. This
	// means that Destroy methods of modules and the Window type must be very well implemented.
	screens := app.Screens()
	for _, screen := range screens {
		var leftModules []modules.ModuleFactory
		var rightModules []modules.ModuleFactory

		rightModules = append(rightModules, clock.NewModuleFactory())

		window := bar.NewWindow(screen)
		window.Render(leftModules, rightModules)
	}

	app.Exec()
}
