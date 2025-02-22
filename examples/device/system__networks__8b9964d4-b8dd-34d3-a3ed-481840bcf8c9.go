package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9Controls struct {
	Name               *control.TextControl
	Uuid               *control.TextControl
	Type               *control.TextControl
	Active             *control.SwitchControl
	Device             *control.TextControl
	State              *control.TextControl
	Address            *control.TextControl
	Connectivity       *control.SwitchControl
	UpDown             *control.PushbuttonControl
	Operator           *control.TextControl
	SignalQuality      *control.TextControl
	AccessTechnologies *control.TextControl
}

type SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9 struct {
	name     string
	Controls *SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9Controls
}

func (w *SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9     sync.Once
	instanceSystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9 *SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9
)

func NewSystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9(client mqtt.ClientInterface) *SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9 {
	onceSystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9.Do(func() {
		name := "system__networks__8b9964d4-b8dd-34d3-a3ed-481840bcf8c9"

		controlList := &SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9Controls{
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
			Operator: control.NewTextControl(client, name, "Operator", control.Meta{
				Type: "text",

				Order:    9,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			SignalQuality: control.NewTextControl(client, name, "SignalQuality", control.Meta{
				Type: "text",

				Order:    10,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Signal Quality`},
			}),
			AccessTechnologies: control.NewTextControl(client, name, "AccessTechnologies", control.Meta{
				Type: "text",

				Order:    11,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Access Technologies`},
			}),
		}

		instanceSystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9 = &SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9
}
