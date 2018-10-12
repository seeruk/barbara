package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/seeruk/barbara/event"
	"github.com/seeruk/barbara/internal"
)

func main() {
	log.Println("Started...")

	config, err := internal.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	timeout := make(chan struct{}, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	resolver := internal.NewResolver(config)

	watcher := resolver.ResolveX11RandrEventWatcher()
	watcher.Watch(context.Background())

	app := resolver.ResolveApplication()

	dispatcher := resolver.ResolveEventDispatcher()
	dispatcher.Dispatch(event.TypeStartup)

	go func() {
		sig := <-signals
		timeout <- struct{}{}

		fmt.Println() // Skip the ^C
		log.Printf("Caught %v. Shutting down...\n", sig)

		dispatcher.Dispatch(event.TypeShutdown)
	}()

	go func() {
		<-timeout
		time.AfterFunc(5*time.Second, func() {
			log.Println("Took too long shutting down. Exiting...")
			os.Exit(1)
		})
	}()

	app.QApplication().Exec() // Block until signal is caught.
}
