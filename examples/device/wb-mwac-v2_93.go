package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMwacV293Controls struct {
	InputF1Counter *control.ValueControl
	InputF2        *control.SwitchControl
	InputF2Counter *control.ValueControl
	InputF4        *control.SwitchControl
	InputF4Counter *control.ValueControl
	InputF5        *control.SwitchControl
	InputF5Counter *control.ValueControl
	CleaningMode   *control.SwitchControl
	P1Volume       *control.ValueControl
	P2Volume       *control.ValueControl
	InputF1        *control.SwitchControl
	InputF3        *control.SwitchControl
	InputF3Counter *control.ValueControl
	InputS6        *control.SwitchControl
	InputS6Counter *control.ValueControl
	OutputK1       *control.SwitchControl
	OutputK2       *control.SwitchControl
	LeakageMode    *control.SwitchControl
	Serial         *control.TextControl
}

type WbMwacV293 struct {
	name     string
	Controls *WbMwacV293Controls
}

func (w *WbMwacV293) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbMwacV293     sync.Once
	instanceWbMwacV293 *WbMwacV293
)

func NewWbMwacV293(client mqtt.ClientInterface) *WbMwacV293 {
	onceWbMwacV293.Do(func() {
		name := "wb-mwac-v2_93"

		controlList := &WbMwacV293Controls{
			InputF1Counter: control.NewValueControl(client, name, "Input F1 Counter", control.Meta{
				Type: "value",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик входа F1`},
			}),
			InputF2: control.NewSwitchControl(client, name, "Input F2", control.Meta{
				Type: "switch",

				Order:    5,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход F2`},
			}),
			InputF2Counter: control.NewValueControl(client, name, "Input F2 Counter", control.Meta{
				Type: "value",

				Order:    6,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик входа F2`},
			}),
			InputF4: control.NewSwitchControl(client, name, "Input F4", control.Meta{
				Type: "switch",

				Order:    9,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход F4`},
			}),
			InputF4Counter: control.NewValueControl(client, name, "Input F4 Counter", control.Meta{
				Type: "value",

				Order:    10,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик входа F4`},
			}),
			InputF5: control.NewSwitchControl(client, name, "Input F5", control.Meta{
				Type: "switch",

				Order:    11,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход F5`},
			}),
			InputF5Counter: control.NewValueControl(client, name, "Input F5 Counter", control.Meta{
				Type: "value",

				Order:    12,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик входа F5`},
			}),
			CleaningMode: control.NewSwitchControl(client, name, "Cleaning Mode", control.Meta{
				Type: "switch",

				Order:    18,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `"Wet cleaning" Mode`, "ru": `Режим "Влажная уборка"`},
			}),
			P1Volume: control.NewValueControl(client, name, "P1 Volume", control.Meta{
				Type:  "value",
				Units: "m^3",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик объема P1`},
			}),
			P2Volume: control.NewValueControl(client, name, "P2 Volume", control.Meta{
				Type:  "value",
				Units: "m^3",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик объема P2`},
			}),
			InputF1: control.NewSwitchControl(client, name, "Input F1", control.Meta{
				Type: "switch",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход F1`},
			}),
			InputF3: control.NewSwitchControl(client, name, "Input F3", control.Meta{
				Type: "switch",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход F3`},
			}),
			InputF3Counter: control.NewValueControl(client, name, "Input F3 Counter", control.Meta{
				Type: "value",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик входа F3`},
			}),
			InputS6: control.NewSwitchControl(client, name, "Input S6", control.Meta{
				Type: "switch",

				Order:    13,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход S6`},
			}),
			InputS6Counter: control.NewValueControl(client, name, "Input S6 Counter", control.Meta{
				Type: "value",

				Order:    14,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик входа S6`},
			}),
			OutputK1: control.NewSwitchControl(client, name, "Output K1", control.Meta{
				Type: "switch",

				Order:    15,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Выход K1`},
			}),
			OutputK2: control.NewSwitchControl(client, name, "Output K2", control.Meta{
				Type: "switch",

				Order:    16,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Выход K2`},
			}),
			LeakageMode: control.NewSwitchControl(client, name, "Leakage Mode", control.Meta{
				Type: "switch",

				Order:    17,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `"Leakage" Mode`, "ru": `Режим "Протечка"`},
			}),
			Serial: control.NewTextControl(client, name, "Serial", control.Meta{
				Type: "text",

				Order:    19,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Серийный номер`},
			}),
		}

		instanceWbMwacV293 = &WbMwacV293{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbMwacV293
}
