package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbLed150Controls struct {
	Cct1             *SwitchControl
	Cct1Temperature  *RangeControl
	Cct1Brightness   *RangeControl
	Cct2             *SwitchControl
	Cct2Temperature  *RangeControl
	Cct2Brightness   *RangeControl
	BoardTemperature *ValueControl
	AllowedPower     *ValueControl
	Overcurrent      *SwitchControl
	Input1           *SwitchControl
	Input2           *SwitchControl
	Input2Counter    *ValueControl
	Input3           *SwitchControl
	Input3Counter    *ValueControl
	Input4           *SwitchControl
	Input4Counter    *ValueControl
	Serial           *TextControl
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
		name := "wb-led"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "150")
		controls := &WbLed150Controls{
			Cct1:             NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "CCT1")),
			Cct1Temperature:  NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "CCT1 Temperature")),
			Cct1Brightness:   NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "CCT1 Brightness")),
			Cct2:             NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "CCT2")),
			Cct2Temperature:  NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "CCT2 Temperature")),
			Cct2Brightness:   NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "CCT2 Brightness")),
			BoardTemperature: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Board Temperature")),
			AllowedPower:     NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Allowed Power")),
			Overcurrent:      NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Overcurrent")),
			Input1:           NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 1")),
			Input2:           NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 2")),
			Input2Counter:    NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 2 Counter")),
			Input3:           NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 3")),
			Input3Counter:    NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 3 Counter")),
			Input4:           NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 4")),
			Input4Counter:    NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 4 Counter")),
			Serial:           NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Serial")),
		}

		instanceWbLed150 = &WbLed150{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceWbLed150
}
