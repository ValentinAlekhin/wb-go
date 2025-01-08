package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMrm2Mini61Controls struct {
	Input1        *control.SwitchControl
	Input1Counter *control.ValueControl
	Input2        *control.SwitchControl
	Input2Counter *control.ValueControl
	K1            *control.SwitchControl
	K2            *control.SwitchControl
	Serial        *control.TextControl
}

type WbMrm2Mini61 struct {
	name     string
	Controls *WbMrm2Mini61Controls
}

func (w *WbMrm2Mini61) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbMrm2Mini61     sync.Once
	instanceWbMrm2Mini61 *WbMrm2Mini61
)

func NewWbMrm2Mini61(client mqtt.ClientInterface) *WbMrm2Mini61 {
	onceWbMrm2Mini61.Do(func() {
		name := "wb-mrm2-mini_61"

		controlList := &WbMrm2Mini61Controls{
			Input1: control.NewSwitchControl(client, name, "Input 1", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 1`},
			}),
			Input1Counter: control.NewValueControl(client, name, "Input 1 counter", control.Meta{
				Type: "value",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 1`},
			}),
			Input2: control.NewSwitchControl(client, name, "Input 2", control.Meta{
				Type: "switch",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 2`},
			}),
			Input2Counter: control.NewValueControl(client, name, "Input 2 counter", control.Meta{
				Type: "value",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 2`},
			}),
			K1: control.NewSwitchControl(client, name, "K1", control.Meta{
				Type: "switch",

				Order:    5,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			K2: control.NewSwitchControl(client, name, "K2", control.Meta{
				Type: "switch",

				Order:    6,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			Serial: control.NewTextControl(client, name, "Serial", control.Meta{
				Type: "text",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Серийный номер`},
			}),
		}

		instanceWbMrm2Mini61 = &WbMrm2Mini61{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbMrm2Mini61
}
