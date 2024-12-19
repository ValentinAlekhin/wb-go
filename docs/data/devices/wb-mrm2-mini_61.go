package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMrm2Mini61Controls struct {
	Input1        *SwitchControl
	Input1Counter *ValueControl
	Input2        *SwitchControl
	Input2Counter *ValueControl
	K1            *SwitchControl
	K2            *SwitchControl
	Serial        *TextControl
}

type WbMrm2Mini61 struct {
	Name     string
	Controls *WbMrm2Mini61Controls
}

var (
	onceWbMrm2Mini61     sync.Once
	instanceWbMrm2Mini61 *WbMrm2Mini61
)

func NewWbMrm2Mini61(client *mqtt.Client) *WbMrm2Mini61 {
	onceWbMrm2Mini61.Do(func() {
		name := "wb-mrm2-mini"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "61")
		controls := &WbMrm2Mini61Controls{
			Input1:        NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 1")),
			Input1Counter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 1 counter")),
			Input2:        NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 2")),
			Input2Counter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 2 counter")),
			K1:            NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "K1")),
			K2:            NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "K2")),
			Serial:        NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Serial")),
		}

		instanceWbMrm2Mini61 = &WbMrm2Mini61{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceWbMrm2Mini61
}
