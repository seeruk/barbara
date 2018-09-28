package main

import (
	"log"
	"os"

	"github.com/seeruk/barbara/barbara"
	"github.com/therecipe/qt/widgets"
)

func main() {
	config, err := barbara.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetStyleSheet(barbara.BuildStyleSheet())

	// Get primary screen so we know which bar config to load.
	primaryScreen := app.PrimaryScreen()

	// Create a bar for each screen.
	// TODO(elliot): React to screens being connected / disconnected and recreate all bars. This
	// means that Destroy methods of modules and the Window type must be very well implemented.
	screens := app.Screens()
	for _, screen := range screens {
		barConfig := config.Secondary
		if primaryScreen != nil && screen.Name() == primaryScreen.Name() {
			barConfig = config.Primary
		}

		leftModules := barbara.BuildModules(barConfig.Left)
		rightModules := barbara.BuildModules(barConfig.Right)

		window := barbara.NewWindow(screen, barConfig.Position)
		window.Render(leftModules, rightModules)
	}

	app.Exec()

	// TODO(elliot): Move all of this out of main, use resolver, use main for signal handling.
}
