package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type HwmonControls struct {
	BoardTemperature *controls.ValueControl
	CpuTemperature   *controls.ValueControl
}

type Hwmon struct {
	Name     string
	Controls *HwmonControls
}

var (
	onceHwmon     sync.Once
	instanceHwmon *Hwmon
)

func NewHwmon(client *mqtt.Client) *Hwmon {
	onceHwmon.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "hwmon", "")
		controlList := &HwmonControls{
			BoardTemperature: controls.NewValueControl(client, deviceName, "Board Temperature"),
			CpuTemperature:   controls.NewValueControl(client, deviceName, "CPU Temperature"),
		}

		instanceHwmon = &Hwmon{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceHwmon
}
