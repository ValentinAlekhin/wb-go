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
	"gorm.io/gorm"
	"math"
	"time"
)

type AdaptiveLight struct {
	client    wb.ClientInterface
	name      string
	fullName  string
	meta      Meta
	metaTopic string
	Controls  AdaptiveLightControls
	ticker    *time.Ticker
	loaded    bool
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
	DB     *gorm.DB
	Client wb.ClientInterface
	Device string
}

func (a *AdaptiveLight) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         a.fullName,
		ControlsInfo: nil,
	}
}

func (a *AdaptiveLight) update() {
	if !a.loaded {
		return
	}

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

	_ = a.client.Publish(wb.PublishPayload{
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

	if config.DB == nil {
		return nil, errors.New("db is nil")
	}

	if config.Device == "" {
		return nil, errors.New("device is empty")
	}

	err := migrate(config.DB)
	if err != nil {
		return nil, err
	}

	deviceFullName := getDeviceFullName(config.Device)

	al := &AdaptiveLight{
		client:    config.Client,
		name:      config.Device,
		fullName:  deviceFullName,
		metaTopic: fmt.Sprintf(conventions.CONV_DEVICE_META_V2_FMT, deviceFullName),
		meta:      Meta{Name: config.Device, Driver: "wb-go"},
		Controls:  AdaptiveLightControls{},
	}

	al.Controls.Enabled = virualcontrol.NewVirtualSwitchControl(virualcontrol.SwitchOptions{
		BaseOptions: virualcontrol.BaseOptions{
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Enabled",
			Meta: control.Meta{
				Order: 1,
				Title: control.MultilingualText{"ru": "Включено"},
			},
		},
		OnHandler: func(p virualcontrol.OnSwitchHandlerPayload) {
			p.Set(p.Value)
			al.update()
		},
		DefaultValue: true,
	})

	al.Controls.MinTemp = virualcontrol.NewVirtualRangeControl(virualcontrol.RangeOptions{
		BaseOptions: virualcontrol.BaseOptions{
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Min Temperature",
			Meta: control.Meta{
				Order: 2,
				Max:   100,
				Min:   0,
				Title: control.MultilingualText{"ru": "Минимальная температура"},
			},
		},
		OnHandler: func(p virualcontrol.OnRangeHandlerPayload) {
			p.Set(p.Value)
			al.update()
		},
		DefaultValue: 0,
	})

	al.Controls.MaxTemp = virualcontrol.NewVirtualRangeControl(virualcontrol.RangeOptions{
		BaseOptions: virualcontrol.BaseOptions{
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Max Temperature",
			Meta: control.Meta{
				Order: 3,
				Max:   100,
				Min:   0,
				Title: control.MultilingualText{"ru": "Максимальная температура"},
			},
		},
		OnHandler: func(p virualcontrol.OnRangeHandlerPayload) {
			p.Set(p.Value)
			al.update()
		},
		DefaultValue: 100,
	})

	al.Controls.CurrentTemp = virualcontrol.NewVirtualRangeControl(virualcontrol.RangeOptions{
		BaseOptions: virualcontrol.BaseOptions{
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Temperature",
			Meta: control.Meta{
				Order: 4,
				Max:   100,
				Min:   0,
				Title: control.MultilingualText{"ru": "Температура"},
			},
		},
		DefaultValue: 100,
	})

	al.Controls.MinBrightness = virualcontrol.NewVirtualRangeControl(virualcontrol.RangeOptions{
		BaseOptions: virualcontrol.BaseOptions{
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Min Brightness",
			Meta: control.Meta{
				Order: 5,
				Max:   100,
				Min:   0,
				Title: control.MultilingualText{"ru": "Минимальная яркость"},
			},
		},
		OnHandler: func(p virualcontrol.OnRangeHandlerPayload) {
			p.Set(p.Value)
			al.update()
		},
		DefaultValue: 0,
	})

	al.Controls.MaxBrightness = virualcontrol.NewVirtualRangeControl(virualcontrol.RangeOptions{
		BaseOptions: virualcontrol.BaseOptions{
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Max Brightness",
			Meta: control.Meta{
				Order: 6,
				Max:   100,
				Min:   0,
				Title: control.MultilingualText{"ru": "Максимальная яркость"},
			},
		},
		OnHandler: func(p virualcontrol.OnRangeHandlerPayload) {
			p.Set(p.Value)
			al.update()
		},
		DefaultValue: 100,
	})

	al.Controls.CurrentBrightness = virualcontrol.NewVirtualRangeControl(virualcontrol.RangeOptions{
		BaseOptions: virualcontrol.BaseOptions{
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Brightness",
			Meta: control.Meta{
				Order: 7,
				Max:   100,
				Min:   0,
				Title: control.MultilingualText{"ru": "Яркость"},
			},
		},
		DefaultValue: 100,
	})

	al.Controls.SleepMode = virualcontrol.NewVirtualSwitchControl(virualcontrol.SwitchOptions{
		BaseOptions: virualcontrol.BaseOptions{
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Sleep Mode",
			Meta: control.Meta{
				ReadOnly: true,
				Order:    8,
				Title:    control.MultilingualText{"ru": "Режим сна"},
			},
		},
		DefaultValue: false,
	})

	al.Controls.Sunrise = virualcontrol.NewVirtualTimeControl(virualcontrol.TimeOptions{
		BaseOptions: virualcontrol.BaseOptions{
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Sunrise",
			Meta: control.Meta{
				Order: 9,
				Title: control.MultilingualText{"ru": "Рассвет"},
			},
		},
		DefaultValue: carbon.ParseByFormat("06:00:00", "H:i:s").StdTime(),
		OnHandler: func(p virualcontrol.OnTimeHandlerPayload) {
			p.Set(p.Value)
			al.update()
		},
	})

	al.Controls.Sunset = virualcontrol.NewVirtualTimeControl(virualcontrol.TimeOptions{
		BaseOptions: virualcontrol.BaseOptions{
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Sunset",
			Meta: control.Meta{
				Order: 10,
				Title: control.MultilingualText{"ru": "Закат"},
			},
		},
		DefaultValue: carbon.ParseByFormat("18:00:00", "H:i:s").StdTime(),
		OnHandler: func(p virualcontrol.OnTimeHandlerPayload) {
			p.Set(p.Value)
			al.update()
		},
	})

	al.Controls.SleepStart = virualcontrol.NewVirtualTimeControl(virualcontrol.TimeOptions{
		BaseOptions: virualcontrol.BaseOptions{
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Slip Start",
			Meta: control.Meta{
				Order: 11,
				Title: control.MultilingualText{"ru": "Начало сна"},
			},
		},
		DefaultValue: carbon.ParseByFormat("23:00:00", "H:i:s").StdTime(),
		OnHandler: func(p virualcontrol.OnTimeHandlerPayload) {
			p.Set(p.Value)
			al.update()
		},
	})

	al.Controls.SleepEnd = virualcontrol.NewVirtualTimeControl(virualcontrol.TimeOptions{
		BaseOptions: virualcontrol.BaseOptions{
			DB:     config.DB,
			Client: config.Client,
			Device: deviceFullName,
			Name:   "Slip End",
			Meta: control.Meta{
				Order: 11,
				Title: control.MultilingualText{"ru": "Конец сна"},
			},
		},
		DefaultValue: carbon.ParseByFormat("06:00:00", "H:i:s").StdTime(),
		OnHandler: func(p virualcontrol.OnTimeHandlerPayload) {
			p.Set(p.Value)
			al.update()
		},
	})

	al.setMeta()
	al.runTicker()

	al.loaded = true

	return al, nil
}
