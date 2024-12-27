package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
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
	name     string
	device   string
	address  string
	Controls *WbAdcControls
}

func (w *WbAdc) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
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
		device := "wb-adc"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &WbAdcControls{
			A1: controls.NewValueControl(client, name, "A1", controls.Meta{
				Type: "voltage",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			A2: controls.NewValueControl(client, name, "A2", controls.Meta{
				Type: "voltage",

				Order:    1,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			A3: controls.NewValueControl(client, name, "A3", controls.Meta{
				Type: "voltage",

				Order:    2,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			A4: controls.NewValueControl(client, name, "A4", controls.Meta{
				Type: "voltage",

				Order:    3,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			Vin: controls.NewValueControl(client, name, "Vin", controls.Meta{
				Type: "voltage",

				Order:    4,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			V33: controls.NewValueControl(client, name, "V3_3", controls.Meta{
				Type: "voltage",

				Order:    5,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			V50: controls.NewValueControl(client, name, "V5_0", controls.Meta{
				Type: "voltage",

				Order:    6,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			VbusDebug: controls.NewValueControl(client, name, "Vbus_debug", controls.Meta{
				Type: "voltage",

				Order:    7,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			VbusNetwork: controls.NewValueControl(client, name, "Vbus_network", controls.Meta{
				Type: "voltage",

				Order:    8,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
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
