package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMs235Controls struct {
	Temperature     *controls.ValueControl
	Humidity        *controls.ValueControl
	AirQualityVoc   *controls.ValueControl
	AirQualityIndex *controls.ValueControl
	Illuminance     *controls.ValueControl
	ExternalSensor1 *controls.ValueControl
	ExternalSensor2 *controls.ValueControl
	Serial          *controls.TextControl
}

type WbMs235 struct {
	Name     string
	Controls *WbMs235Controls
}

var (
	onceWbMs235     sync.Once
	instanceWbMs235 *WbMs235
)

func NewWbMs235(client *mqtt.Client) *WbMs235 {
	onceWbMs235.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "wb-ms", "235")
		controlList := &WbMs235Controls{
			Temperature:     controls.NewValueControl(client, deviceName, "Temperature"),
			Humidity:        controls.NewValueControl(client, deviceName, "Humidity"),
			AirQualityVoc:   controls.NewValueControl(client, deviceName, "Air Quality (VOC)"),
			AirQualityIndex: controls.NewValueControl(client, deviceName, "Air Quality Index"),
			Illuminance:     controls.NewValueControl(client, deviceName, "Illuminance"),
			ExternalSensor1: controls.NewValueControl(client, deviceName, "External Sensor 1"),
			ExternalSensor2: controls.NewValueControl(client, deviceName, "External Sensor 2"),
			Serial:          controls.NewTextControl(client, deviceName, "Serial"),
		}

		instanceWbMs235 = &WbMs235{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceWbMs235
}
