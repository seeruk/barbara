package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
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

	// Listen to events and just dump them to standard out.
	// A more involved approach will have to read the 'U' field of
	// RandrNotifyEvent, which is a union (really a struct) of type
	// RanrNotifyDataUnion.
	for {
		ev, err := xc.WaitForEvent()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ev)
	}
}
