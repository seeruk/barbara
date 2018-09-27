package clock

import (
	"context"
	"time"

	"github.com/seeruk/barbara/barbara"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// Module is a Barbara Module that presents a clock. It uses Go's time formatting, and is basically
// just a label that gets updated every second.
type Module struct {
	ctx context.Context
	cfn context.CancelFunc

	label  *widgets.QLabel
	parent widgets.QWidget_ITF
}

// NewModule returns a new Module instance.
func NewModule(parent widgets.QWidget_ITF) *Module {
	return &Module{
		parent: parent,
	}
}

// Render attempts starts a background process to update the time displayed in a label that is then
// returned to be placed on a bar.
func (m *Module) Render(_ barbara.Alignment, _ barbara.Position) (widgets.QWidget_ITF, error) {
	m.label = widgets.NewQLabel2(time.Now().Format("15:04:05\nMon, 02 Jan"), m.parent, core.Qt__Widget)
	m.label.SetAlignment(core.Qt__AlignCenter)

	m.ctx, m.cfn = context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-m.ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				m.label.SetText(time.Now().Format("15:04:05\nMon, 02 Jan"))
			}
		}
	}()

	return m.label, nil
}

// Destroy stops background processes and frees up resources.
func (m *Module) Destroy() error {
	if m.cfn != nil {
		m.cfn()
	}

	if m.label != nil {
		m.label.Destroy(true, true)
	}

	m.ctx = nil
	m.cfn = nil
	m.label = nil

	return nil
}
