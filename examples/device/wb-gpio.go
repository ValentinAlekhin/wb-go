package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbGpioControls struct {
	A1Out  *control.SwitchControl
	A2Out  *control.SwitchControl
	A3Out  *control.SwitchControl
	A4Out  *control.SwitchControl
	A1In   *control.SwitchControl
	A2In   *control.SwitchControl
	A3In   *control.SwitchControl
	A4In   *control.SwitchControl
	C5VOut *control.SwitchControl
	W1In   *control.SwitchControl
	W2In   *control.SwitchControl
	VOut   *control.SwitchControl
}

type WbGpio struct {
	name     string
	Controls *WbGpioControls
}

func (w *WbGpio) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbGpio) GetControlsInfo() []control.ControlInfo {
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
	onceWbGpio     sync.Once
	instanceWbGpio *WbGpio
)

func NewWbGpio(client *mqtt.Client) *WbGpio {
	onceWbGpio.Do(func() {
		name := "wb-gpio"

		controlList := &WbGpioControls{
			A1Out: control.NewSwitchControl(client, name, "A1_OUT", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			A2Out: control.NewSwitchControl(client, name, "A2_OUT", control.Meta{
				Type: "switch",

				Order:    2,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			A3Out: control.NewSwitchControl(client, name, "A3_OUT", control.Meta{
				Type: "switch",

				Order:    3,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			A4Out: control.NewSwitchControl(client, name, "A4_OUT", control.Meta{
				Type: "switch",

				Order:    4,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			A1In: control.NewSwitchControl(client, name, "A1_IN", control.Meta{
				Type: "switch",

				Order:    5,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			A2In: control.NewSwitchControl(client, name, "A2_IN", control.Meta{
				Type: "switch",

				Order:    6,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			A3In: control.NewSwitchControl(client, name, "A3_IN", control.Meta{
				Type: "switch",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			A4In: control.NewSwitchControl(client, name, "A4_IN", control.Meta{
				Type: "switch",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			C5VOut: control.NewSwitchControl(client, name, "5V_OUT", control.Meta{
				Type: "switch",

				Order:    9,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			W1In: control.NewSwitchControl(client, name, "W1_IN", control.Meta{
				Type: "switch",

				Order:    10,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			W2In: control.NewSwitchControl(client, name, "W2_IN", control.Meta{
				Type: "switch",

				Order:    11,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			VOut: control.NewSwitchControl(client, name, "V_OUT", control.Meta{
				Type: "switch",

				Order:    12,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
		}

		instanceWbGpio = &WbGpio{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbGpio
}
