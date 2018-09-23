package main

import (
	"fmt"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetStyleSheet(`
		QMainWindow {
			background: #1a1a1a;
			margin: 0;
			padding: 0px;
		}

		.board-button {
			background: #1a1a1a;
			color: #e5e5e5;
			padding: 10px;
		}
	`)

	screens := app.Screens()

	for i, screen := range screens {
		geo := screen.Geometry()

		fmt.Printf("Screen: %d\n", i)
		fmt.Printf("Name: %s\n", screen.Name())
		fmt.Printf("X: %d\n", geo.X())
		fmt.Printf("Y: %d\n", geo.Y())
		fmt.Printf("Width: %d\n", geo.Width())
		fmt.Printf("Height: %d\n", geo.Height())
		fmt.Println()

		lbox := widgets.NewQHBoxLayout()
		rbox := widgets.NewQHBoxLayout()

		cbox := widgets.NewQHBoxLayout()
		cbox.AddLayout(lbox, 1)
		cbox.AddLayout(rbox, 1)

		cbox.SetContentsMargins(0, 0, 0, 0)

		window := widgets.NewQMainWindow(nil, core.Qt__Drawer)
		window.SetWindowTitle("Board Example")
		window.SetMaximumHeight(54)

		window.StatusBar().Hide()

		// Turn into dock, and move to bottom of screen.
		window.SetAttribute(core.Qt__WA_X11NetWmWindowTypeDock, true)
		window.Move2(geo.X(), geo.Height()-window.Height())

		// Volume Label Start
		volumeST := widgets.NewQAction2("Volume:", nil)
		volumeST.SetEnabled(false)
		// Volume Label End

		// Volume Slider Start
		slider := widgets.NewQSlider(nil)
		slider.SetOrientation(core.Qt__Horizontal)
		slider.SetTickInterval(100)

		sliderBox := widgets.NewQHBoxLayout()
		sliderBox.AddWidget(slider, 0, core.Qt__AlignJustify)
		sliderBox.SetContentsMargins(26, 7, 26, 7)

		sliderWid := widgets.NewQWidget(nil, 0)
		sliderWid.SetLayout(sliderBox)

		sliderAction := widgets.NewQWidgetAction(nil)
		sliderAction.SetDefaultWidget(sliderWid)
		// Volume Slider End

		menu := widgets.NewQMenu(nil)
		menu.AddActions([]*widgets.QAction{
			volumeST,
			menu.AddSeparator(),
			sliderAction.QAction_PTR(),
		})

		button := widgets.NewQPushButton2("Volume", nil)
		button.SetProperty("class", core.NewQVariant14("board-button"))
		//button.SetMenu(menu)
		button.ConnectClicked(func(checked bool) {
			mw := menu.SizeHint().Width()
			mh := menu.SizeHint().Height()

			x := button.Width() - mw
			y := -mh

			point := button.MapToGlobal(core.NewQPoint2(x, y))

			menu.Popup(point, nil)
		})

		window.SetMaximumHeight(button.SizeHint().Height())
		window.SetFixedHeight(button.SizeHint().Height())

		rbox.AddWidget(button, 0, core.Qt__AlignRight)

		cwid := widgets.NewQWidget(nil, 0)
		cwid.SetLayout(cbox)

		window.SetCentralWidget(cwid)
		window.Show()
	}

	app.Exec()
}
