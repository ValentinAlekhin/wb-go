package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbAdcControls struct {
	A1          *controls.ValueControl
	A2          *controls.ValueControl
	A3          *controls.ValueControl
	A4          *controls.ValueControl
	Vin         *controls.ValueControl
	V33         *controls.ValueControl
	V50         *controls.ValueControl
	VbusDebug   *controls.ValueControl
	VbusNetwork *controls.ValueControl
}

type WbAdc struct {
	Name     string
	Address  string
	Controls *WbAdcControls
}

func (w *WbAdc) GetControlsInfo() []controls.ControlInfo {
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
	onceWbAdc     sync.Once
	instanceWbAdc *WbAdc
)

func NewWbAdc(client *mqtt.Client) *WbAdc {
	onceWbAdc.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "wb-adc", "")
		controlList := &WbAdcControls{
			A1:          controls.NewValueControl(client, deviceName, "A1"),
			A2:          controls.NewValueControl(client, deviceName, "A2"),
			A3:          controls.NewValueControl(client, deviceName, "A3"),
			A4:          controls.NewValueControl(client, deviceName, "A4"),
			Vin:         controls.NewValueControl(client, deviceName, "Vin"),
			V33:         controls.NewValueControl(client, deviceName, "V3_3"),
			V50:         controls.NewValueControl(client, deviceName, "V5_0"),
			VbusDebug:   controls.NewValueControl(client, deviceName, "Vbus_debug"),
			VbusNetwork: controls.NewValueControl(client, deviceName, "Vbus_network"),
		}

		instanceWbAdc = &WbAdc{
			Name:     deviceName,
			Address:  "",
			Controls: controlList,
		}
	})

	return instanceWbAdc
}
