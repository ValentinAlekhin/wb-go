package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworks91F1C71D2D974675886FEcbe52B8451EControls struct {
	Name *control.TextControl
}

type SystemNetworks91F1C71D2D974675886FEcbe52B8451E struct {
	name     string
	Controls *SystemNetworks91F1C71D2D974675886FEcbe52B8451EControls
}

func (w *SystemNetworks91F1C71D2D974675886FEcbe52B8451E) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworks91F1C71D2D974675886FEcbe52B8451E     sync.Once
	instanceSystemNetworks91F1C71D2D974675886FEcbe52B8451E *SystemNetworks91F1C71D2D974675886FEcbe52B8451E
)

func NewSystemNetworks91F1C71D2D974675886FEcbe52B8451E(client mqtt.ClientInterface) *SystemNetworks91F1C71D2D974675886FEcbe52B8451E {
	onceSystemNetworks91F1C71D2D974675886FEcbe52B8451E.Do(func() {
		name := "system__networks__91f1c71d-2d97-4675-886f-ecbe52b8451e"

		controlList := &SystemNetworks91F1C71D2D974675886FEcbe52B8451EControls{
			Name: control.NewTextControl(client, name, "Name", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
		}

		instanceSystemNetworks91F1C71D2D974675886FEcbe52B8451E = &SystemNetworks91F1C71D2D974675886FEcbe52B8451E{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworks91F1C71D2D974675886FEcbe52B8451E
}
