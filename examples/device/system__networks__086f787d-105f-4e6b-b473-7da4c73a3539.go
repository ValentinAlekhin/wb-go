package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworks086F787D105F4E6BB4737Da4C73A3539Controls struct {
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

type SystemNetworks086F787D105F4E6BB4737Da4C73A3539 struct {
	name     string
	Controls *SystemNetworks086F787D105F4E6BB4737Da4C73A3539Controls
}

func (w *SystemNetworks086F787D105F4E6BB4737Da4C73A3539) GetInfo() basedevice.Info {
	controlsInfo := w.GetControlsInfo()

	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: controlsInfo,
	}
}

func (w *SystemNetworks086F787D105F4E6BB4737Da4C73A3539) GetControlsInfo() []control.Info {
	return basedevice.GetControlsInfo(w.Controls)
}

var (
	onceSystemNetworks086F787D105F4E6BB4737Da4C73A3539     sync.Once
	instanceSystemNetworks086F787D105F4E6BB4737Da4C73A3539 *SystemNetworks086F787D105F4E6BB4737Da4C73A3539
)

func NewSystemNetworks086F787D105F4E6BB4737Da4C73A3539(client *mqtt.Client) *SystemNetworks086F787D105F4E6BB4737Da4C73A3539 {
	onceSystemNetworks086F787D105F4E6BB4737Da4C73A3539.Do(func() {
		name := "system__networks__086f787d-105f-4e6b-b473-7da4c73a3539"

		controlList := &SystemNetworks086F787D105F4E6BB4737Da4C73A3539Controls{
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

		instanceSystemNetworks086F787D105F4E6BB4737Da4C73A3539 = &SystemNetworks086F787D105F4E6BB4737Da4C73A3539{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworks086F787D105F4E6BB4737Da4C73A3539
}
