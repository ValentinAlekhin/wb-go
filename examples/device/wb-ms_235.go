package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbMs235controls struct {
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
	device   string
	address  string
	Controls *WbMs235controls
}

func (w *WbMs235) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbMs235) GetControlsInfo() []control.ControlInfo {
	var infoList []control.ControlInfo

	// Получаем значение и тип структуры Controls
	controlsValue := reflect.ValueOf(w.Controls).Elem()
	controlsType := controlsValue.Type()

	// Проходимся по всем полям структуры Controls
	for i := 0; i < controlsValue.NumField(); i++ {
		field := controlsValue.Field(i)

		// Проверяем, что поле является указателем и не nil
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			// Проверяем, реализует ли поле метод GetInfo
			method := field.MethodByName("GetInfo")
			if method.IsValid() {
				// Вызываем метод GetInfo
				info := method.Call(nil)[0].Interface().(control.ControlInfo)
				infoList = append(infoList, info)
			} else {
				fmt.Printf("Field %s does not implement GetInfo\n", controlsType.Field(i).Name)
			}
		}
	}

	return infoList
}

var (
	onceWbMs235     sync.Once
	instanceWbMs235 *WbMs235
)

func NewWbMs235(client *mqtt.Client) *WbMs235 {
	onceWbMs235.Do(func() {
		device := "wb-ms"
		address := "235"
		name := fmt.Sprintf("%s_%s", device, address)

		controlList := &WbMs235controls{
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
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceWbMs235
}
