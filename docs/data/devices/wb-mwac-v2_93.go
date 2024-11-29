package devices

import (
	"fmt"
	"sync"
	"wb-go/pkg/mqtt"
)

type WbMwacV293Controls struct {
	P2Volume       *ValueControl
	InputF1        *SwitchControl
	InputF1Counter *ValueControl
	InputF2        *SwitchControl
	InputF2Counter *ValueControl
	InputF3        *SwitchControl
	InputF3Counter *ValueControl
	InputF4        *SwitchControl
	InputF4Counter *ValueControl
	InputF5        *SwitchControl
	InputF5Counter *ValueControl
	OutputK2       *SwitchControl
	CleaningMode   *SwitchControl
	P1Volume       *ValueControl
	InputS6        *SwitchControl
	InputS6Counter *ValueControl
	OutputK1       *SwitchControl
	LeakageMode    *SwitchControl
	Serial         *TextControl
}

type WbMwacV293 struct {
	Name          string
	ModbusAddress int32
	Controls      *WbMwacV293Controls
}

var (
	onceWbMwacV293     sync.Once
	instanceWbMwacV293 *WbMwacV293
)

func NewWbMwacV293(client *mqtt.Client) *WbMwacV293 {
	onceWbMwacV293.Do(func() {
		name := "wb-mwac-v2"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "93")
		controls := &WbMwacV293Controls{
			P2Volume:       NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "P2 Volume")),
			InputF1:        NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input F1")),
			InputF1Counter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input F1 Counter")),
			InputF2:        NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input F2")),
			InputF2Counter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input F2 Counter")),
			InputF3:        NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input F3")),
			InputF3Counter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input F3 Counter")),
			InputF4:        NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input F4")),
			InputF4Counter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input F4 Counter")),
			InputF5:        NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input F5")),
			InputF5Counter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input F5 Counter")),
			OutputK2:       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Output K2")),
			CleaningMode:   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Cleaning Mode")),
			P1Volume:       NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "P1 Volume")),
			InputS6:        NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input S6")),
			InputS6Counter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input S6 Counter")),
			OutputK1:       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Output K1")),
			LeakageMode:    NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Leakage Mode")),
			Serial:         NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Serial")),
		}

		instanceWbMwacV293 = &WbMwacV293{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceWbMwacV293
}
