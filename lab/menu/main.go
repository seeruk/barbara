package main

import (
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var mainWin *gtk.Window

func main() {
	gtk.Init(nil)

	var err error

	mainWin, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	// This'll fling it to the top of the screen by default probably, at least on i3.
	mainWin.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)
	mainWin.SetDecorated(false)

	// Use dark theme.
	settings, _ := gtk.SettingsGetDefault()
	settings.SetProperty("gtk-application-prefer-dark-theme", true)

	mainWin.SetTitle("Simple Example")
	mainWin.Connect("destroy", func() {
		gtk.MainQuit()
	})

	barBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		panic(err)
	}

	cssProvider, err := gtk.CssProviderNew()
	if err != nil {
		panic(err)
	}

	cssProvider.LoadFromData(`
		.board-window {
			background: #1a1a1a;
			padding: 7px;
		}

		.board-button { background-color: #1a1a1a; }
		.board-button:hover { background-color: #2a2a2a; }
		.board-button:active { background-color: #2a2a2a; }

		.board-datetime { padding: 0 7px; }
	`)

	styles, err := barBox.GetStyleContext()
	if err != nil {
		panic(err)
	}

	styles.AddClass("board-window")
	styles.AddProvider(cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	leftBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		panic(err)
	}

	rightBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		panic(err)
	}

	barBox.SetHAlign(gtk.ALIGN_FILL)
	barBox.SetHExpand(true)

	leftBox.SetHAlign(gtk.ALIGN_START)
	leftBox.SetHExpand(true)

	rightBox.SetHAlign(gtk.ALIGN_END)
	rightBox.SetHExpand(true)

	barBox.Add(leftBox)
	barBox.Add(rightBox)

	ws1, _ := gtk.ButtonNewWithLabel("1")

	styles, _ = ws1.GetStyleContext()
	styles.AddClass("board-button")
	styles.AddProvider(cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	ws7, _ := gtk.ButtonNewWithLabel("7")

	styles, _ = ws7.GetStyleContext()
	styles.AddClass("board-button")
	styles.AddProvider(cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	leftBox.Add(ws1)
	leftBox.Add(ws7)

	rightLbl, _ := gtk.LabelNew(time.Now().Format("Monday, 02 Jan - 15:04:05"))

	styles, _ = rightLbl.GetStyleContext()
	styles.AddClass("board-datetime")
	styles.AddProvider(cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-ticker.C:
				// Add this to the main loop, from this thread.
				glib.IdleAdd(func(label *gtk.Label) bool {
					label.SetLabel(time.Now().Format("Monday, 02 Jan - 15:04:05"))
					return false
				}, rightLbl)
			}
		}
	}()

	button := createUserMenuButton(createUserMenu(
		createLogOffMenuItem(),
		createRebootMenuItem(),
		createShutdownMenuItem(),
	))

	rightBox.Add(rightLbl)
	rightBox.Add(createSoundMenuButton())
	rightBox.Add(button)

	// We have one display per X session normally. Technically we could have more than one screen,
	// but it's very unlikely - so unlikely that a lot of that functionality is deprecated in GTK.
	display, _ := gdk.DisplayGetDefault()

	// Each display can have multiple monitors. There should be a primary monitor.
	monitor, _ := display.GetPrimaryMonitor()

	// The geometry is x and y position, and width and height.
	monGeo := monitor.GetGeometry()

	// So for now, don't affect x position, just move it down the full height of the screen, minus
	// the height of the window, so it's positioned at the bottom of the primary display.
	mainWin.Move(0, monGeo.GetHeight()-mainWin.GetAllocation().GetHeight())

	mainWin.Add(barBox)
	mainWin.ShowAll()

	gtk.Main()
}

