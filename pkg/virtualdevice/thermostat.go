package virtualdevice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ValentinAlekhin/wb-go/internal/db"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/pkg/virualcontrol"
)

// DefaultUpdateInterval определяет стандартный интервал между обновлениями термостата
const DefaultUpdateInterval = 1 * time.Second

// Thermostat представляет виртуальное устройство термостата
type Thermostat struct {
	client              wb.ClientInterface
	Controls            ThermostatControls
	meta                Meta
	metaTopic           string
	temperatureControls []*control.ValueControl
	ticker              *time.Ticker
	updateInterval      time.Duration
	hysteresis          float64
	loaded              bool
}

// ThermostatControls содержит элементы управления термостатом
type ThermostatControls struct {
	TargetTemperature  *virualcontrol.VirtualRangeControl
	CurrentTemperature *virualcontrol.VirtualValueControl
	Enabled            *virualcontrol.VirtualSwitchControl
	Relay              *virualcontrol.VirtualSwitchControl
}

// ThermostatConfig содержит конфигурацию термостата
type ThermostatConfig struct {
	DB                  *sql.DB
	Client              wb.ClientInterface
	Device              string
	TargetTemperature   int
	TemperatureControls []*control.ValueControl
	Hysteresis          float64
}

func (t *Thermostat) update() {
	if !t.loaded {
		return
	}

	t.updateCurrentTemperature()
	t.updateRelay()
}

func (t *Thermostat) updateCurrentTemperature() {
	if len(t.temperatureControls) == 0 {
		t.Controls.CurrentTemperature.SetValue(0)
		return
	}

	// Оптимизация: предварительно выделяем память для суммы
	var sum float64
	for _, temperatureControl := range t.temperatureControls {
		sum += temperatureControl.GetValue()
	}

	currentTemperature := sum / float64(len(t.temperatureControls))
	t.Controls.CurrentTemperature.SetValue(currentTemperature)
}

func (t *Thermostat) updateRelay() {
	if !t.Controls.Enabled.GetValue() {
		t.Controls.Relay.TurnOff()
		return
	}

	currentTemp := t.Controls.CurrentTemperature.GetValue()
	targetTemp := float64(t.Controls.TargetTemperature.GetValue())

	heightTemp := targetTemp + t.hysteresis
	lowTemp := targetTemp - t.hysteresis

	if currentTemp > heightTemp {
		t.Controls.Relay.TurnOff()
	} else if currentTemp < lowTemp {
		t.Controls.Relay.TurnOn()
	}
}

func (t *Thermostat) runTicker(ctx context.Context) {
	t.ticker = time.NewTicker(t.updateInterval)
	go func() {
		defer t.ticker.Stop()
		for {
			select {
			case <-t.ticker.C:
				t.update()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Stop останавливает работу термостата
func (t *Thermostat) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
	}
}

// SetUpdateInterval устанавливает новый интервал обновления термостата
func (t *Thermostat) SetUpdateInterval(interval time.Duration) {
	if interval <= 0 {
		return
	}

	t.updateInterval = interval
	if t.ticker != nil {
		t.ticker.Reset(interval)
	}
}

func (t *Thermostat) setMeta() error {
	byteMeta, err := json.Marshal(t.meta)
	if err != nil {
		return fmt.Errorf("failed to marshal meta: %w", err)
	}

	err = t.client.Publish(wb.PublishPayload{
		Topic:    t.metaTopic,
		Value:    string(byteMeta),
		QOS:      1,
		Retained: true,
	})
	if err != nil {
		return fmt.Errorf("failed to publish meta: %w", err)
	}

	return nil
}

// NewThermostat создает новый экземпляр термостата
func NewThermostat(ctx context.Context, config ThermostatConfig) (*Thermostat, error) {
	if config.Client == nil {
		return nil, errors.New("client is nil")
	}

	if config.DB == nil {
		return nil, errors.New("db is nil")
	}

	if config.Device == "" {
		return nil, errors.New("device is empty")
	}

	if config.Hysteresis < 0 {
		return nil, errors.New("hysteresis must be non-negative")
	}

	err := db.MigrateOnce(config.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	q := db.NewQueries(config.DB)

	deviceFullName := getDeviceFullName(config.Device)

	t := &Thermostat{
		client:              config.Client,
		Controls:            ThermostatControls{},
		metaTopic:           fmt.Sprintf(conventions.CONV_DEVICE_META_V2_FMT, deviceFullName),
		temperatureControls: config.TemperatureControls,
		hysteresis:          config.Hysteresis,
		updateInterval:      DefaultUpdateInterval,
		meta:                Meta{Name: config.Device, Driver: "wb-go"},
	}

	t.Controls.TargetTemperature = virualcontrol.NewVirtualRangeControl(ctx, virualcontrol.RangeOptions{
		BaseOptions: virualcontrol.BaseOptions{
			Queries: q,
			Client:  config.Client,
			Device:  deviceFullName,
			Name:    "Set Point",
			Meta: control.Meta{
				Units:    "°C",
				Order:    1,
				ReadOnly: false,
				Min:      0,
				Max:      100,
				Title:    control.MultilingualText{"ru": "Целевая температура"},
			},
		},
		OnHandler: func(p virualcontrol.RangeHandlerPayload) {
			p.Set(p.Value)
			t.update()
		},
		DefaultValue: config.TargetTemperature,
	})

	t.Controls.CurrentTemperature = virualcontrol.NewVirtualValueControl(ctx, virualcontrol.ValueOptions{
		BaseOptions: virualcontrol.BaseOptions{
			Queries: q,
			Client:  config.Client,
			Device:  deviceFullName,
			Name:    "Current Temperature",
			Meta: control.Meta{
				Units:    "°C",
				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": "Текущая температура"},
			},
		},
	})

	t.Controls.Enabled = virualcontrol.NewVirtualSwitchControl(ctx, virualcontrol.SwitchOptions{
		BaseOptions: virualcontrol.BaseOptions{
			Queries: q,
			Client:  config.Client,
			Device:  deviceFullName,
			Name:    "Enabled",
			Meta: control.Meta{
				ReadOnly: false,
				Order:    3,
				Title:    control.MultilingualText{"ru": "Термостат включен"},
			},
		},
		DefaultValue: true,
		OnHandler: func(p virualcontrol.SwitchHandlerPayload) {
			p.Set(p.Value)
			t.update()
		},
	})

	t.Controls.Relay = virualcontrol.NewVirtualSwitchControl(ctx, virualcontrol.SwitchOptions{
		BaseOptions: virualcontrol.BaseOptions{
			Queries: q,
			Client:  config.Client,
			Device:  deviceFullName,
			Name:    "Relay",
			Meta: control.Meta{
				ReadOnly: true,
				Order:    4,
				Title:    control.MultilingualText{"ru": "Нагрев"},
			},
		},
	})

	if err := t.setMeta(); err != nil {
		return nil, err
	}

	t.runTicker(ctx)
	t.loaded = true

	return t, nil
}
