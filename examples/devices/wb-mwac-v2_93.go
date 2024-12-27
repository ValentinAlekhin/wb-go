package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbMwacV293Controls struct {
	P1Volume       *controls.ValueControl
	P2Volume       *controls.ValueControl
	InputF1        *controls.SwitchControl
	InputF1Counter *controls.ValueControl
	InputF2        *controls.SwitchControl
	InputF2Counter *controls.ValueControl
	InputF3        *controls.SwitchControl
	InputF3Counter *controls.ValueControl
	InputF4        *controls.SwitchControl
	InputF4Counter *controls.ValueControl
	InputF5        *controls.SwitchControl
	InputF5Counter *controls.ValueControl
	InputS6        *controls.SwitchControl
	InputS6Counter *controls.ValueControl
	OutputK1       *controls.SwitchControl
	OutputK2       *controls.SwitchControl
	LeakageMode    *controls.SwitchControl
	CleaningMode   *controls.SwitchControl
	Serial         *controls.TextControl
}

type WbMwacV293 struct {
	name     string
	device   string
	address  string
	Controls *WbMwacV293Controls
}

func (w *WbMwacV293) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbMwacV293) GetControlsInfo() []controls.ControlInfo {
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
	onceWbMwacV293     sync.Once
	instanceWbMwacV293 *WbMwacV293
)

func NewWbMwacV293(client *mqtt.Client) *WbMwacV293 {
	onceWbMwacV293.Do(func() {
		device := "wb-mwac-v2"
		address := "93"
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &WbMwacV293Controls{
			P1Volume: controls.NewValueControl(client, name, "P1 Volume", controls.Meta{
				Type:  "value",
				Units: "m^3",

				Order:    1,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик объема P1`},
			}),
			P2Volume: controls.NewValueControl(client, name, "P2 Volume", controls.Meta{
				Type:  "value",
				Units: "m^3",

				Order:    2,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик объема P2`},
			}),
			InputF1: controls.NewSwitchControl(client, name, "Input F1", controls.Meta{
				Type: "switch",

				Order:    3,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход F1`},
			}),
			InputF1Counter: controls.NewValueControl(client, name, "Input F1 Counter", controls.Meta{
				Type: "value",

				Order:    4,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик входа F1`},
			}),
			InputF2: controls.NewSwitchControl(client, name, "Input F2", controls.Meta{
				Type: "switch",

				Order:    5,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход F2`},
			}),
			InputF2Counter: controls.NewValueControl(client, name, "Input F2 Counter", controls.Meta{
				Type: "value",

				Order:    6,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик входа F2`},
			}),
			InputF3: controls.NewSwitchControl(client, name, "Input F3", controls.Meta{
				Type: "switch",

				Order:    7,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход F3`},
			}),
			InputF3Counter: controls.NewValueControl(client, name, "Input F3 Counter", controls.Meta{
				Type: "value",

				Order:    8,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик входа F3`},
			}),
			InputF4: controls.NewSwitchControl(client, name, "Input F4", controls.Meta{
				Type: "switch",

				Order:    9,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход F4`},
			}),
			InputF4Counter: controls.NewValueControl(client, name, "Input F4 Counter", controls.Meta{
				Type: "value",

				Order:    10,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик входа F4`},
			}),
			InputF5: controls.NewSwitchControl(client, name, "Input F5", controls.Meta{
				Type: "switch",

				Order:    11,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход F5`},
			}),
			InputF5Counter: controls.NewValueControl(client, name, "Input F5 Counter", controls.Meta{
				Type: "value",

				Order:    12,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик входа F5`},
			}),
			InputS6: controls.NewSwitchControl(client, name, "Input S6", controls.Meta{
				Type: "switch",

				Order:    13,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход S6`},
			}),
			InputS6Counter: controls.NewValueControl(client, name, "Input S6 Counter", controls.Meta{
				Type: "value",

				Order:    14,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик входа S6`},
			}),
			OutputK1: controls.NewSwitchControl(client, name, "Output K1", controls.Meta{
				Type: "switch",

				Order:    15,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Выход K1`},
			}),
			OutputK2: controls.NewSwitchControl(client, name, "Output K2", controls.Meta{
				Type: "switch",

				Order:    16,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Выход K2`},
			}),
			LeakageMode: controls.NewSwitchControl(client, name, "Leakage Mode", controls.Meta{
				Type: "switch",

				Order:    17,
				ReadOnly: false,
				Title:    controls.MultilingualText{"en": `"Leakage" Mode`, "ru": `Режим "Протечка"`},
			}),
			CleaningMode: controls.NewSwitchControl(client, name, "Cleaning Mode", controls.Meta{
				Type: "switch",

				Order:    18,
				ReadOnly: false,
				Title:    controls.MultilingualText{"en": `"Wet cleaning" Mode`, "ru": `Режим "Влажная уборка"`},
			}),
			Serial: controls.NewTextControl(client, name, "Serial", controls.Meta{
				Type: "text",

				Order:    19,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Серийный номер`},
			}),
		}

		instanceWbMwacV293 = &WbMwacV293{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceWbMwacV293
}