func createSoundMenuButton() *gtk.Button {
	cssProvider, _ := gtk.CssProviderNew()
	cssProvider.LoadFromData(`
		.board-button { background-color: #1a1a1a; }
		.board-button:hover { background-color: #2a2a2a; }
		.board-button:active { background-color: #2a2a2a; }
	`)

	var open bool

	icon, err := gtk.ImageNewFromIconName("audio-volume-muted", gtk.ICON_SIZE_BUTTON)
	if err != nil {
		panic(err)
	}

	button, _ := gtk.ButtonNew()
	button.Add(icon)

	win, _ := gtk.WindowNew(gtk.WINDOW_POPUP)
	win.SetDecorated(false)
	win.SetKeepAbove(true)
	win.SetPosition(gtk.WIN_POS_NONE)
	win.SetResizable(false)
	win.SetSkipPagerHint(true)
	win.SetSkipTaskbarHint(true)
	win.SetTitle("board volume")
	win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DROPDOWN_MENU)
	win.Stick()

	cssProvider2, _ := gtk.CssProviderNew()
	cssProvider2.LoadFromData(`
		scale {
			border: 1px solid #333;
			min-width: 100px;
			padding: 7px 10px;
		}
	`)

	scale, _ := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 1, 0.02)
	scale.SetDrawValue(false)
	scale.ShowAll()

	scale.Connect("value-changed", func() {
		val := scale.GetValue()

		switch {
		case val > 0.66:
			icon.SetFromIconName("audio-volume-high", gtk.ICON_SIZE_BUTTON)
		case val > 0.33:
			icon.SetFromIconName("audio-volume-medium", gtk.ICON_SIZE_BUTTON)
		case val > 0:
			icon.SetFromIconName("audio-volume-low", gtk.ICON_SIZE_BUTTON)
		case val == 0:
			icon.SetFromIconName("audio-volume-muted", gtk.ICON_SIZE_BUTTON)
		}
	})

	styles2, _ := scale.GetStyleContext()
	styles2.AddProvider(cssProvider2, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	win.Add(scale)
	win.SetAcceptFocus(true)
	win.SetCanFocus(true)
	win.Widget.SetEvents(int(gdk.FOCUS_CHANGE_MASK))

	win.Connect("focus-in-event", func() {
		fmt.Println("GOT FOCUS")
	})

	win.Connect("focus-out-event", func() {
		win.Hide()
		open = false
	})

	button.Connect("clicked", func(btn *gtk.Button) {
		if open {
			win.Hide()
			open = false
		} else {
			win.ShowAll()

			btnAlloc := btn.GetAllocation()
			btnWin, _ := btn.GetWindow()

			_, wy := btnWin.GetRootOrigin()
			_, vwh := win.GetSize()

			spew.Dump(btnAlloc)
			spew.Dump(wy, vwh)

			// Move over button.
			win.Move(
				btnAlloc.GetX()+btnAlloc.GetWidth()-122,
				wy+btnAlloc.GetY()-vwh,
			)

			// Move under button.
			//win.Move(
			//	btnAlloc.GetX()+btnAlloc.GetWidth()-122,
			//	wy+btnAlloc.GetY()+btnAlloc.GetHeight(),
			//)

			open = true
		}

	})

	styles, _ := button.GetStyleContext()
	styles.AddClass("board-button")
	styles.AddProvider(cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return button
}

func createUserMenuButton(menu *gtk.Menu) *gtk.Button {
	cssProvider, _ := gtk.CssProviderNew()
	cssProvider.LoadFromData(`
		.board-button { background-color: #1a1a1a; }
		.board-button:hover { background-color: #2a2a2a; }
		.board-button:active { background-color: #2a2a2a; }
	`)

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	name := usr.Name
	if name == "" {
		name = usr.Username
	}

	// Created the button to show the menu.
	button, _ := gtk.ButtonNewWithLabel(name)
	button.Connect("clicked", func(btn *gtk.Button) {
		menu.PopupAtWidget(btn, gdk.GDK_GRAVITY_NORTH_EAST, gdk.GDK_GRAVITY_SOUTH_EAST, nil)
	})

	button.SetName("board-user-menu")

	styles, _ := button.GetStyleContext()
	styles.AddClass("board-button")
	styles.AddProvider(cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return button
}

func createUserMenu(items ...*gtk.MenuItem) *gtk.Menu {
	scale, _ := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 1, 0.02)
	scale.SetDrawValue(false)
	scale.SetCanFocus(true)
	scale.SetSensitive(true)

	menuItem, _ := gtk.MenuItemNew()
	menuItem.Add(scale)
	menuItem.SetCanFocus(false)
	menuItem.SetFocusChild(scale)

	menu, err := gtk.MenuNew()
	if err != nil {
		panic(err)
	}

	for _, item := range items {
		menu.Append(item)
	}

	menu.Append(menuItem)

	// If this is not called, the menu will be empty.
	menu.ShowAll()

	return menu
}

func createLogOffMenuItem() *gtk.MenuItem {
	return createMenuItem("_Log Off", "system-log-out", func() {
		err := exec.Command("i3-msg", "exit").Start()
		if err != nil {
			log.Println(err)
		}
	})
}

func createRebootMenuItem() *gtk.MenuItem {
	return createMenuItem("_Reboot", "system-reboot", func() {
		err := exec.Command("sudo", "systemctl", "reboot").Start()
		if err != nil {
			log.Println(err)
		}
	})
}

func createShutdownMenuItem() *gtk.MenuItem {
	return createMenuItem("_Shutdown", "system-shutdown", func() {
		err := exec.Command("sudo", "systemctl", "poweroff").Start()
		if err != nil {
			log.Println(err)
		}
	})
}

func createMenuItem(mnemonic, iconName string, activate func()) *gtk.MenuItem {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		panic(err)
	}

	icon, err := gtk.ImageNewFromIconName(iconName, gtk.ICON_SIZE_MENU)
	if err != nil {
		panic(err)
	}

	label, err := gtk.LabelNewWithMnemonic(mnemonic)
	if err != nil {
		panic(err)
	}

	box.Add(icon)
	box.Add(label)

	menuItem, err := gtk.MenuItemNew()
	if err != nil {
		panic(err)
	}

	menuItem.Add(box)
	menuItem.ShowAll()

	menuItem.Connect("activate", activate)

	return menuItem
}
