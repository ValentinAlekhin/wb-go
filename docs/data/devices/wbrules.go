package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbrulesControls struct {
	RuleDebugging *controls.SwitchControl
}

type Wbrules struct {
	Name     string
	Controls *WbrulesControls
}

var (
	onceWbrules     sync.Once
	instanceWbrules *Wbrules
)

func NewWbrules(client *mqtt.Client) *Wbrules {
	onceWbrules.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "wbrules", "")
		controlList := &WbrulesControls{
			RuleDebugging: controls.NewSwitchControl(client, deviceName, "Rule debugging"),
		}

		instanceWbrules = &Wbrules{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceWbrules
}
