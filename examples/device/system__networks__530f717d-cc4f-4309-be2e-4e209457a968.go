package device

import (
	"context"
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworks530F717DCc4F4309Be2E4E209457A968Controls struct {
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

type SystemNetworks530F717DCc4F4309Be2E4E209457A968 struct {
	name     string
	Controls *SystemNetworks530F717DCc4F4309Be2E4E209457A968Controls
}

func (w *SystemNetworks530F717DCc4F4309Be2E4E209457A968) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworks530F717DCc4F4309Be2E4E209457A968     sync.Once
	instanceSystemNetworks530F717DCc4F4309Be2E4E209457A968 *SystemNetworks530F717DCc4F4309Be2E4E209457A968
)

func NewSystemNetworks530F717DCc4F4309Be2E4E209457A968(ctx context.Context, client mqtt.ClientInterface) *SystemNetworks530F717DCc4F4309Be2E4E209457A968 {
	onceSystemNetworks530F717DCc4F4309Be2E4E209457A968.Do(func() {
		name := "system__networks__530f717d-cc4f-4309-be2e-4e209457a968"

		controlList := &SystemNetworks530F717DCc4F4309Be2E4E209457A968Controls{
			Name: control.NewTextControl(ctx, client, name, "Name", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Uuid: control.NewTextControl(ctx, client, name, "UUID", control.Meta{
				Type: "text",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Type: control.NewTextControl(ctx, client, name, "Type", control.Meta{
				Type: "text",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Active: control.NewSwitchControl(ctx, client, name, "Active", control.Meta{
				Type: "switch",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Device: control.NewTextControl(ctx, client, name, "Device", control.Meta{
				Type: "text",

				Order:    5,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			State: control.NewTextControl(ctx, client, name, "State", control.Meta{
				Type: "text",

				Order:    6,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Address: control.NewTextControl(ctx, client, name, "Address", control.Meta{
				Type: "text",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Connectivity: control.NewSwitchControl(ctx, client, name, "Connectivity", control.Meta{
				Type: "switch",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			UpDown: control.NewPushbuttonControl(ctx, client, name, "UpDown", control.Meta{
				Type: "pushbutton",

				Order:    12,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Down`},
			}),
		}

		instanceSystemNetworks530F717DCc4F4309Be2E4E209457A968 = &SystemNetworks530F717DCc4F4309Be2E4E209457A968{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworks530F717DCc4F4309Be2E4E209457A968
}
