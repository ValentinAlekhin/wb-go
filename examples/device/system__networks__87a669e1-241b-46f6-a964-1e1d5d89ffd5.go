package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5Controls struct {
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

type SystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5 struct {
	name     string
	Controls *SystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5Controls
}

func (w *SystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5     sync.Once
	instanceSystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5 *SystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5
)

func NewSystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5(client mqtt.ClientInterface) *SystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5 {
	onceSystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5.Do(func() {
		name := "system__networks__87a669e1-241b-46f6-a964-1e1d5d89ffd5"

		controlList := &SystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5Controls{
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
				Title:    control.MultilingualText{"en": `Down`},
			}),
		}

		instanceSystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5 = &SystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworks87A669E1241B46F6A9641E1D5D89Ffd5
}
