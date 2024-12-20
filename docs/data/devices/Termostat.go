package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type TermostatControls struct {
	R01Ts161Lock     *controls.SwitchControl
	R01Ts161Mode     *controls.SwitchControl
	R01Ts161Setpoint *controls.RangeControl
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
		deviceName := fmt.Sprintf("%s_%s", "Termostat", "")
		controlList := &TermostatControls{
			R01Ts161Lock:     controls.NewSwitchControl(client, deviceName, "R01-TS16-1-lock"),
			R01Ts161Mode:     controls.NewSwitchControl(client, deviceName, "R01-TS16-1-mode"),
			R01Ts161Setpoint: controls.NewRangeControl(client, deviceName, "R01-TS16-1-setpoint"),
		}

		instanceTermostat = &Termostat{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceTermostat
}
