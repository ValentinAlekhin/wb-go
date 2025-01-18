package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMs235Controls struct {
	Temperature *control.ValueControl
}

type WbMs235 struct {
	name     string
	Controls *WbMs235Controls
}

func (w *WbMs235) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbMs235     sync.Once
	instanceWbMs235 *WbMs235
)

func NewWbMs235(client mqtt.ClientInterface) *WbMs235 {
	onceWbMs235.Do(func() {
		name := "wb-ms_235"

		controlList := &WbMs235Controls{
			Temperature: control.NewValueControl(client, name, "Temperature", control.Meta{
				Type: "temperature",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Температура`},
			}),
		}

		instanceWbMs235 = &WbMs235{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbMs235
}
