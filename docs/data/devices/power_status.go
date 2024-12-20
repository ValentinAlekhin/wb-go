package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type PowerstatusControls struct {
	Vin              *controls.ValueControl
	WorkingOnBattery *controls.SwitchControl
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
		deviceName := fmt.Sprintf("%s_%s", "power", "status")
		controlList := &PowerstatusControls{
			Vin:              controls.NewValueControl(client, deviceName, "Vin"),
			WorkingOnBattery: controls.NewSwitchControl(client, deviceName, "working on battery"),
		}

		instancePowerstatus = &Powerstatus{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instancePowerstatus
}
