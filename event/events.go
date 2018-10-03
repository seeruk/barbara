package event

import "sync"

const (
	// TypeStartup is sent when Barbara is first started.
	TypeStartup Type = iota
	// TypeShutdown is sent when Barbara is completely shutting down.
	TypeShutdown
	// TypeWM is an event that comes from a window manager.
	TypeWM
)

// Type represents enumerations of event types.
type Type int

// Dispatcher is a very, very basic event dispatcher. It keeps a map between event types, and
// functions to call when an event is dispatched. There are no background processes involved.
type Dispatcher struct {
	sync.Mutex

	listeners map[Type][]func()
}

// NewDispatcher creates a new Dispatcher instance.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		listeners: make(map[Type][]func()),
	}
}

// Dispatch emits an event via this Dispatcher instance, calling any listener functions associated
// with the given event Type. This method is safe for concurrent use.
func (d *Dispatcher) Dispatch(eventType Type) {
	d.Lock()
	defer d.Unlock()

	if lfns, ok := d.listeners[eventType]; ok {
		for _, lfn := range lfns {
			lfn()
		}
	}
}

// RegisterListener is used to register an event listener. An event listener is basically just a
// function that will be called when a given event Type is dispatched. This method is safe for
// concurrent use.
func (d *Dispatcher) RegisterListener(eventType Type, lfn func()) {
	d.Lock()
	defer d.Unlock()

	d.listeners[eventType] = append(d.listeners[eventType], lfn)
}
