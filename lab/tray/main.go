package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
)

const (
	NotifierHostService            = "org.kde.StatusNotifierHost-%d"
	NotifierItemPath               = "/StatusNotifierItem"
	NotifierItemInterface          = "org.kde.StatusNotifierItem"
	NotifierWatcherPath            = "/StatusNotifierWatcher"
	NotifierWatcherService         = "org.kde.StatusNotifierWatcher"
	NotifierWatcherRegisterHost    = "RegisterStatusNotifierHost"
	NotifierIntrospectionInterface = "org.freedesktop.DBus.Introspectable"
	PropertiesInterface            = "org.freedesktop.DBus.Properties"
)

func senderAndPath(serviceName string, sender dbus.Sender) (string, dbus.ObjectPath) {
	if regexp.MustCompile("^(/\\w+)+$").MatchString(serviceName) {
		return string(sender), dbus.ObjectPath(serviceName)
	} else {
		return string(serviceName), dbus.ObjectPath(NotifierItemPath)
	}
}

type RegisterStatusNotifierItemInput struct {
	Service string
	Sender  dbus.Sender
}

func main() {
	log.Println("Tray starting")

	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	// Add introspection, so that other DBUS-using applications can see which methods we supposedly
	// support, the XML should probably be generated from a struct or something.
	conn.Export(introspect.Introspectable(IntrospectXML), NotifierWatcherPath, NotifierIntrospectionInterface)

	// This channel needs to be buffered, otherwise we'll end up stalling item registration, and
	// potentially deadlocking.
	registrationCh := make(chan RegisterStatusNotifierItemInput, 128)

	go func() {
		for {
			select {
			case registration := <-registrationCh:
				continue

				fmt.Println("\nITEM REGISTERING?")
				fmt.Println("service", registration.Service)
				fmt.Printf("sender %+v\n", registration.Sender)

				method := fmt.Sprintf("%s.Get", PropertiesInterface)
				fmt.Println(method)

				es, ep := senderAndPath(registration.Service, registration.Sender)
				fmt.Println(es, ep)

				itemObj := conn.Object(es, ep)
				call := itemObj.Call(method, dbus.Flags(64), NotifierItemInterface, "Title")
				if call.Err != nil {
					return
				}

				title, _ := itemObj.GetProperty(fmt.Sprintf("%s.Title", NotifierItemInterface))
				status, _ := itemObj.GetProperty(fmt.Sprintf("%s.Status", NotifierItemInterface))
				windowID, _ := itemObj.GetProperty(fmt.Sprintf("%s.WindowId", NotifierItemInterface))
				iconName, _ := itemObj.GetProperty(fmt.Sprintf("%s.IconName", NotifierItemInterface))
				iconPixmap, _ := itemObj.GetProperty(fmt.Sprintf("%s.IconPixmap", NotifierItemInterface))
				overlayIconName, _ := itemObj.GetProperty(fmt.Sprintf("%s.OverlayIconName", NotifierItemInterface))
				overlayIconPixmap, _ := itemObj.GetProperty(fmt.Sprintf("%s.OverlayIconPixmap", NotifierItemInterface))
				attentionIconName, _ := itemObj.GetProperty(fmt.Sprintf("%s.AttentionIconName", NotifierItemInterface))
				attentionIconPixmap, _ := itemObj.GetProperty(fmt.Sprintf("%s.AttentionIconPixmap", NotifierItemInterface))
				attentionMovieName, _ := itemObj.GetProperty(fmt.Sprintf("%s.AttentionMovieName", NotifierItemInterface))
				tooltip, _ := itemObj.GetProperty(fmt.Sprintf("%s.ToolTip", NotifierItemInterface))
				itemIsMenu, _ := itemObj.GetProperty(fmt.Sprintf("%s.ItemIsMenu", NotifierItemInterface))
				menu, _ := itemObj.GetProperty(fmt.Sprintf("%s.Menu", NotifierItemInterface))

				fmt.Println("title", title.String())
				fmt.Println("status", status.String())
				if windowID.Value() != nil {
					fmt.Println("windowID", windowID.String())
				}
				if iconName.Value() != nil {
					fmt.Println("iconName", iconName.String())
				}
				if iconPixmap.Value() != nil {
					fmt.Printf("iconPixmap: %v\n", iconPixmap.Value())
				}
				if overlayIconName.Value() != nil {
					fmt.Println("overlayIconName", overlayIconName.String())
				}
				if overlayIconPixmap.Value() != nil {
					fmt.Printf("overlayIconPixmap: %v\n", overlayIconPixmap.Value())
				}
				if attentionIconName.Value() != nil {
					fmt.Println("attentionIconName", attentionIconName.String())
				}
				if attentionIconPixmap.Value() != nil {
					fmt.Printf("attentionIconPixmap: %v\n", attentionIconPixmap.Value())
				}
				if attentionMovieName.Value() != nil {
					fmt.Println("attentionMovieName", attentionMovieName.String())
				}
				if tooltip.Value() != nil {
					fmt.Println("tooltip", tooltip.String())
				}
				if itemIsMenu.Value() != nil {
					fmt.Println("itemIsMenu", itemIsMenu.String())
				}
				if menu.Value() != nil {
					fmt.Println("menu", menu.String())
				}

				fmt.Println()
			}
		}
	}()

	// Add watcher methods.
	conn.ExportMethodTable(
		map[string]interface{}{
			"RegisterStatusNotifierItem": func(service string, sender dbus.Sender) *dbus.Error {
				registrationCh <- RegisterStatusNotifierItemInput{
					Service: service,
					Sender:  sender,
				}

				return nil
			},
			"RegisterStatusNotifierHost": func(service string, sender dbus.Sender) *dbus.Error {
				fmt.Println("HOST REGISTERING?", service)

				// Emitting signals, maybe we should emit both a freedesktop one and a kde one?
				// Need to test that this works, by starting something that listens for signals.
				err := conn.Emit(NotifierWatcherPath, "org.kde.StatusNotifierWatcher.StatusNotifierHostRegistered")
				if err != nil {
					fmt.Println("SIGNAL ERROR", err)
				}

				return nil
			},
		},
		NotifierWatcherPath,
		NotifierWatcherService,
	)

	var hostRegistered bool

	refreshPops := func() {
		// Properties on the bus are exposed like this. It's a function, because it needs to be
		// pushed again to the bus when anything changes... all at once? (Maybe a diff can be sent,
		// but that doesn't really matter, a full update is just as easy).
		prop.New(conn, NotifierWatcherPath, map[string]map[string]*prop.Prop{
			NotifierWatcherService: {
				"IsStatusNotifierHostRegistered": {hostRegistered, false, prop.EmitTrue, nil},
				"ProtocolVersion":                {1, false, prop.EmitTrue, nil},
				"RegisteredStatusNotifierItems":  {[]string{}, false, prop.EmitTrue, nil},
			},
		})
	}

	// Set initial properties state.
	refreshPops()

	// This might not be how these things are meant to be handled... maybe signals are things that
	// are propagated _from_ the service, not something
	c := make(chan *dbus.Signal, 16)

	// This actually registers the channel, so conn.RemoveSignal should be called.
	conn.Signal(c)

	conn.BusObject().Call("org.freedesktop.StatusNotifierWatcher.StatusNotifierHostRegistered", 0, "")

	go func() {
		for v := range c {
			fmt.Println("SIGNAL RECEIVED")
			spew.Dump(v)

			refreshPops()
		}
	}()

	// Request the name of the watcher service, and do not queue up to wait to become the primary
	// owner of the name. Messages get sent to this service if it is the primary.
	reply, err := conn.RequestName(NotifierWatcherService, dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}

	// If we aren't the primary, we want to bail.
	if reply != dbus.RequestNameReplyPrimaryOwner {
		panic(fmt.Errorf("service %s already taken", NotifierWatcherService))
	}

	fmt.Println("IS PRIMARY WATCHER, CONTINUING")

	host := fmt.Sprintf(NotifierHostService, os.Getpid())

	// Request the name of the host service. The host must only be present on the bus, as primary.
	// After this, the host should be sent items by the watcher. The presence of both the watcher
	// and the host are what allows items to be retrieved (I think?)
	//
	// Or is a host something that "owns" an item to be displayed?
	reply, err = conn.RequestName(host, dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}

	// If this host is already registered (but... how?), then bail.
	if reply != dbus.RequestNameReplyPrimaryOwner {
		panic(fmt.Errorf("service %s already taken", host))
	}

	go func() {
		signals := make(chan *dbus.Signal, 100)
		conn.Signal(signals)

		for signal := range signals {
			fmt.Printf("SIGNAL: %+v\n", signal)
		}
	}()

	hostObj := conn.Object(NotifierWatcherService, NotifierWatcherPath)

	// Register our host in the watcher.
	call := hostObj.Call(NotifierWatcherRegisterHost, 0, host)
	if call.Err != nil {
		panic(call.Err)
	}

	fmt.Printf("Call: %+v\n", call)

	time.Sleep(5 * time.Hour)
}
