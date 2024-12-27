package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbLed106Controls struct {
	Overcurrent        *controls.SwitchControl
	RgbStrip           *controls.SwitchControl
	RgbPalette         *controls.RgbControl
	RgbStripHue        *controls.RangeControl
	RgbStripSaturation *controls.RangeControl
	RgbStripBrightness *controls.RangeControl
	HueChanging        *controls.SwitchControl
	HueChangingRate    *controls.ValueControl
	Channel4           *controls.SwitchControl
	Channel4Brightness *controls.RangeControl
	BoardTemperature   *controls.ValueControl
	AllowedPower       *controls.ValueControl
	Input1             *controls.SwitchControl
	Input2             *controls.SwitchControl
	Input2Counter      *controls.ValueControl
	Input1Counter      *controls.ValueControl
	Input3             *controls.SwitchControl
	Input3Counter      *controls.ValueControl
	Input4             *controls.SwitchControl
	Input4Counter      *controls.ValueControl
	Serial             *controls.TextControl
}

type WbLed106 struct {
	name     string
	device   string
	address  string
	Controls *WbLed106Controls
}

func (w *WbLed106) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbLed106) GetControlsInfo() []controls.ControlInfo {
	var infoList []controls.ControlInfo

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
				info := method.Call(nil)[0].Interface().(controls.ControlInfo)
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
			Overcurrent: controls.NewSwitchControl(client, name, "Overcurrent", controls.Meta{
				Type: "switch",

				Order:    12,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Перегрузка по току`},
			}),
			RgbStrip: controls.NewSwitchControl(client, name, "RGB Strip", controls.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Лента RGB`},
			}),
			RgbPalette: controls.NewRgbControl(client, name, "RGB Palette", controls.Meta{
				Type: "rgb",

				Order:    2,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Палитра RGB`},
			}),
			RgbStripHue: controls.NewRangeControl(client, name, "RGB Strip Hue", controls.Meta{
				Type: "range",

				Max: 360,

				Order:    3,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Оттенок RGB ленты (H)`},
			}),
			RgbStripSaturation: controls.NewRangeControl(client, name, "RGB Strip Saturation", controls.Meta{
				Type: "range",

				Max: 100,

				Order:    4,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Насыщенность цвета RGB ленты (S)`},
			}),
			RgbStripBrightness: controls.NewRangeControl(client, name, "RGB Strip Brightness", controls.Meta{
				Type: "range",

				Max: 100,

				Order:    5,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Яркость RGB ленты (V)`},
			}),
			HueChanging: controls.NewSwitchControl(client, name, "Hue Changing", controls.Meta{
				Type: "switch",

				Order:    6,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Изменение оттенка RGB ленты`},
			}),
			HueChangingRate: controls.NewValueControl(client, name, "Hue Changing Rate", controls.Meta{
				Type: "value",

				Max: 10000,
				Min: 3,

				Order:    7,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Время изменения оттенка (Hue) на 1° (мс)`},
			}),
			Channel4: controls.NewSwitchControl(client, name, "Channel 4", controls.Meta{
				Type: "switch",

				Order:    8,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Канал 4`},
			}),
			Channel4Brightness: controls.NewRangeControl(client, name, "Channel 4 Brightness", controls.Meta{
				Type: "range",

				Max: 100,

				Order:    9,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Яркость канала 4`},
			}),
			BoardTemperature: controls.NewValueControl(client, name, "Board Temperature", controls.Meta{
				Type:  "value",
				Units: "deg C",

				Order:    10,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Температура платы`},
			}),
			AllowedPower: controls.NewValueControl(client, name, "Allowed Power", controls.Meta{
				Type:  "value",
				Units: "%",

				Order:    11,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Разрешенная мощность`},
			}),
			Input1: controls.NewSwitchControl(client, name, "Input 1", controls.Meta{
				Type: "switch",

				Order:    13,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 1`},
			}),
			Input2: controls.NewSwitchControl(client, name, "Input 2", controls.Meta{
				Type: "switch",

				Order:    15,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 2`},
			}),
			Input2Counter: controls.NewValueControl(client, name, "Input 2 Counter", controls.Meta{
				Type: "value",

				Order:    16,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 2`},
			}),
			Input1Counter: controls.NewValueControl(client, name, "Input 1 Counter", controls.Meta{
				Type: "value",

				Order:    14,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 1`},
			}),
			Input3: controls.NewSwitchControl(client, name, "Input 3", controls.Meta{
				Type: "switch",

				Order:    17,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 3`},
			}),
			Input3Counter: controls.NewValueControl(client, name, "Input 3 Counter", controls.Meta{
				Type: "value",

				Order:    18,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 3`},
			}),
			Input4: controls.NewSwitchControl(client, name, "Input 4", controls.Meta{
				Type: "switch",

				Order:    19,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 4`},
			}),
			Input4Counter: controls.NewValueControl(client, name, "Input 4 Counter", controls.Meta{
				Type: "value",

				Order:    20,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 4`},
			}),
			Serial: controls.NewTextControl(client, name, "Serial", controls.Meta{
				Type: "text",

				Order:    21,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Серийный номер`},
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
