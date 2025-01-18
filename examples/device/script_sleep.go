package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type ScriptSleepControls struct {
	Current *control.ValueControl
}

type ScriptSleep struct {
	name     string
	Controls *ScriptSleepControls
}

func (w *ScriptSleep) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceScriptSleep     sync.Once
	instanceScriptSleep *ScriptSleep
)

func NewScriptSleep(client mqtt.ClientInterface) *ScriptSleep {
	onceScriptSleep.Do(func() {
		name := "script_sleep"

		controlList := &ScriptSleepControls{
			Current: control.NewValueControl(client, name, "current", control.Meta{
				Type: "temperature",

				Order:    16,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Current Temperature`, "ru": `Current Temperature`},
			}),
		}

		instanceScriptSleep = &ScriptSleep{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceScriptSleep
}
