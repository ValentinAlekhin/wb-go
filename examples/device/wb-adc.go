package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbAdcControls struct {
	A1          *control.ValueControl
	A2          *control.ValueControl
	A3          *control.ValueControl
	A4          *control.ValueControl
	Vin         *control.ValueControl
	V33         *control.ValueControl
	V50         *control.ValueControl
	VbusDebug   *control.ValueControl
	VbusNetwork *control.ValueControl
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

func NewWbAdc(client *mqtt.Client) *WbAdc {
	onceWbAdc.Do(func() {
		name := "wb-adc"

		controlList := &WbAdcControls{
			A1: control.NewValueControl(client, name, "A1", control.Meta{
				Type: "voltage",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			A2: control.NewValueControl(client, name, "A2", control.Meta{
				Type: "voltage",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			A3: control.NewValueControl(client, name, "A3", control.Meta{
				Type: "voltage",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			A4: control.NewValueControl(client, name, "A4", control.Meta{
				Type: "voltage",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Vin: control.NewValueControl(client, name, "Vin", control.Meta{
				Type: "voltage",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			V33: control.NewValueControl(client, name, "V3_3", control.Meta{
				Type: "voltage",

				Order:    5,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			V50: control.NewValueControl(client, name, "V5_0", control.Meta{
				Type: "voltage",

				Order:    6,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			VbusDebug: control.NewValueControl(client, name, "Vbus_debug", control.Meta{
				Type: "voltage",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			VbusNetwork: control.NewValueControl(client, name, "Vbus_network", control.Meta{
				Type: "voltage",

				Order:    8,
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
