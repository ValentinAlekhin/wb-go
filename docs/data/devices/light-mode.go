package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type LightModeControls struct {
	Enabled *controls.SwitchControl
	State   *controls.ValueControl
}

type LightMode struct {
	Name     string
	Controls *LightModeControls
}

var (
	onceLightMode     sync.Once
	instanceLightMode *LightMode
)

func NewLightMode(client *mqtt.Client) *LightMode {
	onceLightMode.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "light-mode", "")
		controlList := &LightModeControls{
			Enabled: controls.NewSwitchControl(client, deviceName, "enabled"),
			State:   controls.NewValueControl(client, deviceName, "state"),
		}

		instanceLightMode = &LightMode{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceLightMode
}
