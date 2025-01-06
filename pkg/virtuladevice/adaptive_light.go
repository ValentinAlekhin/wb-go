package virtuladevice

import (
	"encoding/json"
	"fmt"
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
	meta      Meta
	metaTopic string
	controls  AdaptiveLightControls
	ticker    *time.Ticker
}

type AdaptiveLightControls struct {
	enabled *virualcontrol.VirtualSwitchControl

	minTemp     *virualcontrol.VirtualRangeControl
	maxTemp     *virualcontrol.VirtualRangeControl
	currentTemp *virualcontrol.VirtualRangeControl

	minBrightness     *virualcontrol.VirtualRangeControl
	maxBrightness     *virualcontrol.VirtualRangeControl
	currentBrightness *virualcontrol.VirtualRangeControl

	sleepMode *virualcontrol.VirtualSwitchControl

	Sunrise *virualcontrol.VirtualTimeControl
	Sunset  *virualcontrol.VirtualTimeControl

	SleepStart *virualcontrol.VirtualTimeControl
	SleepEnd   *virualcontrol.VirtualTimeControl
}

type AdaptiveLightConfig struct {
	Client *wb.Client
	Device string
}

func (a *AdaptiveLight) update() {
	if !a.controls.enabled.GetValue() {
		return
	}

	now := carbon.Now().SetDate(0, 1, 1)

	a.setSleepMode(now)
	a.setBrightness()
	a.setColorTemp(now)
}

func (a *AdaptiveLight) setSleepMode(now carbon.Carbon) {
	slipStart := carbon.CreateFromStdTime(a.controls.SleepStart.GetValue())
	slipEnd := carbon.CreateFromStdTime(a.controls.SleepEnd.GetValue())

	if now.Gte(slipStart) || now.Lte(slipEnd) {
		a.controls.sleepMode.SetValue(true)
	} else {
		a.controls.sleepMode.SetValue(false)
	}
}

func (a *AdaptiveLight) setBrightness() {
	maxBrightness := a.controls.maxBrightness.GetValue()
	minBrightness := a.controls.minBrightness.GetValue()

	if a.controls.sleepMode.GetValue() {
		a.controls.currentBrightness.SetValue(minBrightness)
	} else {
		a.controls.currentBrightness.SetValue(maxBrightness)
	}
}

