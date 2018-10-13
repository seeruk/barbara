package battery

import (
	"context"
	"encoding/json"

	"github.com/seeruk/barbara/barbara"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// Module ...
type Module struct {
	ctx context.Context
	cfn context.CancelFunc

	config Config
	layout *widgets.QHBoxLayout
	label  *widgets.QLabel
	icon   *gui.QIcon
}

// NewModule returns a new battery Module instance.
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

// Render ...
func (m *Module) Render() (widgets.QLayout_ITF, error) {
	m.layout = widgets.NewQHBoxLayout()
	m.label = widgets.NewQLabel(m.layout.Widget(), core.Qt__Widget)
	m.label.SetText("BAT0")

	m.layout.AddWidget(m.label, 0, core.Qt__AlignJustify)

	return m.layout, nil
}

// Destroy ...
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

	if m.icon != nil {
		m.icon.DestroyQIcon()
	}

	m.ctx = nil
	m.cfn = nil
	m.layout = nil
	m.label = nil
	m.icon = nil

	return nil
}
