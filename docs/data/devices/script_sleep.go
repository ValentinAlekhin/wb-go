package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type ScriptsleepControls struct {
	Current          *ValueControl
	Enable           *SwitchControl
	State            *TextControl
	Target           *RangeControl
	Zone1RelayStatus *SwitchControl
	Zone1Status      *ValueControl
}

type Scriptsleep struct {
	Name     string
	Controls *ScriptsleepControls
}

var (
	onceScriptsleep     sync.Once
	instanceScriptsleep *Scriptsleep
)

func NewScriptsleep(client *mqtt.Client) *Scriptsleep {
	onceScriptsleep.Do(func() {
		name := "script"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "sleep")
		controls := &ScriptsleepControls{
			Current:          NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "current")),
			Enable:           NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "enable")),
			State:            NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "state")),
			Target:           NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "target")),
			Zone1RelayStatus: NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "zone1_relay_status")),
			Zone1Status:      NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "zone1_status")),
		}

		instanceScriptsleep = &Scriptsleep{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceScriptsleep
}
