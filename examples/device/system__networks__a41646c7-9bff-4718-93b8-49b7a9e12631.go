package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworksA41646C79Bff471893B849B7A9E12631Controls struct {
	Name *control.TextControl
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
		}

		instanceSystemNetworksA41646C79Bff471893B849B7A9E12631 = &SystemNetworksA41646C79Bff471893B849B7A9E12631{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworksA41646C79Bff471893B849B7A9E12631
}
