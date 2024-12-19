package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type TermostatControls struct {
	R01Ts161Lock     *SwitchControl
	R01Ts161Mode     *SwitchControl
	R01Ts161Setpoint *RangeControl
}

type Termostat struct {
	Name     string
	Controls *TermostatControls
}

var (
	onceTermostat     sync.Once
	instanceTermostat *Termostat
)

func NewTermostat(client *mqtt.Client) *Termostat {
	onceTermostat.Do(func() {
		name := "Termostat"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "")
		controls := &TermostatControls{
			R01Ts161Lock:     NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "R01-TS16-1-lock")),
			R01Ts161Mode:     NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "R01-TS16-1-mode")),
			R01Ts161Setpoint: NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "R01-TS16-1-setpoint")),
		}

		instanceTermostat = &Termostat{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceTermostat
}
