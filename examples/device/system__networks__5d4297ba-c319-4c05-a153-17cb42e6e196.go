package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworks5D4297BaC3194C05A15317Cb42E6E196Controls struct {
	Name *control.TextControl
}

type SystemNetworks5D4297BaC3194C05A15317Cb42E6E196 struct {
	name     string
	Controls *SystemNetworks5D4297BaC3194C05A15317Cb42E6E196Controls
}

func (w *SystemNetworks5D4297BaC3194C05A15317Cb42E6E196) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworks5D4297BaC3194C05A15317Cb42E6E196     sync.Once
	instanceSystemNetworks5D4297BaC3194C05A15317Cb42E6E196 *SystemNetworks5D4297BaC3194C05A15317Cb42E6E196
)

func NewSystemNetworks5D4297BaC3194C05A15317Cb42E6E196(client mqtt.ClientInterface) *SystemNetworks5D4297BaC3194C05A15317Cb42E6E196 {
	onceSystemNetworks5D4297BaC3194C05A15317Cb42E6E196.Do(func() {
		name := "system__networks__5d4297ba-c319-4c05-a153-17cb42e6e196"

		controlList := &SystemNetworks5D4297BaC3194C05A15317Cb42E6E196Controls{
			Name: control.NewTextControl(client, name, "Name", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
		}

		instanceSystemNetworks5D4297BaC3194C05A15317Cb42E6E196 = &SystemNetworks5D4297BaC3194C05A15317Cb42E6E196{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworks5D4297BaC3194C05A15317Cb42E6E196
}
