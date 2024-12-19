package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbrulesControls struct {
	RuleDebugging *SwitchControl
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
		name := "wbrules"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "")
		controls := &WbrulesControls{
			RuleDebugging: NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Rule debugging")),
		}

		instanceWbrules = &Wbrules{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceWbrules
}
