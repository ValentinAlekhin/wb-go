package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbLed106Controls struct {
	Overcurrent *control.SwitchControl
}

type WbLed106 struct {
	name     string
	Controls *WbLed106Controls
}

func (w *WbLed106) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbLed106     sync.Once
	instanceWbLed106 *WbLed106
)

func NewWbLed106(client mqtt.ClientInterface) *WbLed106 {
	onceWbLed106.Do(func() {
		name := "wb-led_106"

		controlList := &WbLed106Controls{
			Overcurrent: control.NewSwitchControl(client, name, "Overcurrent", control.Meta{
				Type: "switch",

				Order:    12,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Перегрузка по току`},
			}),
		}

		instanceWbLed106 = &WbLed106{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbLed106
}
