package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworksC3E384059C174155Ad70664311B49066Controls struct {
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

type SystemNetworksC3E384059C174155Ad70664311B49066 struct {
	name     string
	Controls *SystemNetworksC3E384059C174155Ad70664311B49066Controls
}

func (w *SystemNetworksC3E384059C174155Ad70664311B49066) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworksC3E384059C174155Ad70664311B49066     sync.Once
	instanceSystemNetworksC3E384059C174155Ad70664311B49066 *SystemNetworksC3E384059C174155Ad70664311B49066
)

func NewSystemNetworksC3E384059C174155Ad70664311B49066(client *mqtt.Client) *SystemNetworksC3E384059C174155Ad70664311B49066 {
	onceSystemNetworksC3E384059C174155Ad70664311B49066.Do(func() {
		name := "system__networks__c3e38405-9c17-4155-ad70-664311b49066"

		controlList := &SystemNetworksC3E384059C174155Ad70664311B49066Controls{
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

		instanceSystemNetworksC3E384059C174155Ad70664311B49066 = &SystemNetworksC3E384059C174155Ad70664311B49066{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworksC3E384059C174155Ad70664311B49066
}