func (a *AdaptiveLight) setColorTemp(now carbon.Carbon) {
	maxTemp := a.controls.maxTemp.GetValue()
	minTemp := a.controls.minTemp.GetValue()

	if a.controls.sleepMode.GetValue() {
		a.controls.currentTemp.SetValue(minTemp)
		return
	}

	sunrise := carbon.CreateFromStdTime(a.controls.Sunrise.GetValue())
	sunset := carbon.CreateFromStdTime(a.controls.Sunset.GetValue())

	if now.Gte(sunset) || now.Lte(sunrise) {
		a.controls.currentTemp.SetValue(minTemp)
		return
	}

	sunriseMinutes := sunrise.Hour()*60 + sunrise.Minute()
	sunsetMinutes := sunset.Hour()*60 + sunset.Minute()
	currentMinutes := now.Hour()*60 + now.Minute()

	dayLength := sunsetMinutes - sunriseMinutes
	minutesSinceSunrise := currentMinutes - sunriseMinutes
	ratio := float64(minutesSinceSunrise) / float64(dayLength)

	temp := int(float64(maxTemp) + float64(minTemp-maxTemp)*math.Pow(2*ratio-1, 2))

	fmt.Println("exp: ", temp)

	a.controls.currentTemp.SetValue(temp)
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

func NewAdaptiveLight(config AdaptiveLightConfig) *AdaptiveLight {
	enabled := virualcontrol.NewVirtualSwitchControl(config.Client, config.Device, "Enabled", control.Meta{
		Order: 1,
		Title: control.MultilingualText{"ru": "Включено"},
	}, func(p virualcontrol.OnSwitchHandlerPayload) {
		p.Set(p.Value)
	})

	minTemp := virualcontrol.NewVirtualRangeControl(config.Client, config.Device, "Min Temperature", control.Meta{
		Order: 2,
		Max:   100,
		Min:   0,
		Title: control.MultilingualText{"ru": "Минимальная температура"},
	}, func(p virualcontrol.OnRangeHandlerPayload) {
		p.Set(p.Value)
	})

	maxTemp := virualcontrol.NewVirtualRangeControl(config.Client, config.Device, "Max Temperature", control.Meta{
		Order: 3,
		Max:   100,
		Min:   0,
		Title: control.MultilingualText{"ru": "Максимальная температура"},
	}, func(p virualcontrol.OnRangeHandlerPayload) {
		p.Set(p.Value)
	})

	currentTemp := virualcontrol.NewVirtualRangeControl(config.Client, config.Device, "Temperature", control.Meta{
		ReadOnly: true,
		Order:    4,
		Max:      100,
		Min:      0,
		Title:    control.MultilingualText{"ru": "Температура"},
	}, func(p virualcontrol.OnRangeHandlerPayload) {
	})

	minBrightness := virualcontrol.NewVirtualRangeControl(config.Client, config.Device, "Min Brightness", control.Meta{
		Order: 5,
		Max:   100,
		Min:   0,
		Title: control.MultilingualText{"ru": "Минимальная яркость"},
	}, func(p virualcontrol.OnRangeHandlerPayload) {
		p.Set(p.Value)
	})

	maxBrightness := virualcontrol.NewVirtualRangeControl(config.Client, config.Device, "Max Brightness", control.Meta{
		Order: 6,
		Max:   100,
		Min:   0,
		Title: control.MultilingualText{"ru": "Максимальная яркость"},
	}, func(p virualcontrol.OnRangeHandlerPayload) {
		p.Set(p.Value)
	})

	currentBrightness := virualcontrol.NewVirtualRangeControl(config.Client, config.Device, "Brightness", control.Meta{
		ReadOnly: true,
		Order:    7,
		Max:      100,
		Min:      0,
		Title:    control.MultilingualText{"ru": "Яркость"},
	}, func(p virualcontrol.OnRangeHandlerPayload) {
	})

	sleepMode := virualcontrol.NewVirtualSwitchControl(config.Client, config.Device, "Sleep Mode", control.Meta{
		ReadOnly: true,
		Order:    8,
		Title:    control.MultilingualText{"ru": "Режим сна"},
	}, func(p virualcontrol.OnSwitchHandlerPayload) {
		p.Set(p.Value)
	})

	sunrise := virualcontrol.NewVirtualTimeValueControl(config.Client, config.Device, "Sunrise", control.Meta{
		Order: 9,
		Title: control.MultilingualText{"ru": "Рассвет"},
	}, func(p virualcontrol.OnTimeHandlerPayload) {
		p.Set(p.Value)
	})

	sunset := virualcontrol.NewVirtualTimeValueControl(config.Client, config.Device, "Sunset", control.Meta{
		Order: 10,
		Title: control.MultilingualText{"ru": "Закат"},
	}, func(p virualcontrol.OnTimeHandlerPayload) {
		p.Set(p.Value)
	})

	slipStart := virualcontrol.NewVirtualTimeValueControl(config.Client, config.Device, "Slip Start", control.Meta{
		Order: 11,
		Title: control.MultilingualText{"ru": "Начало сна"},
	}, func(p virualcontrol.OnTimeHandlerPayload) {
		p.Set(p.Value)
	})

	slipEnd := virualcontrol.NewVirtualTimeValueControl(config.Client, config.Device, "Slip End", control.Meta{
		Order: 12,
		Title: control.MultilingualText{"ru": "Конец сна"},
	}, func(p virualcontrol.OnTimeHandlerPayload) {
		if p.Error != nil {
			return
		}
		p.Set(p.Value)
		fmt.Println(p.Value.String())
	})

	al := &AdaptiveLight{
		client:    config.Client,
		metaTopic: fmt.Sprintf(conventions.CONV_DEVICE_META_V2_FMT, config.Device),
		meta:      Meta{Name: config.Device, Driver: "wb-go"},
		controls: AdaptiveLightControls{
			enabled:           enabled,
			minTemp:           minTemp,
			maxTemp:           maxTemp,
			currentTemp:       currentTemp,
			minBrightness:     minBrightness,
			maxBrightness:     maxBrightness,
			currentBrightness: currentBrightness,
			sleepMode:         sleepMode,
			Sunrise:           sunrise,
			Sunset:            sunset,
			SleepEnd:          slipEnd,
			SleepStart:        slipStart,
		},
	}

	al.setMeta()
	al.runTicker()

	return al
}
