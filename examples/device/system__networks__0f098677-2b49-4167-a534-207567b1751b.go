package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworks0F0986772B494167A534207567B1751BControls struct {
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

type SystemNetworks0F0986772B494167A534207567B1751B struct {
	name     string
	Controls *SystemNetworks0F0986772B494167A534207567B1751BControls
}

func (w *SystemNetworks0F0986772B494167A534207567B1751B) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworks0F0986772B494167A534207567B1751B     sync.Once
	instanceSystemNetworks0F0986772B494167A534207567B1751B *SystemNetworks0F0986772B494167A534207567B1751B
)

func NewSystemNetworks0F0986772B494167A534207567B1751B(client mqtt.ClientInterface) *SystemNetworks0F0986772B494167A534207567B1751B {
	onceSystemNetworks0F0986772B494167A534207567B1751B.Do(func() {
		name := "system__networks__0f098677-2b49-4167-a534-207567b1751b"

		controlList := &SystemNetworks0F0986772B494167A534207567B1751BControls{
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

		instanceSystemNetworks0F0986772B494167A534207567B1751B = &SystemNetworks0F0986772B494167A534207567B1751B{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworks0F0986772B494167A534207567B1751B
}
