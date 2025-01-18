package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbLed150Controls struct {
	Cct1 *control.SwitchControl
}

type WbLed150 struct {
	name     string
	Controls *WbLed150Controls
}

func (w *WbLed150) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbLed150     sync.Once
	instanceWbLed150 *WbLed150
)

func NewWbLed150(client mqtt.ClientInterface) *WbLed150 {
	onceWbLed150.Do(func() {
		name := "wb-led_150"

		controlList := &WbLed150Controls{
			Cct1: control.NewSwitchControl(client, name, "CCT1", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Лента CCT1`},
			}),
		}

		instanceWbLed150 = &WbLed150{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbLed150
}
