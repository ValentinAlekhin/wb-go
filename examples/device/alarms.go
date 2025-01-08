package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type AlarmsControls struct {
	Log *control.TextControl
}

type Alarms struct {
	name     string
	Controls *AlarmsControls
}

func (w *Alarms) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceAlarms     sync.Once
	instanceAlarms *Alarms
)

func NewAlarms(client mqtt.ClientInterface) *Alarms {
	onceAlarms.Do(func() {
		name := "alarms"

		controlList := &AlarmsControls{
			Log: control.NewTextControl(client, name, "log", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Log`, "ru": `Лог`},
			}),
		}

		instanceAlarms = &Alarms{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceAlarms
}
