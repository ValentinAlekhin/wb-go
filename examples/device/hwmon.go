package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type HwmonControls struct {
	BoardTemperature *control.ValueControl
}

type Hwmon struct {
	name     string
	Controls *HwmonControls
}

func (w *Hwmon) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceHwmon     sync.Once
	instanceHwmon *Hwmon
)

func NewHwmon(client mqtt.ClientInterface) *Hwmon {
	onceHwmon.Do(func() {
		name := "hwmon"

		controlList := &HwmonControls{
			BoardTemperature: control.NewValueControl(client, name, "Board Temperature", control.Meta{
				Type: "temperature",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
		}

		instanceHwmon = &Hwmon{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceHwmon
}
