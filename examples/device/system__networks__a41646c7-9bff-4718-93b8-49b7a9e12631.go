package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworksA41646C79Bff471893B849B7A9E12631Controls struct {
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

type SystemNetworksA41646C79Bff471893B849B7A9E12631 struct {
	name     string
	Controls *SystemNetworksA41646C79Bff471893B849B7A9E12631Controls
}

func (w *SystemNetworksA41646C79Bff471893B849B7A9E12631) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworksA41646C79Bff471893B849B7A9E12631     sync.Once
	instanceSystemNetworksA41646C79Bff471893B849B7A9E12631 *SystemNetworksA41646C79Bff471893B849B7A9E12631
)

func NewSystemNetworksA41646C79Bff471893B849B7A9E12631(client mqtt.ClientInterface) *SystemNetworksA41646C79Bff471893B849B7A9E12631 {
	onceSystemNetworksA41646C79Bff471893B849B7A9E12631.Do(func() {
		name := "system__networks__a41646c7-9bff-4718-93b8-49b7a9e12631"

		controlList := &SystemNetworksA41646C79Bff471893B849B7A9E12631Controls{
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

		instanceSystemNetworksA41646C79Bff471893B849B7A9E12631 = &SystemNetworksA41646C79Bff471893B849B7A9E12631{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworksA41646C79Bff471893B849B7A9E12631
}
