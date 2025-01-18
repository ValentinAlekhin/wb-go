package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMs235Controls struct {
	Temperature     *control.ValueControl
	Humidity        *control.ValueControl
	AirQualityVoc   *control.ValueControl
	AirQualityIndex *control.ValueControl
	Illuminance     *control.ValueControl
	ExternalSensor1 *control.ValueControl
	ExternalSensor2 *control.ValueControl
	Serial          *control.TextControl
}

type WbMs235 struct {
	name     string
	Controls *WbMs235Controls
}

func (w *WbMs235) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceWbMs235     sync.Once
	instanceWbMs235 *WbMs235
)

func NewWbMs235(client mqtt.ClientInterface) *WbMs235 {
	onceWbMs235.Do(func() {
		name := "wb-ms_235"

		controlList := &WbMs235Controls{
			Temperature: control.NewValueControl(client, name, "Temperature", control.Meta{
				Type: "temperature",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Температура`},
			}),
			Humidity: control.NewValueControl(client, name, "Humidity", control.Meta{
				Type: "rel_humidity",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Влажность`},
			}),
			AirQualityVoc: control.NewValueControl(client, name, "Air Quality (VOC)", control.Meta{
				Type:  "value",
				Units: "ppb",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Качество воздуха (VOC)`},
			}),
			AirQualityIndex: control.NewValueControl(client, name, "Air Quality Index", control.Meta{
				Type: "value",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Индекс качества воздуха (AQI)`},
			}),
			Illuminance: control.NewValueControl(client, name, "Illuminance", control.Meta{
				Type: "lux",

				Order:    5,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Освещенность`},
			}),
			ExternalSensor1: control.NewValueControl(client, name, "External Sensor 1", control.Meta{
				Type: "temperature",

				Order:    6,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Датчик температуры 1`},
			}),
			ExternalSensor2: control.NewValueControl(client, name, "External Sensor 2", control.Meta{
				Type: "temperature",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Датчик температуры 2`},
			}),
			Serial: control.NewTextControl(client, name, "Serial", control.Meta{
				Type: "text",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Серийный номер`},
			}),
		}

		instanceWbMs235 = &WbMs235{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbMs235
}
