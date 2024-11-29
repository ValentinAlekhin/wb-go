package devices

import (
	"fmt"
	"sync"
	"wb-go/pkg/mqtt"
)

type WbAdcControls struct {
	A1          *ValueControl
	A2          *ValueControl
	A3          *ValueControl
	A4          *ValueControl
	Vin         *ValueControl
	V33         *ValueControl
	V50         *ValueControl
	VbusDebug   *ValueControl
	VbusNetwork *ValueControl
}

type WbAdc struct {
	Name          string
	ModbusAddress int32
	Controls      *WbAdcControls
}

var (
	onceWbAdc     sync.Once
	instanceWbAdc *WbAdc
)

func NewWbAdc(client *mqtt.Client) *WbAdc {
	onceWbAdc.Do(func() {
		name := "wb-adc"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "")
		controls := &WbAdcControls{
			A1:          NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "A1")),
			A2:          NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "A2")),
			A3:          NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "A3")),
			A4:          NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "A4")),
			Vin:         NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Vin")),
			V33:         NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "V3_3")),
			V50:         NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "V5_0")),
			VbusDebug:   NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Vbus_debug")),
			VbusNetwork: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Vbus_network")),
		}

		instanceWbAdc = &WbAdc{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceWbAdc
}
