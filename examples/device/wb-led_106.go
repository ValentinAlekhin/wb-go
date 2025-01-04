package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbLed106Controls struct {
	Overcurrent        *control.SwitchControl
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
	Input1             *control.SwitchControl
	Input2             *control.SwitchControl
	Input2Counter      *control.ValueControl
	Input1Counter      *control.ValueControl
	Input3             *control.SwitchControl
	Input3Counter      *control.ValueControl
	Input4             *control.SwitchControl
	Input4Counter      *control.ValueControl
	Serial             *control.TextControl
}

type WbLed106 struct {
	name     string
	device   string
	address  string
	Controls *WbLed106Controls
}

func (w *WbLed106) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbLed106) GetControlsInfo() []control.ControlInfo {
	var infoList []control.ControlInfo

	// Получаем значение и тип структуры Controls
	controlsValue := reflect.ValueOf(w.Controls).Elem()
	controlsType := controlsValue.Type()

	// Проходимся по всем полям структуры Controls
	for i := 0; i < controlsValue.NumField(); i++ {
		field := controlsValue.Field(i)

		// Проверяем, что поле является указателем и не nil
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			// Проверяем, реализует ли поле метод GetInfo
			method := field.MethodByName("GetInfo")
			if method.IsValid() {
				// Вызываем метод GetInfo
				info := method.Call(nil)[0].Interface().(control.ControlInfo)
				infoList = append(infoList, info)
			} else {
				fmt.Printf("Field %s does not implement GetInfo\n", controlsType.Field(i).Name)
			}
		}
	}

	return infoList
}

var (
	onceWbLed106     sync.Once
	instanceWbLed106 *WbLed106
)

func NewWbLed106(client *mqtt.Client) *WbLed106 {
	onceWbLed106.Do(func() {
		device := "wb-led"
		address := "106"
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &WbLed106Controls{
			Overcurrent: control.NewSwitchControl(client, name, "Overcurrent", control.Meta{
				Type: "switch",

				Order:    12,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Перегрузка по току`},
			}),
			RgbStrip: control.NewSwitchControl(client, name, "RGB Strip", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Лента RGB`},
			}),
			RgbPalette: control.NewRgbControl(client, name, "RGB Palette", control.Meta{
				Type: "rgb",

				Order:    2,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Палитра RGB`},
			}),
			RgbStripHue: control.NewRangeControl(client, name, "RGB Strip Hue", control.Meta{
				Type: "range",

				Max: 360,

				Order:    3,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Оттенок RGB ленты (H)`},
			}),
			RgbStripSaturation: control.NewRangeControl(client, name, "RGB Strip Saturation", control.Meta{
				Type: "range",

				Max: 100,

				Order:    4,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Насыщенность цвета RGB ленты (S)`},
			}),
			RgbStripBrightness: control.NewRangeControl(client, name, "RGB Strip Brightness", control.Meta{
				Type: "range",

				Max: 100,

				Order:    5,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Яркость RGB ленты (V)`},
			}),
			HueChanging: control.NewSwitchControl(client, name, "Hue Changing", control.Meta{
				Type: "switch",

				Order:    6,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Изменение оттенка RGB ленты`},
			}),
			HueChangingRate: control.NewValueControl(client, name, "Hue Changing Rate", control.Meta{
				Type: "value",

				Max: 10000,
				Min: 3,

				Order:    7,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Время изменения оттенка (Hue) на 1° (мс)`},
			}),
			Channel4: control.NewSwitchControl(client, name, "Channel 4", control.Meta{
				Type: "switch",

				Order:    8,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Канал 4`},
			}),
			Channel4Brightness: control.NewRangeControl(client, name, "Channel 4 Brightness", control.Meta{
				Type: "range",

				Max: 100,

				Order:    9,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Яркость канала 4`},
			}),
			BoardTemperature: control.NewValueControl(client, name, "Board Temperature", control.Meta{
				Type:  "value",
				Units: "deg C",

				Order:    10,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Температура платы`},
			}),
			AllowedPower: control.NewValueControl(client, name, "Allowed Power", control.Meta{
				Type:  "value",
				Units: "%",

				Order:    11,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Разрешенная мощность`},
			}),
			Input1: control.NewSwitchControl(client, name, "Input 1", control.Meta{
				Type: "switch",

				Order:    13,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 1`},
			}),
			Input2: control.NewSwitchControl(client, name, "Input 2", control.Meta{
				Type: "switch",

				Order:    15,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 2`},
			}),
			Input2Counter: control.NewValueControl(client, name, "Input 2 Counter", control.Meta{
				Type: "value",

				Order:    16,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 2`},
			}),
			Input1Counter: control.NewValueControl(client, name, "Input 1 Counter", control.Meta{
				Type: "value",

				Order:    14,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 1`},
			}),
			Input3: control.NewSwitchControl(client, name, "Input 3", control.Meta{
				Type: "switch",

				Order:    17,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 3`},
			}),
			Input3Counter: control.NewValueControl(client, name, "Input 3 Counter", control.Meta{
				Type: "value",

				Order:    18,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 3`},
			}),
			Input4: control.NewSwitchControl(client, name, "Input 4", control.Meta{
				Type: "switch",

				Order:    19,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 4`},
			}),
			Input4Counter: control.NewValueControl(client, name, "Input 4 Counter", control.Meta{
				Type: "value",

				Order:    20,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 4`},
			}),
			Serial: control.NewTextControl(client, name, "Serial", control.Meta{
				Type: "text",

				Order:    21,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Серийный номер`},
			}),
		}

		instanceWbLed106 = &WbLed106{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceWbLed106
}
