package battery

import (
	"context"
	"sync"
	"time"
)

// InfoNotifier is a type used to propagate battery information to types that want to be notified of
// new battery status information.
// TODO(elliot): Use file notification, e.g. inotify, instead of time intervals - of course!
type InfoNotifier struct {
	// TODO(elliot): Logger.
	// TODO(elliot): Do we share a filesystem watcher too? Maybe not.
	reader *InfoReader

	ctx context.Context
	cfn context.CancelFunc

	cs   []chan<- Info
	csMu *sync.Mutex
	ps   string
}

// NewInfoNotifier returns a new InfoNotifier instance.
func NewInfoNotifier(powerSupply string) *InfoNotifier {
	return &InfoNotifier{
		csMu: &sync.Mutex{},
		ps:   powerSupplyPath + "/" + powerSupply,
	}
}

// Notify adds another channel that should be notified when new battery information is available. It
// will be given the updated Info.
func (n *InfoNotifier) Notify(c chan<- Info) {
	n.csMu.Lock()
	defer n.csMu.Unlock()

	n.cs = append(n.cs, c)
}

// Start begins a background process
func (n *InfoNotifier) Start() {
	n.ctx, n.cfn = context.WithCancel(context.Background())

	// TODO(elliot): Configurable interval.
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			// TODO(elliot): Case for status update from fs notifications.
			select {
			case <-n.ctx.Done():
				// TODO(elliot): Probably log here.
				return
			case <-ticker.C:
				n.doNotify()
			}
		}
	}()
}

// Stop attempts to stop the background processes started by this InfoNotifier.
func (n *InfoNotifier) Stop() {
	if n.ctx == nil || n.cfn == nil {
		return
	}

	n.cfn()
	n.ctx = nil
	n.cfn = nil
}

// doNotify actually propagates the battery information to "listener" channels.
func (n *InfoNotifier) doNotify() {
	n.csMu.Lock()
	defer n.csMu.Unlock()

	// TODO(elliot): Read info...
	info := Info{}

	// Notify all listening channels.
	for _, c := range n.cs {
		c <- info
	}
}

// InfoNotifierFactory is a type that keeps track of the instantiated battery InfoNotifier types. If
// the something wants to be notified about battery info for a power supply that's already had an
// InfoNotifier created for it, then we will re-use that existing notifier as it's more efficient.
type InfoNotifierFactory struct {
	infoNotifiers map[string]*InfoNotifier
}

// NewInfoNotifierFactory returns a new InfoNotifierFactory instance.
func NewInfoNotifierFactory() *InfoNotifierFactory {
	return &InfoNotifierFactory{
		infoNotifiers: make(map[string]*InfoNotifier),
	}
}

// Build will create or re-use InfoNotifier instances. New ones will be created if the power supply
// passed in has not already had an InfoNotifier created for it before.
//
// TODO(elliot): How do we destroy these if they're no longer used? Maybe we can hook into events to
// wipe out the map, and stop background processes in the Notifiers? This is a pretty critical one.
func (f *InfoNotifierFactory) Build(powerSupply string) *InfoNotifier {
	if _, ok := f.infoNotifiers[powerSupply]; !ok {
		f.infoNotifiers[powerSupply] = NewInfoNotifier(powerSupply)
	}

	return f.infoNotifiers[powerSupply]
}
