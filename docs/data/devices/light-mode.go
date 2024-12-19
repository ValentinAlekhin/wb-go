package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type LightModeControls struct {
	Enabled *SwitchControl
	State   *ValueControl
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
		name := "light-mode"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "")
		controls := &LightModeControls{
			Enabled: NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "enabled")),
			State:   NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "state")),
		}

		instanceLightMode = &LightMode{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceLightMode
}
