package virtuladevice

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/pkg/virualcontrol"
	"github.com/dromara/carbon/v2"
	"math"
	"time"
)

type AdaptiveLight struct {
	client    *wb.Client
	name      string
	fullName  string
	meta      Meta
	metaTopic string
	Controls  AdaptiveLightControls
	ticker    *time.Ticker
}

type AdaptiveLightControls struct {
	Enabled *virualcontrol.VirtualSwitchControl

	MinTemp     *virualcontrol.VirtualRangeControl
	MaxTemp     *virualcontrol.VirtualRangeControl
	CurrentTemp *virualcontrol.VirtualRangeControl

	MinBrightness     *virualcontrol.VirtualRangeControl
	MaxBrightness     *virualcontrol.VirtualRangeControl
	CurrentBrightness *virualcontrol.VirtualRangeControl

	SleepMode *virualcontrol.VirtualSwitchControl

	Sunrise *virualcontrol.VirtualTimeControl
	Sunset  *virualcontrol.VirtualTimeControl

	SleepStart *virualcontrol.VirtualTimeControl
	SleepEnd   *virualcontrol.VirtualTimeControl
}

type AdaptiveLightConfig struct {
	Client *wb.Client
	Device string
}

func (a *AdaptiveLight) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         a.fullName,
		ControlsInfo: nil,
	}
}

func (a *AdaptiveLight) update() {
	if !a.Controls.Enabled.GetValue() {
		return
	}

	now := carbon.Now().SetDate(0, 1, 1)

	a.setSleepMode(now)
	a.setBrightness()
	a.setColorTemp(now)
}

func (a *AdaptiveLight) setSleepMode(now carbon.Carbon) {
	slipStart := carbon.CreateFromStdTime(a.Controls.SleepStart.GetValue())
	slipEnd := carbon.CreateFromStdTime(a.Controls.SleepEnd.GetValue())

	if now.Gte(slipStart) || now.Lte(slipEnd) {
		a.Controls.SleepMode.SetValue(true)
	} else {
		a.Controls.SleepMode.SetValue(false)
	}
}

func (a *AdaptiveLight) setBrightness() {
	maxBrightness := a.Controls.MaxBrightness.GetValue()
	minBrightness := a.Controls.MinBrightness.GetValue()

	if a.Controls.SleepMode.GetValue() {
		a.Controls.CurrentBrightness.SetValue(minBrightness)
	} else {
		a.Controls.CurrentBrightness.SetValue(maxBrightness)
	}
}

func (a *AdaptiveLight) setColorTemp(now carbon.Carbon) {
	maxTemp := a.Controls.MaxTemp.GetValue()
	minTemp := a.Controls.MinTemp.GetValue()

	if a.Controls.SleepMode.GetValue() {
		a.Controls.CurrentTemp.SetValue(minTemp)
		return
	}

	sunrise := carbon.CreateFromStdTime(a.Controls.Sunrise.GetValue())
	sunset := carbon.CreateFromStdTime(a.Controls.Sunset.GetValue())

	if now.Gte(sunset) || now.Lte(sunrise) {
		a.Controls.CurrentTemp.SetValue(minTemp)
		return
	}

	sunriseMinutes := sunrise.Hour()*60 + sunrise.Minute()
	sunsetMinutes := sunset.Hour()*60 + sunset.Minute()
	currentMinutes := now.Hour()*60 + now.Minute()

	dayLength := sunsetMinutes - sunriseMinutes
	minutesSinceSunrise := currentMinutes - sunriseMinutes
	ratio := float64(minutesSinceSunrise) / float64(dayLength)

	temp := int(float64(maxTemp) + float64(minTemp-maxTemp)*math.Pow(2*ratio-1, 2))

	a.Controls.CurrentTemp.SetValue(temp)
}

func (a *AdaptiveLight) runTicker() {
	a.ticker = time.NewTicker(1 * time.Second)
	go func() {
		for range a.ticker.C {
			a.update()
		}
	}()
}

func (a *AdaptiveLight) setMeta() {
	byteMeta, err := json.Marshal(a.meta)
	if err != nil {
		fmt.Println(err)
	}

	a.client.Publish(wb.PublishPayload{
		Topic:    a.metaTopic,
		Value:    string(byteMeta),
		QOS:      1,
		Retained: true,
	})
}

