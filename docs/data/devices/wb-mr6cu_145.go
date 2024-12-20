package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMr6Cu145Controls struct {
	K1     *controls.SwitchControl
	K2     *controls.SwitchControl
	K3     *controls.SwitchControl
	K4     *controls.SwitchControl
	K5     *controls.SwitchControl
	K6     *controls.SwitchControl
	Serial *controls.TextControl
}

type WbMr6Cu145 struct {
	Name     string
	Controls *WbMr6Cu145Controls
}

var (
	onceWbMr6Cu145     sync.Once
	instanceWbMr6Cu145 *WbMr6Cu145
)

func NewWbMr6Cu145(client *mqtt.Client) *WbMr6Cu145 {
	onceWbMr6Cu145.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "wb-mr6cu", "145")
		controlList := &WbMr6Cu145Controls{
			K1:     controls.NewSwitchControl(client, deviceName, "K1"),
			K2:     controls.NewSwitchControl(client, deviceName, "K2"),
			K3:     controls.NewSwitchControl(client, deviceName, "K3"),
			K4:     controls.NewSwitchControl(client, deviceName, "K4"),
			K5:     controls.NewSwitchControl(client, deviceName, "K5"),
			K6:     controls.NewSwitchControl(client, deviceName, "K6"),
			Serial: controls.NewTextControl(client, deviceName, "Serial"),
		}

		instanceWbMr6Cu145 = &WbMr6Cu145{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceWbMr6Cu145
}
