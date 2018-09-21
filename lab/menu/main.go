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

	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	button := createButton(createMenu(createMenuItem()))

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
