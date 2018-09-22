package main

import (
	"log"
	"os/exec"
	"os/user"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	// This'll fling it to the top of the screen by default probably, at least on i3.
	win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)
	win.SetDecorated(false)

	// Put the window at the bottom?
	win.Move(0, 1080)

	// Use dark theme.
	settings, _ := gtk.SettingsGetDefault()
	settings.SetProperty("gtk-application-prefer-dark-theme", true)

	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
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

	midBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
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

	midBox.SetHAlign(gtk.ALIGN_CENTER)
	midBox.SetHExpand(true)

	rightBox.SetHAlign(gtk.ALIGN_END)
	rightBox.SetHExpand(true)

	barBox.Add(leftBox)
	barBox.Add(midBox)
	barBox.Add(rightBox)

	leftLbl, _ := gtk.LabelNew("This is the left")
	leftBox.Add(leftLbl)

	midLbl, _ := gtk.LabelNew("")
	midBox.Add(midLbl)

	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-ticker.C:
				midLbl.SetLabel(time.Now().Format("Monday, 02 Jan - 15:04:05"))
			}
		}
	}()

	button := createUserMenuButton(createMenu(
		createLogOffMenuItem(),
		createRebootMenuItem(),
		createShutdownMenuItem(),
	))

	rightBox.Add(button)

	win.Add(barBox)
	win.ShowAll()

	gtk.Main()
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

	// Created the button to show the menu.
	button, _ := gtk.ButtonNewWithLabel(usr.Name)
	button.Connect("clicked", func(btn *gtk.Button) {
		menu.PopupAtWidget(btn, gdk.GDK_GRAVITY_NORTH_EAST, gdk.GDK_GRAVITY_SOUTH_EAST, nil)
	})

	button.SetName("board-user-menu")

	styles, _ := button.GetStyleContext()
	styles.AddClass("board-button")
	styles.AddProvider(cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return button
}

func createMenu(items ...*gtk.MenuItem) *gtk.Menu {
	menu, err := gtk.MenuNew()
	if err != nil {
		panic(err)
	}

	for _, item := range items {
		menu.Append(item)
	}

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
