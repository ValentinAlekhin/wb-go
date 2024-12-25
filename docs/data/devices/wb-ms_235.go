package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
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
	Name     string
	Address  string
	Controls *WbMs235Controls
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
			Address:  "235",
			Controls: controlList,
		}
	})

	return instanceWbMs235
}
