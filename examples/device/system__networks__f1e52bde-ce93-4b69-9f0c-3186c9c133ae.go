package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworksF1E52BdeCe934B699F0C3186C9C133AeControls struct {
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

type SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae struct {
	name     string
	Controls *SystemNetworksF1E52BdeCe934B699F0C3186C9C133AeControls
}

func (w *SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae     sync.Once
	instanceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae *SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae
)

func NewSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae(client mqtt.ClientInterface) *SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae {
	onceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae.Do(func() {
		name := "system__networks__f1e52bde-ce93-4b69-9f0c-3186c9c133ae"

		controlList := &SystemNetworksF1E52BdeCe934B699F0C3186C9C133AeControls{
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

		instanceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae = &SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae
}
