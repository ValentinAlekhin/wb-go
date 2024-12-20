package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbGpioControls struct {
	A1Out  *controls.SwitchControl
	A2Out  *controls.SwitchControl
	A3Out  *controls.SwitchControl
	A4Out  *controls.SwitchControl
	A1In   *controls.SwitchControl
	A2In   *controls.SwitchControl
	A3In   *controls.SwitchControl
	A4In   *controls.SwitchControl
	C5VOut *controls.SwitchControl
	W1In   *controls.SwitchControl
	W2In   *controls.SwitchControl
	VOut   *controls.SwitchControl
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
		deviceName := fmt.Sprintf("%s_%s", "wb-gpio", "")
		controlList := &WbGpioControls{
			A1Out:  controls.NewSwitchControl(client, deviceName, "A1_OUT"),
			A2Out:  controls.NewSwitchControl(client, deviceName, "A2_OUT"),
			A3Out:  controls.NewSwitchControl(client, deviceName, "A3_OUT"),
			A4Out:  controls.NewSwitchControl(client, deviceName, "A4_OUT"),
			A1In:   controls.NewSwitchControl(client, deviceName, "A1_IN"),
			A2In:   controls.NewSwitchControl(client, deviceName, "A2_IN"),
			A3In:   controls.NewSwitchControl(client, deviceName, "A3_IN"),
			A4In:   controls.NewSwitchControl(client, deviceName, "A4_IN"),
			C5VOut: controls.NewSwitchControl(client, deviceName, "5V_OUT"),
			W1In:   controls.NewSwitchControl(client, deviceName, "W1_IN"),
			W2In:   controls.NewSwitchControl(client, deviceName, "W2_IN"),
			VOut:   controls.NewSwitchControl(client, deviceName, "V_OUT"),
		}

		instanceWbGpio = &WbGpio{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceWbGpio
}
