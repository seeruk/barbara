package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/therecipe/qt/core"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/seeruk/barbara/barbara"
	"github.com/seeruk/barbara/wm/x11"
	"github.com/therecipe/qt/widgets"

	_ "github.com/seeruk/barbara/modules/clock"
	_ "github.com/seeruk/barbara/modules/menu"
)

var (
	// EventCreateWindows ...
	EventCreateWindows = core.QEvent__Type(2000)
	// EventDestroyWindows ...
	EventDestroyWindows = core.QEvent__Type(2001)
)

func main() {
	xc, _ := xgb.NewConn()

	// Every extension must be initialized before it can be used.
	err := randr.Init(xc)
	if err != nil {
		log.Fatal(err)
	}

	// Get the root window on the default screen.
	root := xproto.Setup(xc).DefaultScreen(xc).Root

	// Tell RandR to send us events. (I think these are all of them, as of 1.3.)
	err = randr.SelectInputChecked(xc, root,
		randr.NotifyMaskScreenChange|
			randr.NotifyMaskCrtcChange|
			randr.NotifyMaskOutputChange|
			randr.NotifyMaskOutputProperty).Check()

	if err != nil {
		log.Fatal(err)
	}

	watcher := x11.NewRandrEventWatcher(xc)
	xevCh := watcher.Start(context.Background())

	config, err := barbara.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	_ = config

	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetStyleSheet(barbara.BuildStyleSheet())

	deletedCh := make(chan struct{})

	go func() {
		for {
			<-xevCh

			// Send event to destroy existing windows.
			app.PostEvent(app, core.NewQEvent(EventDestroyWindows), 0)

			<-deletedCh

			// Send event to create new windows.
			app.PostEvent(app, core.NewQEvent(EventCreateWindows), 0)
		}
	}()

	var windows []*barbara.Window

	app.ConnectEvent(func(e *core.QEvent) bool {
		switch e.Type() {
		case EventCreateWindows:
			windows = make([]*barbara.Window, 0, 0) // Reset

			fmt.Println("1")

			// Get primary screen so we know which bar config to load.
			primaryScreen := app.PrimaryScreen()
			fmt.Println("2")

			// Create a bar for each screen.
			screens := app.Screens()
			fmt.Println("3")
			for _, screen := range screens {
				barConfig := config.Secondary
				if primaryScreen != nil && screen.Name() == primaryScreen.Name() {
					barConfig = config.Primary
				}
				fmt.Println("4")

				leftModules := barbara.BuildModules(barConfig.Left)
				rightModules := barbara.BuildModules(barConfig.Right)

				fmt.Println("5")
				window := barbara.NewWindow(screen, barConfig.Position)
				window.Render(leftModules, rightModules)

				fmt.Println("6")
				windows = append(windows, window)
			}
		case EventDestroyWindows:
			if windows == nil {
				return false
			}

			for _, window := range windows {
				window.Destroy()
			}

			deletedCh <- struct{}{}
		}

		return true
	})

	// Send event to initially create windows.
	app.SendEvent(app, core.NewQEvent(EventCreateWindows))
	app.Exec()

	// TODO(elliot): Move all of this out of main, use resolver, use main for signal handling.
}
