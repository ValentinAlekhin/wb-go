package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMwacV293Controls struct {
	InputF1Counter *control.ValueControl
}

type WbMwacV293 struct {
	name     string
	Controls *WbMwacV293Controls
}

func (w *WbMwacV293) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbMwacV293     sync.Once
	instanceWbMwacV293 *WbMwacV293
)

func NewWbMwacV293(client mqtt.ClientInterface) *WbMwacV293 {
	onceWbMwacV293.Do(func() {
		name := "wb-mwac-v2_93"

		controlList := &WbMwacV293Controls{
			InputF1Counter: control.NewValueControl(client, name, "Input F1 Counter", control.Meta{
				Type: "value",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик входа F1`},
			}),
		}

		instanceWbMwacV293 = &WbMwacV293{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbMwacV293
}
