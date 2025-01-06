package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbAdccontrols struct {
	A1          *control.ValueControl
	A2          *control.ValueControl
	A3          *control.ValueControl
	A4          *control.ValueControl
	Vin         *control.ValueControl
	V33         *control.ValueControl
	V50         *control.ValueControl
	VbusDebug   *control.ValueControl
	VbusNetwork *control.ValueControl
}

type WbAdc struct {
	name     string
	device   string
	address  string
	Controls *WbAdccontrols
}

func (w *WbAdc) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbAdc) GetControlsInfo() []control.ControlInfo {
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
	onceWbAdc     sync.Once
	instanceWbAdc *WbAdc
)

func NewWbAdc(client *mqtt.Client) *WbAdc {
	onceWbAdc.Do(func() {
		device := "wb-adc"
		address := ""
		name := device

		controlList := &WbAdccontrols{
			A1: control.NewValueControl(client, name, "A1", control.Meta{
				Type: "voltage",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			A2: control.NewValueControl(client, name, "A2", control.Meta{
				Type: "voltage",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			A3: control.NewValueControl(client, name, "A3", control.Meta{
				Type: "voltage",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			A4: control.NewValueControl(client, name, "A4", control.Meta{
				Type: "voltage",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Vin: control.NewValueControl(client, name, "Vin", control.Meta{
				Type: "voltage",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			V33: control.NewValueControl(client, name, "V3_3", control.Meta{
				Type: "voltage",

				Order:    5,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			V50: control.NewValueControl(client, name, "V5_0", control.Meta{
				Type: "voltage",

				Order:    6,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			VbusDebug: control.NewValueControl(client, name, "Vbus_debug", control.Meta{
				Type: "voltage",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			VbusNetwork: control.NewValueControl(client, name, "Vbus_network", control.Meta{
				Type: "voltage",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
		}

		instanceWbAdc = &WbAdc{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceWbAdc
}
