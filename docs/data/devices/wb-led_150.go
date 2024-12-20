package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbLed150Controls struct {
	Cct1             *controls.SwitchControl
	Cct1Temperature  *controls.RangeControl
	Cct1Brightness   *controls.RangeControl
	Cct2             *controls.SwitchControl
	Cct2Temperature  *controls.RangeControl
	Cct2Brightness   *controls.RangeControl
	BoardTemperature *controls.ValueControl
	AllowedPower     *controls.ValueControl
	Overcurrent      *controls.SwitchControl
	Input1           *controls.SwitchControl
	Input2           *controls.SwitchControl
	Input2Counter    *controls.ValueControl
	Input3           *controls.SwitchControl
	Input3Counter    *controls.ValueControl
	Input4           *controls.SwitchControl
	Input4Counter    *controls.ValueControl
	Serial           *controls.TextControl
}

type WbLed150 struct {
	Name     string
	Controls *WbLed150Controls
}

var (
	onceWbLed150     sync.Once
	instanceWbLed150 *WbLed150
)

func NewWbLed150(client *mqtt.Client) *WbLed150 {
	onceWbLed150.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "wb-led", "150")
		controlList := &WbLed150Controls{
			Cct1:             controls.NewSwitchControl(client, deviceName, "CCT1"),
			Cct1Temperature:  controls.NewRangeControl(client, deviceName, "CCT1 Temperature"),
			Cct1Brightness:   controls.NewRangeControl(client, deviceName, "CCT1 Brightness"),
			Cct2:             controls.NewSwitchControl(client, deviceName, "CCT2"),
			Cct2Temperature:  controls.NewRangeControl(client, deviceName, "CCT2 Temperature"),
			Cct2Brightness:   controls.NewRangeControl(client, deviceName, "CCT2 Brightness"),
			BoardTemperature: controls.NewValueControl(client, deviceName, "Board Temperature"),
			AllowedPower:     controls.NewValueControl(client, deviceName, "Allowed Power"),
			Overcurrent:      controls.NewSwitchControl(client, deviceName, "Overcurrent"),
			Input1:           controls.NewSwitchControl(client, deviceName, "Input 1"),
			Input2:           controls.NewSwitchControl(client, deviceName, "Input 2"),
			Input2Counter:    controls.NewValueControl(client, deviceName, "Input 2 Counter"),
			Input3:           controls.NewSwitchControl(client, deviceName, "Input 3"),
			Input3Counter:    controls.NewValueControl(client, deviceName, "Input 3 Counter"),
			Input4:           controls.NewSwitchControl(client, deviceName, "Input 4"),
			Input4Counter:    controls.NewValueControl(client, deviceName, "Input 4 Counter"),
			Serial:           controls.NewTextControl(client, deviceName, "Serial"),
		}

		instanceWbLed150 = &WbLed150{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceWbLed150
}
