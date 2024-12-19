package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMs235Controls struct {
	Temperature     *ValueControl
	Humidity        *ValueControl
	AirQualityVoc   *ValueControl
	AirQualityIndex *ValueControl
	Illuminance     *ValueControl
	ExternalSensor1 *ValueControl
	ExternalSensor2 *ValueControl
	Serial          *TextControl
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
		name := "wb-ms"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "235")
		controls := &WbMs235Controls{
			Temperature:     NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Temperature")),
			Humidity:        NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Humidity")),
			AirQualityVoc:   NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Air Quality (VOC)")),
			AirQualityIndex: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Air Quality Index")),
			Illuminance:     NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Illuminance")),
			ExternalSensor1: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "External Sensor 1")),
			ExternalSensor2: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "External Sensor 2")),
			Serial:          NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Serial")),
		}

		instanceWbMs235 = &WbMs235{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceWbMs235
}
