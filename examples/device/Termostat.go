package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type TermostatControls struct {
	R01Ts161Lock     *control.SwitchControl
	R01Ts161Mode     *control.SwitchControl
	R01Ts161Setpoint *control.RangeControl
}

type Termostat struct {
	name     string
	Controls *TermostatControls
}

func (w *Termostat) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceTermostat     sync.Once
	instanceTermostat *Termostat
)

func NewTermostat(client mqtt.ClientInterface) *Termostat {
	onceTermostat.Do(func() {
		name := "Termostat"

		controlList := &TermostatControls{
			R01Ts161Lock: control.NewSwitchControl(client, name, "R01-TS16-1-lock", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			R01Ts161Mode: control.NewSwitchControl(client, name, "R01-TS16-1-mode", control.Meta{
				Type: "switch",

				Order:    2,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			R01Ts161Setpoint: control.NewRangeControl(client, name, "R01-TS16-1-setpoint", control.Meta{
				Type: "range",

				Max: 30,

				Order:    3,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
		}

		instanceTermostat = &Termostat{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceTermostat
}
