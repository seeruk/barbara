package main

import (
	"context"
	"fmt"
	"log"

	"github.com/seeruk/barbara/event"
	"github.com/seeruk/barbara/internal"
)

func main() {
	fmt.Println("Started...")

	config, err := internal.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	resolver := internal.NewResolver(config)

	watcher := resolver.ResolveX11RandrEventWatcher()
	watcher.Watch(context.Background())

	dispatcher := resolver.ResolveEventDispatcher()
	dispatcher.Dispatch(event.TypeStartup)

	app := resolver.ResolveApplication()
	app.QApplication().Exec()

	// TODO(elliot): Signal handling?
}
