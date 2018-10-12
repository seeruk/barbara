package internal

import (
	"fmt"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/seeruk/barbara/barbara"
	"github.com/seeruk/barbara/event"
	"github.com/seeruk/barbara/modules/clock"
	"github.com/seeruk/barbara/modules/menu"
	"github.com/seeruk/barbara/wm/x11"
)

// Resolver is a type that resolves Barbara's runtime dependencies. It handles wiring up types in
// the application, using plain Go.
type Resolver struct {
	config Config

	app        *barbara.Application
	dispatcher *event.Dispatcher
	xc         *xgb.Conn
}

// NewResolver returns a new instance of Resolver.
func NewResolver(config Config) *Resolver {
	resolver := &Resolver{config: config}
	resolver.resolveEager()

	return resolver
}

// ResolveApplication resolves the Application instance.
func (r *Resolver) ResolveApplication() *barbara.Application {
	if r.app == nil {
		r.app = barbara.NewApplication(
			r.ResolveModuleFactory(),
			r.config.Primary,
			r.config.Secondary,
		)

		// Register application events in dispatcher.
		dispatcher := r.ResolveEventDispatcher()
		dispatcher.RegisterListener(event.TypeStartup, r.app.CreateWindows)
		dispatcher.RegisterListener(event.TypeShutdown, r.app.DestroyWindows)
		dispatcher.RegisterListener(event.TypeWM, r.app.RecreateWindows)
	}

	return r.app
}

// ResolveEventDispatcher resolves the application's event dispatcher.
func (r *Resolver) ResolveEventDispatcher() *event.Dispatcher {
	if r.dispatcher == nil {
		r.dispatcher = event.NewDispatcher()
	}

	return r.dispatcher
}

// ResolveModuleFactory resolves a new barbara.ModuleFactory instance, with available modules
// already registered with it.
func (r *Resolver) ResolveModuleFactory() *barbara.ModuleFactory {
	// Register all available modules, modules might have special constructors, so this approach
	// needs to be taken over an approach similar to sql.DB drivers. Modules may have dependencies
	// on shared services (e.g. some kind of API client, for example).
	mbf := barbara.NewModuleFactory()
	mbf.RegisterConstructor("clock", clock.NewModule)
	mbf.RegisterConstructor("menu", menu.NewModule)

	return mbf
}

// ResolveXConnection resolves the application's X connection, setting up extensions, etc.
func (r *Resolver) ResolveXConnection() *xgb.Conn {
	if r.xc == nil {
		xc, _ := xgb.NewConn()

		// Every extension must be initialized before it can be used.
		err := randr.Init(xc)
		if err != nil {
			panic(fmt.Errorf("failed to initalise randr: %v", err))
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
			panic(fmt.Errorf("failed to subscribe to randr events: %v", err))
		}

		r.xc = xc
	}

	return r.xc
}

// ResolveX11RandrEventWatcher resolves a new x11.RandrEventWatcher instance.
func (r *Resolver) ResolveX11RandrEventWatcher() *x11.RandrEventWatcher {
	return x11.NewRandrEventWatcher(
		r.ResolveEventDispatcher(),
		r.ResolveXConnection(),
	)
}

// resolveEager resolves the types that may error, so that an errors can be found at startup,
// instead of later on when the application may have started. It also resolves types that should be
// resolved at startup in general, as their resolvers may be used to do things like register events
// listeners in the event dispatcher.
func (r *Resolver) resolveEager() {
	r.ResolveApplication()
	r.ResolveXConnection()
}
