package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbLed150Controls struct {
	Cct1             *control.SwitchControl
	Cct1Temperature  *control.RangeControl
	Cct1Brightness   *control.RangeControl
	Cct2             *control.SwitchControl
	Cct2Temperature  *control.RangeControl
	Cct2Brightness   *control.RangeControl
	BoardTemperature *control.ValueControl
	AllowedPower     *control.ValueControl
	Overcurrent      *control.SwitchControl
	Input1           *control.SwitchControl
	Input2           *control.SwitchControl
	Input2Counter    *control.ValueControl
	Input3           *control.SwitchControl
	Input3Counter    *control.ValueControl
	Input4           *control.SwitchControl
	Input4Counter    *control.ValueControl
	Serial           *control.TextControl
}

type WbLed150 struct {
	name     string
	Controls *WbLed150Controls
}

func (w *WbLed150) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbLed150) GetControlsInfo() []control.ControlInfo {
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
	onceWbLed150     sync.Once
	instanceWbLed150 *WbLed150
)

func NewWbLed150(client *mqtt.Client) *WbLed150 {
	onceWbLed150.Do(func() {
		name := "wb-led_150"

		controlList := &WbLed150Controls{
			Cct1: control.NewSwitchControl(client, name, "CCT1", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Лента CCT1`},
			}),
			Cct1Temperature: control.NewRangeControl(client, name, "CCT1 Temperature", control.Meta{
				Type: "range",

				Max: 100,

				Order:    2,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Цветовая температура ленты CCT1`},
			}),
			Cct1Brightness: control.NewRangeControl(client, name, "CCT1 Brightness", control.Meta{
				Type: "range",

				Max: 100,

				Order:    3,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Яркость ленты CCT1`},
			}),
			Cct2: control.NewSwitchControl(client, name, "CCT2", control.Meta{
				Type: "switch",

				Order:    4,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Лента CCT2`},
			}),
			Cct2Temperature: control.NewRangeControl(client, name, "CCT2 Temperature", control.Meta{
				Type: "range",

				Max: 100,

				Order:    5,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Цветовая температура ленты CCT2`},
			}),
			Cct2Brightness: control.NewRangeControl(client, name, "CCT2 Brightness", control.Meta{
				Type: "range",

				Max: 100,

				Order:    6,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Яркость ленты CCT2`},
			}),
			BoardTemperature: control.NewValueControl(client, name, "Board Temperature", control.Meta{
				Type:  "value",
				Units: "deg C",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Температура платы`},
			}),
			AllowedPower: control.NewValueControl(client, name, "Allowed Power", control.Meta{
				Type:  "value",
				Units: "%",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Разрешенная мощность`},
			}),
			Overcurrent: control.NewSwitchControl(client, name, "Overcurrent", control.Meta{
				Type: "switch",

				Order:    9,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Перегрузка по току`},
			}),
			Input1: control.NewSwitchControl(client, name, "Input 1", control.Meta{
				Type: "switch",

				Order:    10,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 1`},
			}),
			Input2: control.NewSwitchControl(client, name, "Input 2", control.Meta{
				Type: "switch",

				Order:    11,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 2`},
			}),
			Input2Counter: control.NewValueControl(client, name, "Input 2 Counter", control.Meta{
				Type: "value",

				Order:    12,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 2`},
			}),
			Input3: control.NewSwitchControl(client, name, "Input 3", control.Meta{
				Type: "switch",

				Order:    13,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 3`},
			}),
			Input3Counter: control.NewValueControl(client, name, "Input 3 Counter", control.Meta{
				Type: "value",

				Order:    14,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 3`},
			}),
			Input4: control.NewSwitchControl(client, name, "Input 4", control.Meta{
				Type: "switch",

				Order:    15,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 4`},
			}),
			Input4Counter: control.NewValueControl(client, name, "Input 4 Counter", control.Meta{
				Type: "value",

				Order:    16,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 4`},
			}),
			Serial: control.NewTextControl(client, name, "Serial", control.Meta{
				Type: "text",

				Order:    17,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Серийный номер`},
			}),
		}

		instanceWbLed150 = &WbLed150{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbLed150
}
