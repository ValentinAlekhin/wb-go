package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type LightModeControls struct {
	Enabled *control.SwitchControl
	State   *control.ValueControl
}

type LightMode struct {
	name     string
	Controls *LightModeControls
}

func (w *LightMode) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceLightMode     sync.Once
	instanceLightMode *LightMode
)

func NewLightMode(client *mqtt.Client) *LightMode {
	onceLightMode.Do(func() {
		name := "light-mode"

		controlList := &LightModeControls{
			Enabled: control.NewSwitchControl(client, name, "enabled", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			State: control.NewValueControl(client, name, "state", control.Meta{
				Type: "value",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `State`},
			}),
		}

		instanceLightMode = &LightMode{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceLightMode
}
