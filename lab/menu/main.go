package main

import (
	"fmt"
	"log"

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

	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	button := createButton(createMenu(createMenuItem(), createImageMenuItem()))

	win.Add(button)
	win.ShowAll()

	gtk.Main()
}

func createButton(menu *gtk.Menu) *gtk.Button {
	// Created the button to show the menu.
	button, err := gtk.ButtonNewWithLabel("Show Menu")
	if err != nil {
		panic(err)
	}

	button.Connect("clicked", func(btn *gtk.Button) {
		menu.PopupAtPointer(nil)
	})

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

func createImageMenuItem() *gtk.MenuItem {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		panic(err)
	}

	icon, err := gtk.ImageNewFromIconName("blueman-tray", gtk.ICON_SIZE_MENU)
	if err != nil {
		panic(err)
	}

	label, err := gtk.LabelNewWithMnemonic("_Image Item")
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

	return menuItem
}

func createMenuItem() *gtk.MenuItem {
	// Mnemonics allow us to use the _Label notation to allow us to press the key following an _ to
	// activate the menu item.
	menuItem, err := gtk.MenuItemNewWithMnemonic("_Test Item")
	if err != nil {
		panic(err)
	}

	// Used for shortcut keys later...
	acgroup, err := gtk.AccelGroupNew()
	if err != nil {
		panic(err)
	}

	// Show shortcuts on menu items, something like this (not sure actually what an AccelGroup is...
	menuItem.AddAccelerator("activate", acgroup, gdk.KEY_f, gdk.GDK_CONTROL_MASK, gtk.ACCEL_VISIBLE)

	// When the item is clicked, what should happen?
	menuItem.Connect("activate", func() {
		fmt.Println("clicked menu item")
	})

	return menuItem
}
