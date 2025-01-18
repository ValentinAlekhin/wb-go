package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbGpioControls struct {
	A1Out *control.SwitchControl
}

type WbGpio struct {
	name     string
	Controls *WbGpioControls
}

func (w *WbGpio) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbGpio     sync.Once
	instanceWbGpio *WbGpio
)

func NewWbGpio(client mqtt.ClientInterface) *WbGpio {
	onceWbGpio.Do(func() {
		name := "wb-gpio"

		controlList := &WbGpioControls{
			A1Out: control.NewSwitchControl(client, name, "A1_OUT", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
		}

		instanceWbGpio = &WbGpio{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbGpio
}
