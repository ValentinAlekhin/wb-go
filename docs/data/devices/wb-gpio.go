package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbGpioControls struct {
	A1Out  *SwitchControl
	A2Out  *SwitchControl
	A3Out  *SwitchControl
	A4Out  *SwitchControl
	A1In   *SwitchControl
	A2In   *SwitchControl
	A3In   *SwitchControl
	A4In   *SwitchControl
	C5VOut *SwitchControl
	W1In   *SwitchControl
	W2In   *SwitchControl
	VOut   *SwitchControl
}

type WbGpio struct {
	Name     string
	Controls *WbGpioControls
}

var (
	onceWbGpio     sync.Once
	instanceWbGpio *WbGpio
)

func NewWbGpio(client *mqtt.Client) *WbGpio {
	onceWbGpio.Do(func() {
		name := "wb-gpio"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "")
		controls := &WbGpioControls{
			A1Out:  NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "A1_OUT")),
			A2Out:  NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "A2_OUT")),
			A3Out:  NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "A3_OUT")),
			A4Out:  NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "A4_OUT")),
			A1In:   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "A1_IN")),
			A2In:   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "A2_IN")),
			A3In:   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "A3_IN")),
			A4In:   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "A4_IN")),
			C5VOut: NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "5V_OUT")),
			W1In:   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "W1_IN")),
			W2In:   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "W2_IN")),
			VOut:   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "V_OUT")),
		}

		instanceWbGpio = &WbGpio{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceWbGpio
}
