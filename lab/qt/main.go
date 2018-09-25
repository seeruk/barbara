package main

import (
	"fmt"
	"os"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
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

		QLabel {
			color: #e5e5e5;
			font-family: "Fira Sans";
			font-size: 13px;
			padding: 0 0 0 7px;
			text-align: center;
		}

		.board-button {
			background-color: #1a1a1a;
			color: #e5e5e5;
			font-family: "Fira Sans";
			font-size: 13px;
			padding: 7px;
		}

		.board-button:flat {
			border: 1px solid #333;
			border-radius: 3px;
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
		lbox.SetAlign(core.Qt__AlignLeft)

		rbox := widgets.NewQHBoxLayout()
		rbox.SetAlign(core.Qt__AlignRight)

		cbox := widgets.NewQHBoxLayout()
		cbox.AddLayout(lbox, 1)
		cbox.AddLayout(rbox, 1)

		cbox.SetContentsMargins(7, 7, 7, 7)

		window := widgets.NewQMainWindow(nil, core.Qt__Window)
		window.SetWindowTitle("Board Example")
		window.StatusBar().Hide()

		// Turn into dock, and move to bottom of screen.
		window.SetAttribute(core.Qt__WA_X11NetWmWindowTypeDock, true)
		window.Move2(geo.X(), geo.Height()-window.Height())

		dateLabel := widgets.NewQLabel2(time.Now().Format("15:04:05\nMon, 02 Jan"), nil, core.Qt__Widget)
		dateLabel.SetAlignment(core.Qt__AlignCenter)

		go func() {
			ticker := time.NewTicker(time.Second)

			for {
				select {
				case <-ticker.C:
					dateLabel.SetText(time.Now().Format("15:04:05\nMon, 02 Jan"))
				}
			}
		}()

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

		icon := gui.NewQIcon5("/usr/share/icons/Paper/512x512/status/audio-volume-muted.png")

		button := widgets.NewQPushButton3(icon, "", nil)
		button.SetFlat(true)
		button.SetProperty("class", core.NewQVariant14("board-button"))
		//button.SetMenu(menu) // Can't position menu if I use this.
		button.ConnectClicked(func(checked bool) {
			// But I can move the menu when it's clicked, based on the button position, etc.
			mw := menu.SizeHint().Width()
			mh := menu.SizeHint().Height()

			x := button.Width() - mw
			y := -mh

			point := button.MapToGlobal(core.NewQPoint2(x, y))

			menu.Popup(point, nil)
		})

		slider.ConnectValueChanged(func(value int) {
			switch {
			case value > 66:
				icon := gui.NewQIcon5("/usr/share/icons/Paper/512x512/status/audio-volume-high.png")
				button.SetIcon(icon)
			case value > 33:
				icon := gui.NewQIcon5("/usr/share/icons/Paper/512x512/status/audio-volume-medium.png")
				button.SetIcon(icon)
			case value > 0:
				icon := gui.NewQIcon5("/usr/share/icons/Paper/512x512/status/audio-volume-low.png")
				button.SetIcon(icon)
			case value == 0:
				icon := gui.NewQIcon5("/usr/share/icons/Paper/512x512/status/audio-volume-muted.png")
				button.SetIcon(icon)
			}
		})

		window.SetFixedHeight(button.SizeHint().Height() + 14)

		rbox.AddWidget(button, 0, core.Qt__AlignRight)

		userButton := widgets.NewQPushButton2("Elliot Wright", nil)
		userButton.SetFlat(true)
		userButton.SetProperty("class", core.NewQVariant14("board-button"))

		rbox.AddWidget(userButton, 0, core.Qt__AlignRight)

		rbox.AddWidget(dateLabel, 0, core.Qt__AlignRight)

		cwid := widgets.NewQWidget(nil, 0)
		cwid.SetLayout(cbox)

		window.SetCentralWidget(cwid)
		window.Show()
	}

	app.Exec()
}
