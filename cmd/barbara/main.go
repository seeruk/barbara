package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/davecgh/go-spew/spew"
	"github.com/seeruk/board/barbara"
	"github.com/seeruk/board/modules/clock"
	"github.com/seeruk/board/modules/menu"
	"github.com/therecipe/qt/widgets"
)

func main() {
	config, err := barbara.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(config)

	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetStyleSheet(barbara.BuildStyleSheet())

	// TODO(elliot): Use this:
	//app.PrimaryScreen()

	// Create a bar for each screen.
	// TODO(elliot): React to screens being connected / disconnected and recreate all bars. This
	// means that Destroy methods of modules and the Window type must be very well implemented.
	screens := app.Screens()
	for _, screen := range screens {
		var leftModules []barbara.ModuleFactory
		var rightModules []barbara.ModuleFactory

		// TODO(elliot): Is screen primary or not?
		for _, raw := range config.Primary.Right {
			var config barbara.ModuleConfig

			err := json.Unmarshal(raw, &config)
			if err != nil {
				// TODO(elliot): Some context?
				log.Fatal(err)
			}

			// TODO(elliot): This needs to be... better, and somewhere else. This whole block does.
			switch config.Kind {
			case "clock":
				rightModules = append(rightModules, clock.NewModuleFactory())
			case "menu":
				rightModules = append(rightModules, menu.NewModuleFactory(raw))
			}
		}

		window := barbara.NewWindow(screen)
		window.Render(leftModules, rightModules)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.Handle("/free", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			debug.FreeOSMemory()
		}))

		if err := http.ListenAndServe(":4000", nil); err != nil {
			log.Fatal(err)
		}
	}()

	app.Exec()

	// TODO(elliot): Move all of this out of main, use resolver, use main for signal handling.
}
