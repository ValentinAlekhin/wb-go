package devices

import (
	"fmt"
	"sync"
	"wb-go/pkg/mqtt"
)

type WbMr6Cu145Controls struct {
	K1     *SwitchControl
	K2     *SwitchControl
	K3     *SwitchControl
	K4     *SwitchControl
	K5     *SwitchControl
	K6     *SwitchControl
	Serial *TextControl
}

type WbMr6Cu145 struct {
	Name          string
	ModbusAddress int32
	Controls      *WbMr6Cu145Controls
}

var (
	onceWbMr6Cu145     sync.Once
	instanceWbMr6Cu145 *WbMr6Cu145
)

func NewWbMr6Cu145(client *mqtt.Client) *WbMr6Cu145 {
	onceWbMr6Cu145.Do(func() {
		name := "wb-mr6cu"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "145")
		controls := &WbMr6Cu145Controls{
			K1:     NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "K1")),
			K2:     NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "K2")),
			K3:     NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "K3")),
			K4:     NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "K4")),
			K5:     NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "K5")),
			K6:     NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "K6")),
			Serial: NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Serial")),
		}

		instanceWbMr6Cu145 = &WbMr6Cu145{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceWbMr6Cu145
}
