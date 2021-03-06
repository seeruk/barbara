package clock

import (
	"context"
	"encoding/json"
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

	config Config
	layout *widgets.QHBoxLayout
	label  *widgets.QLabel
}

// NewModule returns a new clock Module instance.
func NewModule(mctx barbara.ModuleContext) (barbara.Module, error) {
	var config Config

	err := json.Unmarshal(mctx.Config, &config)
	if err != nil {
		// TODO(elliot): More context.
		return nil, err
	}

	return &Module{
		config: config,
	}, nil
}

// Render attempts starts a background process to update the time displayed in a label that is then
// returned to be placed on a bar.
func (m *Module) Render() (widgets.QLayout_ITF, error) {
	m.layout = widgets.NewQHBoxLayout()

	m.label = widgets.NewQLabel2(time.Now().Format(m.config.Format), nil, core.Qt__Widget)
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
				m.label.SetText(time.Now().Format(m.config.Format))
			}
		}
	}()

	m.layout.AddWidget(m.label, 0, core.Qt__AlignJustify)

	return m.layout, nil
}

// Destroy stops background processes and frees up resources.
func (m *Module) Destroy() error {
	if m.cfn != nil {
		m.cfn()
	}

	if m.layout != nil {
		m.layout.DestroyQHBoxLayout()
	}

	if m.label != nil {
		m.label.Destroy(true, true)
	}

	m.ctx = nil
	m.cfn = nil
	m.layout = nil
	m.label = nil

	return nil
}
