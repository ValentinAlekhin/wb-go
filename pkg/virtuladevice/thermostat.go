package virtuladevice

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/pkg/virualcontrol"
	"time"
)

type Thermostat struct {
	client              wb.ClientInterface
	controls            ThermostatControls
	meta                Meta
	metaTopic           string
	temperatureControls []*control.ValueControl
	ticker              *time.Ticker
	hysteresis          float64
}

type ThermostatControls struct {
	targetTemperature  *virualcontrol.VirtualRangeControl
	currentTemperature *virualcontrol.VirtualValueControl
	enabled            *virualcontrol.VirtualSwitchControl
	heater             *virualcontrol.VirtualSwitchControl
}

type ThermostatConfig struct {
	Client              wb.ClientInterface
	Device              string
	TargetTemperature   int
	TemperatureControls []*control.ValueControl
	Hysteresis          float64
}

func (t *Thermostat) TurnOn() {
	t.controls.enabled.TurnOn()
}

func (t *Thermostat) TurnOff() {
	t.controls.enabled.TurnOff()
}

func (t *Thermostat) SetTarget(value int) {
	t.controls.targetTemperature.SetValue(value)
}

func (t *Thermostat) AddHeaterWatcher(f func(payload control.SwitchControlWatcherPayload)) {
	t.controls.heater.AddWatcher(f)
}

func (t *Thermostat) updateCurrentTemperature() {
	var sum float64 = 0

	for _, temperatureControl := range t.temperatureControls {
		sum += temperatureControl.GetValue()
	}

	currentTemperature := sum / float64(len(t.temperatureControls))
	t.controls.currentTemperature.SetValue(currentTemperature)

	t.updateHeater()
}

func (t *Thermostat) updateHeater() {
	if !t.controls.enabled.GetValue() {
		t.controls.heater.TurnOff()
		return
	}

	currentTemp := t.controls.currentTemperature.GetValue()
	targetTemp := float64(t.controls.targetTemperature.GetValue())

	heightTemp := targetTemp + t.hysteresis
	lowTemp := targetTemp - t.hysteresis

	if currentTemp > heightTemp {
		t.controls.heater.TurnOff()
	} else if currentTemp < lowTemp {
		t.controls.heater.TurnOn()
	}
}

func (t *Thermostat) runTicker() {
	t.ticker = time.NewTicker(1 * time.Second)
	go func() {
		for range t.ticker.C {
			t.updateCurrentTemperature()
		}
	}()
}

func (t *Thermostat) setMeta() {
	byteMeta, err := json.Marshal(t.meta)
	if err != nil {
		fmt.Println(err)
	}

	_ = t.client.Publish(wb.PublishPayload{
		Topic:    t.metaTopic,
		Value:    string(byteMeta),
		QOS:      1,
		Retained: true,
	})
}

func NewThermostat(config ThermostatConfig) (*Thermostat, error) {
	if config.Client == nil {
		return &Thermostat{}, errors.New("client is nil")
	}

	if config.Device == "" {
		return &Thermostat{}, errors.New("device is empty")
	}

	deviceFullName := getDeviceFullName(config.Device)

	targetTemperature := virualcontrol.NewVirtualRangeControl(config.Client, deviceFullName, "Target Temperature", control.Meta{
		Units:    "°C",
		Order:    1,
		ReadOnly: false,
		Min:      0,
		Max:      100,
		Title:    control.MultilingualText{"ru": "Целевая температура"},
	}, func(p virualcontrol.OnRangeHandlerPayload) {
		p.Set(p.Value)
	})

	currentTemperature := virualcontrol.NewVirtualValueControl(config.Client, deviceFullName, "Current Temperature", control.Meta{
		Units:    "°C",
		Order:    2,
		ReadOnly: true,
		Title:    control.MultilingualText{"ru": "Текущая температура"},
	}, nil)

	enabled := virualcontrol.NewVirtualSwitchControl(config.Client, deviceFullName, "Enabled", control.Meta{
		ReadOnly: false,
		Order:    3,
		Title:    control.MultilingualText{"ru": "Термостат включен"},
	}, func(p virualcontrol.OnSwitchHandlerPayload) {
		p.Set(p.Value)
	})

	on := virualcontrol.NewVirtualSwitchControl(config.Client, deviceFullName, "On", control.Meta{
		ReadOnly: true,
		Order:    4,
		Title:    control.MultilingualText{"ru": "Нагрев"},
	}, nil)

	t := &Thermostat{
		client: config.Client,
		controls: ThermostatControls{
			targetTemperature:  targetTemperature,
			currentTemperature: currentTemperature,
			enabled:            enabled,
			heater:             on,
		},
		metaTopic:           fmt.Sprintf(conventions.CONV_DEVICE_META_V2_FMT, deviceFullName),
		temperatureControls: config.TemperatureControls,
		hysteresis:          config.Hysteresis,
		meta:                Meta{Name: config.Device, Driver: "wb-go"},
	}

	t.setMeta()
	t.runTicker()

	return t, nil
}
