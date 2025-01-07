package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type ScriptSleepControls struct {
	Current          *control.ValueControl
	Enable           *control.SwitchControl
	State            *control.TextControl
	Target           *control.RangeControl
	Zone1RelayStatus *control.SwitchControl
	Zone1Status      *control.ValueControl
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

func NewScriptSleep(client *mqtt.Client) *ScriptSleep {
	onceScriptSleep.Do(func() {
		name := "script_sleep"

		controlList := &ScriptSleepControls{
			Current: control.NewValueControl(client, name, "current", control.Meta{
				Type: "temperature",

				Order:    16,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Current Temperature`, "ru": `Current Temperature`},
			}),
			Enable: control.NewSwitchControl(client, name, "enable", control.Meta{
				Type: "switch",

				Order:    10,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Enable`, "ru": `Enable`},
			}),
			State: control.NewTextControl(client, name, "state", control.Meta{
				Type: "text",

				Order:    15,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Mode State`, "ru": `State`},
			}),
			Target: control.NewRangeControl(client, name, "target", control.Meta{
				Type: "range",

				Max: 30,
				Min: 14,

				Order:    20,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Target Temperature`, "ru": `Target Temperature`},
			}),
			Zone1RelayStatus: control.NewSwitchControl(client, name, "zone1_relay_status", control.Meta{
				Type: "switch",

				Order:    50,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Zone Relay Status`, "ru": `Zone Relay Status`},
			}),
			Zone1Status: control.NewValueControl(client, name, "zone1_status", control.Meta{
				Type: "temperature",

				Order:    30,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Zone Temperature`, "ru": `Zone Temperature`},
			}),
		}

		instanceScriptSleep = &ScriptSleep{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceScriptSleep
}
