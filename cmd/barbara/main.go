package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/seeruk/board/barbara"
	"github.com/seeruk/board/modules/clock"
	"github.com/seeruk/board/modules/menu"
	"github.com/therecipe/qt/widgets"
)

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetStyleSheet(barbara.BuildStyleSheet())

	// Create a bar for each screen.
	// TODO(elliot): React to screens being connected / disconnected and recreate all bars. This
	// means that Destroy methods of modules and the Window type must be very well implemented.
	screens := app.Screens()
	for _, screen := range screens {
		var leftModules []barbara.ModuleFactory
		var rightModules []barbara.ModuleFactory

		// TODO(elliot): This should definitely come from a configuration file.
		rightModules = append(rightModules, menu.NewModuleFactory())
		rightModules = append(rightModules, clock.NewModuleFactory())

		window := barbara.NewWindow(screen)
		window.Render(leftModules, rightModules)
	}

	go func() {
		for {
			var memstats runtime.MemStats

			runtime.ReadMemStats(&memstats)

			fmt.Println("CGo Calls:", runtime.NumCgoCall())
			fmt.Println("Routines:", runtime.NumGoroutine())
			fmt.Println("Heap:")
			fmt.Println(memstats.HeapAlloc)
			fmt.Println(memstats.HeapIdle)
			fmt.Println(memstats.HeapInuse)
			fmt.Println(memstats.HeapObjects)
			fmt.Println(memstats.HeapReleased)
			fmt.Println(memstats.HeapSys)
			fmt.Println()

			time.Sleep(10 * time.Second)
		}
	}()

	app.Exec()

	// TODO(elliot): Move all of this out of main, use resolver, use main for signal handling.
}
