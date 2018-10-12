package x11

import (
	"context"
	"log"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/seeruk/barbara/event"
)

// RandrEventWatcher watches for randr events in X, allowing other parts of the application to react
// to randr events (e.g. for re-rendering bars).
type RandrEventWatcher struct {
	dispatcher *event.Dispatcher
	xc         *xgb.Conn
}

// NewRandrEventWatcher returns a new RandrEventWatcher instance.
func NewRandrEventWatcher(dispatcher *event.Dispatcher, xc *xgb.Conn) *RandrEventWatcher {
	return &RandrEventWatcher{
		dispatcher: dispatcher,
		xc:         xc,
	}
}

// Watch starts watching for randr events from the X server. It will debounce events, waiting for 1
// second of no event activity before sending a WM event in the event dispatcher.
func (w *RandrEventWatcher) Watch(ctx context.Context) {
	debounceCh := make(chan struct{}, 128)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			_, err := w.xc.WaitForEvent()
			if err != nil {
				// TODO(elliot): What to do...
				log.Fatal(err)
			}

			debounceCh <- struct{}{}
		}
	}()

	go func() {
		var timerCh <-chan time.Time

		for {
			select {
			case <-ctx.Done():
				return
			case <-debounceCh:
				timerCh = time.After(time.Second)
			case <-timerCh:
				// Dispatch WM event.
				w.dispatcher.Dispatch(event.TypeWM)
			}
		}
	}()
}
