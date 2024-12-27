package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbGpioControls struct {
	A1Out  *controls.SwitchControl
	A2Out  *controls.SwitchControl
	A3Out  *controls.SwitchControl
	A4Out  *controls.SwitchControl
	A1In   *controls.SwitchControl
	A2In   *controls.SwitchControl
	A3In   *controls.SwitchControl
	A4In   *controls.SwitchControl
	C5VOut *controls.SwitchControl
	W1In   *controls.SwitchControl
	W2In   *controls.SwitchControl
	VOut   *controls.SwitchControl
}

type WbGpio struct {
	name     string
	device   string
	address  string
	Controls *WbGpioControls
}

func (w *WbGpio) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbGpio) GetControlsInfo() []controls.ControlInfo {
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
	onceWbGpio     sync.Once
	instanceWbGpio *WbGpio
)

func NewWbGpio(client *mqtt.Client) *WbGpio {
	onceWbGpio.Do(func() {
		device := "wb-gpio"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &WbGpioControls{
			A1Out: controls.NewSwitchControl(client, name, "A1_OUT", controls.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			A2Out: controls.NewSwitchControl(client, name, "A2_OUT", controls.Meta{
				Type: "switch",

				Order:    2,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			A3Out: controls.NewSwitchControl(client, name, "A3_OUT", controls.Meta{
				Type: "switch",

				Order:    3,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			A4Out: controls.NewSwitchControl(client, name, "A4_OUT", controls.Meta{
				Type: "switch",

				Order:    4,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			A1In: controls.NewSwitchControl(client, name, "A1_IN", controls.Meta{
				Type: "switch",

				Order:    5,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			A2In: controls.NewSwitchControl(client, name, "A2_IN", controls.Meta{
				Type: "switch",

				Order:    6,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			A3In: controls.NewSwitchControl(client, name, "A3_IN", controls.Meta{
				Type: "switch",

				Order:    7,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			A4In: controls.NewSwitchControl(client, name, "A4_IN", controls.Meta{
				Type: "switch",

				Order:    8,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			C5VOut: controls.NewSwitchControl(client, name, "5V_OUT", controls.Meta{
				Type: "switch",

				Order:    9,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			W1In: controls.NewSwitchControl(client, name, "W1_IN", controls.Meta{
				Type: "switch",

				Order:    10,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			W2In: controls.NewSwitchControl(client, name, "W2_IN", controls.Meta{
				Type: "switch",

				Order:    11,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			VOut: controls.NewSwitchControl(client, name, "V_OUT", controls.Meta{
				Type: "switch",

				Order:    12,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
		}

		instanceWbGpio = &WbGpio{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceWbGpio
}
