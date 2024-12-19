package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type PowerstatusControls struct {
	Vin              *ValueControl
	WorkingOnBattery *SwitchControl
}

type Powerstatus struct {
	Name     string
	Controls *PowerstatusControls
}

var (
	oncePowerstatus     sync.Once
	instancePowerstatus *Powerstatus
)

func NewPowerstatus(client *mqtt.Client) *Powerstatus {
	oncePowerstatus.Do(func() {
		name := "power"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "status")
		controls := &PowerstatusControls{
			Vin:              NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Vin")),
			WorkingOnBattery: NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "working on battery")),
		}

		instancePowerstatus = &Powerstatus{
			Name:     name,
			Controls: controls,
		}
	})

	return instancePowerstatus
}
