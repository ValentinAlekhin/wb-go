package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMrm2Mini61Controls struct {
	Input1 *control.SwitchControl
}

type WbMrm2Mini61 struct {
	name     string
	Controls *WbMrm2Mini61Controls
}

func (w *WbMrm2Mini61) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbMrm2Mini61     sync.Once
	instanceWbMrm2Mini61 *WbMrm2Mini61
)

func NewWbMrm2Mini61(client mqtt.ClientInterface) *WbMrm2Mini61 {
	onceWbMrm2Mini61.Do(func() {
		name := "wb-mrm2-mini_61"

		controlList := &WbMrm2Mini61Controls{
			Input1: control.NewSwitchControl(client, name, "Input 1", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 1`},
			}),
		}

		instanceWbMrm2Mini61 = &WbMrm2Mini61{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbMrm2Mini61
}