func NewAdaptiveLight(config AdaptiveLightConfig) (*AdaptiveLight, error) {
	if config.Client == nil {
		return nil, errors.New("client is nil")
	}

	if config.Device == "" {
		return nil, errors.New("device is empty")
	}

	deviceFullName := getDeviceFullName(config.Device)

	enabled := virualcontrol.NewVirtualSwitchControl(config.Client, deviceFullName, "Enabled", control.Meta{
		Order: 1,
		Title: control.MultilingualText{"ru": "Включено"},
	}, func(p virualcontrol.OnSwitchHandlerPayload) {
		p.Set(p.Value)
	})

	minTemp := virualcontrol.NewVirtualRangeControl(config.Client, deviceFullName, "Min Temperature", control.Meta{
		Order: 2,
		Max:   100,
		Min:   0,
		Title: control.MultilingualText{"ru": "Минимальная температура"},
	}, func(p virualcontrol.OnRangeHandlerPayload) {
		p.Set(p.Value)
	})

	maxTemp := virualcontrol.NewVirtualRangeControl(config.Client, deviceFullName, "Max Temperature", control.Meta{
		Order: 3,
		Max:   100,
		Min:   0,
		Title: control.MultilingualText{"ru": "Максимальная температура"},
	}, func(p virualcontrol.OnRangeHandlerPayload) {
		p.Set(p.Value)
	})

	currentTemp := virualcontrol.NewVirtualRangeControl(config.Client, deviceFullName, "Temperature", control.Meta{
		ReadOnly: true,
		Order:    4,
		Max:      100,
		Min:      0,
		Title:    control.MultilingualText{"ru": "Температура"},
	}, nil)

	minBrightness := virualcontrol.NewVirtualRangeControl(config.Client, deviceFullName, "Min Brightness", control.Meta{
		Order: 5,
		Max:   100,
		Min:   0,
		Title: control.MultilingualText{"ru": "Минимальная яркость"},
	}, func(p virualcontrol.OnRangeHandlerPayload) {
		p.Set(p.Value)
	})

	maxBrightness := virualcontrol.NewVirtualRangeControl(config.Client, deviceFullName, "Max Brightness", control.Meta{
		Order: 6,
		Max:   100,
		Min:   0,
		Title: control.MultilingualText{"ru": "Максимальная яркость"},
	}, func(p virualcontrol.OnRangeHandlerPayload) {
		p.Set(p.Value)
	})

	currentBrightness := virualcontrol.NewVirtualRangeControl(config.Client, deviceFullName, "Brightness", control.Meta{
		ReadOnly: true,
		Order:    7,
		Max:      100,
		Min:      0,
		Title:    control.MultilingualText{"ru": "Яркость"},
	}, nil)

	sleepMode := virualcontrol.NewVirtualSwitchControl(config.Client, deviceFullName, "Sleep Mode", control.Meta{
		ReadOnly: true,
		Order:    8,
		Title:    control.MultilingualText{"ru": "Режим сна"},
	}, func(p virualcontrol.OnSwitchHandlerPayload) {
		p.Set(p.Value)
	})

	sunrise := virualcontrol.NewVirtualTimeValueControl(config.Client, deviceFullName, "Sunrise", control.Meta{
		Order: 9,
		Title: control.MultilingualText{"ru": "Рассвет"},
	}, func(p virualcontrol.OnTimeHandlerPayload) {
		p.Set(p.Value)
	})

	sunset := virualcontrol.NewVirtualTimeValueControl(config.Client, deviceFullName, "Sunset", control.Meta{
		Order: 10,
		Title: control.MultilingualText{"ru": "Закат"},
	}, func(p virualcontrol.OnTimeHandlerPayload) {
		p.Set(p.Value)
	})

	slipStart := virualcontrol.NewVirtualTimeValueControl(config.Client, deviceFullName, "Slip Start", control.Meta{
		Order: 11,
		Title: control.MultilingualText{"ru": "Начало сна"},
	}, func(p virualcontrol.OnTimeHandlerPayload) {
		p.Set(p.Value)
	})

	slipEnd := virualcontrol.NewVirtualTimeValueControl(config.Client, deviceFullName, "Slip End", control.Meta{
		Order: 12,
		Title: control.MultilingualText{"ru": "Конец сна"},
	}, func(p virualcontrol.OnTimeHandlerPayload) {
		if p.Error != nil {
			return
		}
		p.Set(p.Value)
	})

	al := &AdaptiveLight{
		client:    config.Client,
		name:      config.Device,
		fullName:  deviceFullName,
		metaTopic: fmt.Sprintf(conventions.CONV_DEVICE_META_V2_FMT, deviceFullName),
		meta:      Meta{Name: config.Device, Driver: "wb-go"},
		Controls: AdaptiveLightControls{
			Enabled:           enabled,
			MinTemp:           minTemp,
			MaxTemp:           maxTemp,
			CurrentTemp:       currentTemp,
			MinBrightness:     minBrightness,
			MaxBrightness:     maxBrightness,
			CurrentBrightness: currentBrightness,
			SleepMode:         sleepMode,
			Sunrise:           sunrise,
			Sunset:            sunset,
			SleepEnd:          slipEnd,
			SleepStart:        slipStart,
		},
	}

	al.setMeta()
	al.runTicker()

	return al, nil
}
