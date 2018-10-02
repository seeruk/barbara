package x11

import (
	"context"
	"log"
	"time"

	"github.com/BurntSushi/xgb"
)

// RandrEventWatcher ...
type RandrEventWatcher struct {
	xc *xgb.Conn
}

// NewRandrEventWatcher returns a new RandrEventWatcher instance.
func NewRandrEventWatcher(xc *xgb.Conn) *RandrEventWatcher {
	return &RandrEventWatcher{
		xc: xc,
	}
}

// Start ...
func (w *RandrEventWatcher) Start(ctx context.Context) chan struct{} {
	eventCh := make(chan struct{}, 32)
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
				eventCh <- struct{}{}
			}
		}
	}()

	return eventCh
}
