package x11

import (
	"context"
	"log"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/seeruk/barbara/event"
)

// RandrEventWatcher ...
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

// Watch ...
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
