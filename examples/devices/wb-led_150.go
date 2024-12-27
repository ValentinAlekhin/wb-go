package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbLed150Controls struct {
	Cct1             *controls.SwitchControl
	Cct1Temperature  *controls.RangeControl
	Cct1Brightness   *controls.RangeControl
	Cct2             *controls.SwitchControl
	Cct2Temperature  *controls.RangeControl
	Cct2Brightness   *controls.RangeControl
	BoardTemperature *controls.ValueControl
	AllowedPower     *controls.ValueControl
	Overcurrent      *controls.SwitchControl
	Input1           *controls.SwitchControl
	Input2           *controls.SwitchControl
	Input2Counter    *controls.ValueControl
	Input3           *controls.SwitchControl
	Input3Counter    *controls.ValueControl
	Input4           *controls.SwitchControl
	Input4Counter    *controls.ValueControl
	Serial           *controls.TextControl
}

type WbLed150 struct {
	name     string
	device   string
	address  string
	Controls *WbLed150Controls
}

func (w *WbLed150) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbLed150) GetControlsInfo() []controls.ControlInfo {
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
	onceWbLed150     sync.Once
	instanceWbLed150 *WbLed150
)

func NewWbLed150(client *mqtt.Client) *WbLed150 {
	onceWbLed150.Do(func() {
		device := "wb-led"
		address := "150"
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &WbLed150Controls{
			Cct1: controls.NewSwitchControl(client, name, "CCT1", controls.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Лента CCT1`},
			}),
			Cct1Temperature: controls.NewRangeControl(client, name, "CCT1 Temperature", controls.Meta{
				Type: "range",

				Max: 100,

				Order:    2,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Цветовая температура ленты CCT1`},
			}),
			Cct1Brightness: controls.NewRangeControl(client, name, "CCT1 Brightness", controls.Meta{
				Type: "range",

				Max: 100,

				Order:    3,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Яркость ленты CCT1`},
			}),
			Cct2: controls.NewSwitchControl(client, name, "CCT2", controls.Meta{
				Type: "switch",

				Order:    4,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Лента CCT2`},
			}),
			Cct2Temperature: controls.NewRangeControl(client, name, "CCT2 Temperature", controls.Meta{
				Type: "range",

				Max: 100,

				Order:    5,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Цветовая температура ленты CCT2`},
			}),
			Cct2Brightness: controls.NewRangeControl(client, name, "CCT2 Brightness", controls.Meta{
				Type: "range",

				Max: 100,

				Order:    6,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Яркость ленты CCT2`},
			}),
			BoardTemperature: controls.NewValueControl(client, name, "Board Temperature", controls.Meta{
				Type:  "value",
				Units: "deg C",

				Order:    7,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Температура платы`},
			}),
			AllowedPower: controls.NewValueControl(client, name, "Allowed Power", controls.Meta{
				Type:  "value",
				Units: "%",

				Order:    8,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Разрешенная мощность`},
			}),
			Overcurrent: controls.NewSwitchControl(client, name, "Overcurrent", controls.Meta{
				Type: "switch",

				Order:    9,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Перегрузка по току`},
			}),
			Input1: controls.NewSwitchControl(client, name, "Input 1", controls.Meta{
				Type: "switch",

				Order:    10,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 1`},
			}),
			Input2: controls.NewSwitchControl(client, name, "Input 2", controls.Meta{
				Type: "switch",

				Order:    11,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 2`},
			}),
			Input2Counter: controls.NewValueControl(client, name, "Input 2 Counter", controls.Meta{
				Type: "value",

				Order:    12,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 2`},
			}),
			Input3: controls.NewSwitchControl(client, name, "Input 3", controls.Meta{
				Type: "switch",

				Order:    13,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 3`},
			}),
			Input3Counter: controls.NewValueControl(client, name, "Input 3 Counter", controls.Meta{
				Type: "value",

				Order:    14,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 3`},
			}),
			Input4: controls.NewSwitchControl(client, name, "Input 4", controls.Meta{
				Type: "switch",

				Order:    15,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 4`},
			}),
			Input4Counter: controls.NewValueControl(client, name, "Input 4 Counter", controls.Meta{
				Type: "value",

				Order:    16,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 4`},
			}),
			Serial: controls.NewTextControl(client, name, "Serial", controls.Meta{
				Type: "text",

				Order:    17,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Серийный номер`},
			}),
		}

		instanceWbLed150 = &WbLed150{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceWbLed150
}
