package virtualdevice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/pkg/virualcontrol"
	"gorm.io/gorm"
	"time"
)

type Thermostat struct {
	client              wb.ClientInterface
	Controls            ThermostatControls
	meta                Meta
	metaTopic           string
	temperatureControls []*control.ValueControl
	ticker              *time.Ticker
	hysteresis          float64
	loaded              bool
}

type ThermostatControls struct {
	TargetTemperature  *virualcontrol.VirtualRangeControl
	CurrentTemperature *virualcontrol.VirtualValueControl
	Enabled            *virualcontrol.VirtualSwitchControl
	Relay              *virualcontrol.VirtualSwitchControl
}

type ThermostatConfig struct {
	DB                  *gorm.DB
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

	var sum float64 = 0
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
	t.ticker = time.NewTicker(1 * time.Second)
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

func (t *Thermostat) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
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

func NewThermostat(ctx context.Context, config ThermostatConfig) (*Thermostat, error) {
	if config.Client == nil {
		return &Thermostat{}, errors.New("client is nil")
	}

	if config.DB == nil {
		return &Thermostat{}, errors.New("db is nil")
	}

	if config.Device == "" {
		return &Thermostat{}, errors.New("device is empty")
	}

	err := migrate(config.DB)
	if err != nil {
		return &Thermostat{}, err
	}

	deviceFullName := getDeviceFullName(config.Device)

	t := &Thermostat{
		client:              config.Client,
		Controls:            ThermostatControls{},
		metaTopic:           fmt.Sprintf(conventions.CONV_DEVICE_META_V2_FMT, deviceFullName),
		temperatureControls: config.TemperatureControls,
		hysteresis:          config.Hysteresis,
		meta:                Meta{Name: config.Device, Driver: "wb-go"},
	}

	t.Controls.TargetTemperature = virualcontrol.NewVirtualRangeControl(ctx, virualcontrol.RangeOptions{
		BaseOptions: virualcontrol.BaseOptions{
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Set Point",
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
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Current Temperature",
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
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Enabled",
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
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Relay",
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
