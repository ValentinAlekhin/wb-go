package device

import (
	"context"
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbLed106Controls struct {
	RgbStrip           *control.SwitchControl
	RgbPalette         *control.RgbControl
	RgbStripHue        *control.RangeControl
	RgbStripSaturation *control.RangeControl
	RgbStripBrightness *control.RangeControl
	HueChanging        *control.SwitchControl
	HueChangingRate    *control.ValueControl
	Channel4           *control.SwitchControl
	Channel4Brightness *control.RangeControl
	BoardTemperature   *control.ValueControl
	AllowedPower       *control.ValueControl
	Overcurrent        *control.SwitchControl
	Input1             *control.SwitchControl
	Input1Counter      *control.ValueControl
	Input2             *control.SwitchControl
	Input2Counter      *control.ValueControl
	Input3             *control.SwitchControl
	Input3Counter      *control.ValueControl
	Input4             *control.SwitchControl
	Input4Counter      *control.ValueControl
	Serial             *control.TextControl
}

type WbLed106 struct {
	name     string
	Controls *WbLed106Controls
}

func (w *WbLed106) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbLed106     sync.Once
	instanceWbLed106 *WbLed106
)

func NewWbLed106(ctx context.Context, client mqtt.ClientInterface) *WbLed106 {
	onceWbLed106.Do(func() {
		name := "wb-led_106"

		controlList := &WbLed106Controls{
			RgbStrip: control.NewSwitchControl(ctx, client, name, "RGB Strip", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Лента RGB`},
			}),
			RgbPalette: control.NewRgbControl(ctx, client, name, "RGB Palette", control.Meta{
				Type: "rgb",

				Order:    2,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Палитра RGB`},
			}),
			RgbStripHue: control.NewRangeControl(ctx, client, name, "RGB Strip Hue", control.Meta{
				Type: "range",

				Max: 360,

				Order:    3,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Оттенок RGB ленты (H)`},
			}),
			RgbStripSaturation: control.NewRangeControl(ctx, client, name, "RGB Strip Saturation", control.Meta{
				Type: "range",

				Max: 100,

				Order:    4,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Насыщенность цвета RGB ленты (S)`},
			}),
			RgbStripBrightness: control.NewRangeControl(ctx, client, name, "RGB Strip Brightness", control.Meta{
				Type: "range",

				Max: 100,

				Order:    5,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Яркость RGB ленты (V)`},
			}),
			HueChanging: control.NewSwitchControl(ctx, client, name, "Hue Changing", control.Meta{
				Type: "switch",

				Order:    6,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Изменение оттенка RGB ленты`},
			}),
			HueChangingRate: control.NewValueControl(ctx, client, name, "Hue Changing Rate", control.Meta{
				Type: "value",

				Max: 10000,
				Min: 3,

				Order:    7,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Время изменения оттенка (Hue) на 1° (мс)`},
			}),
			Channel4: control.NewSwitchControl(ctx, client, name, "Channel 4", control.Meta{
				Type: "switch",

				Order:    8,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Канал 4`},
			}),
			Channel4Brightness: control.NewRangeControl(ctx, client, name, "Channel 4 Brightness", control.Meta{
				Type: "range",

				Max: 100,

				Order:    9,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Яркость канала 4`},
			}),
			BoardTemperature: control.NewValueControl(ctx, client, name, "Board Temperature", control.Meta{
				Type:  "value",
				Units: "deg C",

				Order:    10,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Температура платы`},
			}),
			AllowedPower: control.NewValueControl(ctx, client, name, "Allowed Power", control.Meta{
				Type:  "value",
				Units: "%",

				Order:    11,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Разрешенная мощность`},
			}),
			Overcurrent: control.NewSwitchControl(ctx, client, name, "Overcurrent", control.Meta{
				Type: "switch",

				Order:    12,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Перегрузка по току`},
			}),
			Input1: control.NewSwitchControl(ctx, client, name, "Input 1", control.Meta{
				Type: "switch",

				Order:    13,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 1`},
			}),
			Input1Counter: control.NewValueControl(ctx, client, name, "Input 1 Counter", control.Meta{
				Type: "value",

				Order:    14,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 1`},
			}),
			Input2: control.NewSwitchControl(ctx, client, name, "Input 2", control.Meta{
				Type: "switch",

				Order:    15,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 2`},
			}),
			Input2Counter: control.NewValueControl(ctx, client, name, "Input 2 Counter", control.Meta{
				Type: "value",

				Order:    16,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 2`},
			}),
			Input3: control.NewSwitchControl(ctx, client, name, "Input 3", control.Meta{
				Type: "switch",

				Order:    17,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 3`},
			}),
			Input3Counter: control.NewValueControl(ctx, client, name, "Input 3 Counter", control.Meta{
				Type: "value",

				Order:    18,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 3`},
			}),
			Input4: control.NewSwitchControl(ctx, client, name, "Input 4", control.Meta{
				Type: "switch",

				Order:    19,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 4`},
			}),
			Input4Counter: control.NewValueControl(ctx, client, name, "Input 4 Counter", control.Meta{
				Type: "value",

				Order:    20,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 4`},
			}),
			Serial: control.NewTextControl(ctx, client, name, "Serial", control.Meta{
				Type: "text",

				Order:    21,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Серийный номер`},
			}),
		}

		instanceWbLed106 = &WbLed106{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbLed106
}
