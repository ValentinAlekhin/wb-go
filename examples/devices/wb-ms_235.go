package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
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
	name     string
	device   string
	address  string
	Controls *WbMs235Controls
}

func (w *WbMs235) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbMs235) GetControlsInfo() []controls.ControlInfo {
	var infoList []controls.ControlInfo

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
				info := method.Call(nil)[0].Interface().(controls.ControlInfo)
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
		controlList := &WbMs235Controls{
			Temperature: controls.NewValueControl(client, name, "Temperature", controls.Meta{
				Type: "temperature",

				Order:    1,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Температура`},
			}),
			Humidity: controls.NewValueControl(client, name, "Humidity", controls.Meta{
				Type: "rel_humidity",

				Order:    2,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Влажность`},
			}),
			AirQualityVoc: controls.NewValueControl(client, name, "Air Quality (VOC)", controls.Meta{
				Type:  "value",
				Units: "ppb",

				Order:    3,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Качество воздуха (VOC)`},
			}),
			AirQualityIndex: controls.NewValueControl(client, name, "Air Quality Index", controls.Meta{
				Type: "value",

				Order:    4,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Индекс качества воздуха (AQI)`},
			}),
			Illuminance: controls.NewValueControl(client, name, "Illuminance", controls.Meta{
				Type: "lux",

				Order:    5,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Освещенность`},
			}),
			ExternalSensor1: controls.NewValueControl(client, name, "External Sensor 1", controls.Meta{
				Type: "temperature",

				Order:    6,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Датчик температуры 1`},
			}),
			ExternalSensor2: controls.NewValueControl(client, name, "External Sensor 2", controls.Meta{
				Type: "temperature",

				Order:    7,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Датчик температуры 2`},
			}),
			Serial: controls.NewTextControl(client, name, "Serial", controls.Meta{
				Type: "text",

				Order:    8,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Серийный номер`},
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
