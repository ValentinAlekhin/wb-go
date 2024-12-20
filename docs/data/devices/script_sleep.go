package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type ScriptsleepControls struct {
	Current          *controls.ValueControl
	Enable           *controls.SwitchControl
	State            *controls.TextControl
	Target           *controls.RangeControl
	Zone1RelayStatus *controls.SwitchControl
	Zone1Status      *controls.ValueControl
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
		deviceName := fmt.Sprintf("%s_%s", "script", "sleep")
		controlList := &ScriptsleepControls{
			Current:          controls.NewValueControl(client, deviceName, "current"),
			Enable:           controls.NewSwitchControl(client, deviceName, "enable"),
			State:            controls.NewTextControl(client, deviceName, "state"),
			Target:           controls.NewRangeControl(client, deviceName, "target"),
			Zone1RelayStatus: controls.NewSwitchControl(client, deviceName, "zone1_relay_status"),
			Zone1Status:      controls.NewValueControl(client, deviceName, "zone1_status"),
		}

		instanceScriptsleep = &Scriptsleep{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceScriptsleep
}
