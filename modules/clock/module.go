package clock

import (
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type Module struct {
}

func NewModule() *Module {
	return &Module{}
}

// Render ...
func (m *Module) Render() (widgets.QWidget_ITF, error) {
	dateLabel := widgets.NewQLabel2(time.Now().Format("15:04:05\nMon, 02 Jan"), nil, core.Qt__Widget)
	dateLabel.SetAlignment(core.Qt__AlignCenter)

	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-ticker.C:
				dateLabel.SetText(time.Now().Format("15:04:05\nMon, 02 Jan"))
			}
		}
	}()

	return dateLabel, nil
}

// Destroy ...
func (m *Module) Destroy() error {
	return nil
}
