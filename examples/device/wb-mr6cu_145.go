package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMr6Cu145Controls struct {
	K1 *control.SwitchControl
}

type WbMr6Cu145 struct {
	name     string
	Controls *WbMr6Cu145Controls
}

func (w *WbMr6Cu145) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbMr6Cu145     sync.Once
	instanceWbMr6Cu145 *WbMr6Cu145
)

func NewWbMr6Cu145(client mqtt.ClientInterface) *WbMr6Cu145 {
	onceWbMr6Cu145.Do(func() {
		name := "wb-mr6cu_145"

		controlList := &WbMr6Cu145Controls{
			K1: control.NewSwitchControl(client, name, "K1", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
		}

		instanceWbMr6Cu145 = &WbMr6Cu145{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbMr6Cu145
}
