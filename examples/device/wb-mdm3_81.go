package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMdm381Controls struct {
	Input5Counter *control.ValueControl
}

type WbMdm381 struct {
	name     string
	Controls *WbMdm381Controls
}

func (w *WbMdm381) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbMdm381     sync.Once
	instanceWbMdm381 *WbMdm381
)

func NewWbMdm381(client mqtt.ClientInterface) *WbMdm381 {
	onceWbMdm381.Do(func() {
		name := "wb-mdm3_81"

		controlList := &WbMdm381Controls{
			Input5Counter: control.NewValueControl(client, name, "Input 5 counter", control.Meta{
				Type: "value",

				Order:    18,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 5`},
			}),
		}

		instanceWbMdm381 = &WbMdm381{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbMdm381
}
