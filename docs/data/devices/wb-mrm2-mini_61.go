package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbMrm2Mini61Controls struct {
	Input1        *controls.SwitchControl
	Input1Counter *controls.ValueControl
	Input2        *controls.SwitchControl
	Input2Counter *controls.ValueControl
	K1            *controls.SwitchControl
	K2            *controls.SwitchControl
	Serial        *controls.TextControl
}

type WbMrm2Mini61 struct {
	Name     string
	Address  string
	Controls *WbMrm2Mini61Controls
}

func (w *WbMrm2Mini61) GetControlsInfo() []controls.ControlInfo {
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
	onceWbMrm2Mini61     sync.Once
	instanceWbMrm2Mini61 *WbMrm2Mini61
)

func NewWbMrm2Mini61(client *mqtt.Client) *WbMrm2Mini61 {
	onceWbMrm2Mini61.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "wb-mrm2-mini", "61")
		controlList := &WbMrm2Mini61Controls{
			Input1:        controls.NewSwitchControl(client, deviceName, "Input 1"),
			Input1Counter: controls.NewValueControl(client, deviceName, "Input 1 counter"),
			Input2:        controls.NewSwitchControl(client, deviceName, "Input 2"),
			Input2Counter: controls.NewValueControl(client, deviceName, "Input 2 counter"),
			K1:            controls.NewSwitchControl(client, deviceName, "K1"),
			K2:            controls.NewSwitchControl(client, deviceName, "K2"),
			Serial:        controls.NewTextControl(client, deviceName, "Serial"),
		}

		instanceWbMrm2Mini61 = &WbMrm2Mini61{
			Name:     deviceName,
			Address:  "61",
			Controls: controlList,
		}
	})

	return instanceWbMrm2Mini61
}
