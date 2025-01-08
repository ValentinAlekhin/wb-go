package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbrulesControls struct {
	RuleDebugging *control.SwitchControl
}

type Wbrules struct {
	name     string
	Controls *WbrulesControls
}

func (w *Wbrules) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbrules     sync.Once
	instanceWbrules *Wbrules
)

func NewWbrules(client mqtt.ClientInterface) *Wbrules {
	onceWbrules.Do(func() {
		name := "wbrules"

		controlList := &WbrulesControls{
			RuleDebugging: control.NewSwitchControl(client, name, "Rule debugging", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Rule debugging`, "ru": `Отладка правил`},
			}),
		}

		instanceWbrules = &Wbrules{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbrules
}
