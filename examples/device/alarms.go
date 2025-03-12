package device

import (
	"context"
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

func NewAlarms(ctx context.Context, client mqtt.ClientInterface) *Alarms {
	onceAlarms.Do(func() {
		name := "alarms"

		controlList := &AlarmsControls{
			Log: control.NewTextControl(ctx, client, name, "log", control.Meta{
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
