package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMswV4151Controls struct {
	Temperature *control.ValueControl
}

type WbMswV4151 struct {
	name     string
	Controls *WbMswV4151Controls
}

func (w *WbMswV4151) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbMswV4151     sync.Once
	instanceWbMswV4151 *WbMswV4151
)

func NewWbMswV4151(client mqtt.ClientInterface) *WbMswV4151 {
	onceWbMswV4151.Do(func() {
		name := "wb-msw-v4_151"

		controlList := &WbMswV4151Controls{
			Temperature: control.NewValueControl(client, name, "Temperature", control.Meta{
				Type:  "value",
				Units: "deg C",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Температура`},
			}),
		}

		instanceWbMswV4151 = &WbMswV4151{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbMswV4151
}
