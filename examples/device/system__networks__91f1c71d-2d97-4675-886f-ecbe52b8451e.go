package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworks91F1C71D2D974675886FEcbe52B8451EControls struct {
	Name         *control.TextControl
	Uuid         *control.TextControl
	Type         *control.TextControl
	Active       *control.SwitchControl
	Device       *control.TextControl
	State        *control.TextControl
	Address      *control.TextControl
	Connectivity *control.SwitchControl
	UpDown       *control.PushbuttonControl
}

type SystemNetworks91F1C71D2D974675886FEcbe52B8451E struct {
	name     string
	Controls *SystemNetworks91F1C71D2D974675886FEcbe52B8451EControls
}

func (w *SystemNetworks91F1C71D2D974675886FEcbe52B8451E) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworks91F1C71D2D974675886FEcbe52B8451E     sync.Once
	instanceSystemNetworks91F1C71D2D974675886FEcbe52B8451E *SystemNetworks91F1C71D2D974675886FEcbe52B8451E
)

func NewSystemNetworks91F1C71D2D974675886FEcbe52B8451E(client *mqtt.Client) *SystemNetworks91F1C71D2D974675886FEcbe52B8451E {
	onceSystemNetworks91F1C71D2D974675886FEcbe52B8451E.Do(func() {
		name := "system__networks__91f1c71d-2d97-4675-886f-ecbe52b8451e"

		controlList := &SystemNetworks91F1C71D2D974675886FEcbe52B8451EControls{
			Name: control.NewTextControl(client, name, "Name", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Uuid: control.NewTextControl(client, name, "UUID", control.Meta{
				Type: "text",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Type: control.NewTextControl(client, name, "Type", control.Meta{
				Type: "text",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Active: control.NewSwitchControl(client, name, "Active", control.Meta{
				Type: "switch",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Device: control.NewTextControl(client, name, "Device", control.Meta{
				Type: "text",

				Order:    5,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			State: control.NewTextControl(client, name, "State", control.Meta{
				Type: "text",

				Order:    6,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Address: control.NewTextControl(client, name, "Address", control.Meta{
				Type: "text",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Connectivity: control.NewSwitchControl(client, name, "Connectivity", control.Meta{
				Type: "switch",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			UpDown: control.NewPushbuttonControl(client, name, "UpDown", control.Meta{
				Type: "pushbutton",

				Order:    12,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Up`},
			}),
		}

		instanceSystemNetworks91F1C71D2D974675886FEcbe52B8451E = &SystemNetworks91F1C71D2D974675886FEcbe52B8451E{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworks91F1C71D2D974675886FEcbe52B8451E
}
