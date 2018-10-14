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

	config    Config
	layout    *widgets.QHBoxLayout
	iconLabel *widgets.QLabel
	label     *widgets.QLabel
}

// NewModule returns a new battery Module instance.
// TODO(elliot): This is a mess. Let's make a service for battery information, then just use this as
// the UI component for it, moving all of the actual logic out of here...
// TODO(elliot): If we're at 100%, probably no need to show a time remaining...
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
	m.iconLabel = widgets.NewQLabel(nil, core.Qt__Widget)
	m.label = widgets.NewQLabel(nil, core.Qt__Widget)

	m.ctx, m.cfn = context.WithCancel(context.Background())
	m.onTick()

	statusCh := make(chan struct{}, 1)

	go func() {
		// TODO(elliot): Config.
		ticker := time.NewTicker(time.Second)
		oldStatus := m.getBatteryStatus()

		for {
			select {
			case <-m.ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				status := m.getBatteryStatus()
				if status != oldStatus {
					statusCh <- struct{}{}
				}

				oldStatus = status
			}
		}
	}()

	go func() {
		// TODO(elliot): Config.
		ticker := time.NewTicker(10 * time.Second)

		for {
			select {
			case <-m.ctx.Done():
				ticker.Stop()
				return
			case <-statusCh:
				fmt.Println("Updating display")
				m.onTick()
			case <-ticker.C:
				m.onTick()
			}
		}
	}()

	m.layout.AddWidget(m.iconLabel, 0, core.Qt__AlignJustify)
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

	if m.iconLabel != nil {
		m.iconLabel.Destroy(true, true)
	}

	m.ctx = nil
	m.cfn = nil
	m.layout = nil
	m.label = nil
	m.iconLabel = nil

	return nil
}

// onTick ...
func (m *Module) onTick() {
	status := m.getBatteryStatus()

	percentage := m.getBatteryPercentage()
	labelText := m.getLabelText(percentage, status)

	var iconLevel string
	switch {
	case percentage == 100:
		iconLevel = "full"
	case percentage >= 66:
		iconLevel = "good"
	case percentage >= 33:
		iconLevel = "medium"
	case percentage > 0:
		iconLevel = "low"
	case percentage == 0:
		iconLevel = "empty"
	}

	var iconStatus string
	if status == "charging" {
		iconStatus = "-charging"
	}

	icon := gui.NewQIcon5(fmt.Sprintf(
		"/usr/share/icons/Paper-Mono-Dark/24x24/panel/battery-%s%s.svg",
		iconLevel,
		iconStatus,
	))

	m.iconLabel.SetPixmap(icon.Pixmap2(24, 24, gui.QIcon__Normal, gui.QIcon__On))
	m.label.SetText(labelText)
}

// getLabelText ...
func (m *Module) getLabelText(percentage float64, status string) string {
	timeRemaining := m.getTimeRemaining(status == "charging")

	return fmt.Sprintf("%.0f%% (%s)", percentage, timeRemaining)
}

// getBatteryPercentage ...
func (m *Module) getBatteryPercentage() float64 {
	return fileToFloat(fmt.Sprintf("%s/%s/capacity", powerSupplyPath, m.config.PowerSupply))
}

// getBatteryStatus ...
func (m *Module) getBatteryStatus() string {
	status, err := fileToStr(fmt.Sprintf("%s/%s/status", powerSupplyPath, m.config.PowerSupply))
	if err != nil {
		return "unknown"
	}

	return strings.ToLower(status)
}

// getTimeRemaining ...
func (m *Module) getTimeRemaining(isCharging bool) string {
	ps := m.config.PowerSupply

	chargeNow := fileToFloat(fmt.Sprintf("%s/%s/charge_now", powerSupplyPath, ps))
	currentNow := fileToFloat(fmt.Sprintf("%s/%s/current_now", powerSupplyPath, ps))

	hours := chargeNow / currentNow

	if isCharging {
		chargeFull := fileToFloat(fmt.Sprintf("%s/%s/charge_full", powerSupplyPath, ps))

		hours = (chargeFull - chargeNow) / currentNow
	}

	minutes := 0.6 * (hours * 100)
	duration := time.Duration(minutes) * time.Minute

	minutes = duration.Minutes() - (math.Floor(duration.Hours()) * 60)

	return fmt.Sprintf("%02.0f:%02.0f", duration.Hours(), minutes)
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
func fileToFloat(fileName string) float64 {
	// TODO(elliot): Log in here probably?

	str, err := fileToStr(fileName)
	if err != nil {
		return 0
	}

	flt, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}

	return flt
}
