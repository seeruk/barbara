package battery

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/seeruk/barbara/barbara"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// powerSupplyPath ...
const powerSupplyPath = "/sys/class/power_supply"

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

	m.ctx, m.cfn = context.WithCancel(context.Background())
	m.onTick()

	go func() {
		// TODO(elliot): Config.
		ticker := time.NewTicker(10 * time.Second)

		for {
			select {
			case <-m.ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				m.onTick()
			}
		}
	}()

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

func (m *Module) onTick() {
	labelText, err := m.getLabelText()
	if err != nil {
		return
	}

	m.label.SetText(labelText)
}

// getLabelText ...
func (m *Module) getLabelText() (string, error) {
	percentage, err := m.getBatteryPercentage()
	if err != nil {
		return "", err
	}

	timeRemaining, err := m.getTimeRemaining(false)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%% (%s)", percentage, timeRemaining), nil
}

// getBatteryPercentage ...
func (m *Module) getBatteryPercentage() (string, error) {
	return fileToStr(fmt.Sprintf("%s/%s/capacity", powerSupplyPath, m.config.PowerSupply))
}

// getTimeRemaining ...
func (m *Module) getTimeRemaining(isCharging bool) (string, error) {
	ps := m.config.PowerSupply

	chargeNow, err := fileToFloat(fmt.Sprintf("%s/%s/charge_now", powerSupplyPath, ps))
	if err != nil {
		// TODO(elliot): Context.
		return "", err
	}

	currentNow, err := fileToFloat(fmt.Sprintf("%s/%s/current_now", powerSupplyPath, ps))
	if err != nil {
		// TODO(elliot): Context.
		return "", err
	}

	hours := chargeNow / currentNow

	if isCharging {
		chargeFull, err := fileToFloat(fmt.Sprintf("%s/%s/charge_full", powerSupplyPath, ps))
		if err != nil {
			// TODO(elliot): Context.
			return "", err
		}

		hours = (chargeFull - chargeNow) / currentNow
	}

	minutes := 0.6 * (hours * 100)
	duration := time.Duration(minutes) * time.Minute

	minutes = duration.Minutes() - (math.Floor(duration.Hours()) * 60)

	return fmt.Sprintf("%02.0f:%02.0f", duration.Hours(), minutes), nil
}

// fileToStr ...
func fileToStr(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		// TODO(elliot): Context.
		return "", err
	}

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		// TODO(elliot): Context.
		return "", err
	}

	return strings.TrimSpace(string(bs)), nil
}

// fileToFloat ...
func fileToFloat(fileName string) (float64, error) {
	str, err := fileToStr(fileName)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(str, 64)
}
