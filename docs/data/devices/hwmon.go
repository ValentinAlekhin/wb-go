package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type HwmonControls struct {
	BoardTemperature *ValueControl
	CpuTemperature   *ValueControl
}

type Hwmon struct {
	Name          string
	ModbusAddress int32
	Controls      *HwmonControls
}

var (
	onceHwmon     sync.Once
	instanceHwmon *Hwmon
)

func NewHwmon(client *mqtt.Client) *Hwmon {
	onceHwmon.Do(func() {
		name := "hwmon"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "")
		controls := &HwmonControls{
			BoardTemperature: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Board Temperature")),
			CpuTemperature:   NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "CPU Temperature")),
		}

		instanceHwmon = &Hwmon{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceHwmon
}
