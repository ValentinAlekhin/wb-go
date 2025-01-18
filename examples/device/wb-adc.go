package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbAdcControls struct {
	A1 *control.ValueControl
}

type WbAdc struct {
	name     string
	Controls *WbAdcControls
}

func (w *WbAdc) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbAdc     sync.Once
	instanceWbAdc *WbAdc
)

func NewWbAdc(client mqtt.ClientInterface) *WbAdc {
	onceWbAdc.Do(func() {
		name := "wb-adc"

		controlList := &WbAdcControls{
			A1: control.NewValueControl(client, name, "A1", control.Meta{
				Type: "voltage",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
		}

		instanceWbAdc = &WbAdc{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbAdc
}
