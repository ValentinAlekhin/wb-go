package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemControls struct {
	BatchNo *control.TextControl
}

type System struct {
	name     string
	Controls *SystemControls
}

func (w *System) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystem     sync.Once
	instanceSystem *System
)

func NewSystem(client mqtt.ClientInterface) *System {
	onceSystem.Do(func() {
		name := "system"

		controlList := &SystemControls{
			BatchNo: control.NewTextControl(client, name, "Batch No", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Batch No`, "ru": `Номер партии`},
			}),
		}

		instanceSystem = &System{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystem
}
