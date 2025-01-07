package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbMr6Cu145Controls struct {
	K1     *control.SwitchControl
	K2     *control.SwitchControl
	K3     *control.SwitchControl
	K4     *control.SwitchControl
	K5     *control.SwitchControl
	K6     *control.SwitchControl
	Serial *control.TextControl
}

type WbMr6Cu145 struct {
	name     string
	Controls *WbMr6Cu145Controls
}

func (w *WbMr6Cu145) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbMr6Cu145) GetControlsInfo() []control.ControlInfo {
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
	onceWbMr6Cu145     sync.Once
	instanceWbMr6Cu145 *WbMr6Cu145
)

func NewWbMr6Cu145(client *mqtt.Client) *WbMr6Cu145 {
	onceWbMr6Cu145.Do(func() {
		name := "wb-mr6cu_145"

		controlList := &WbMr6Cu145Controls{
			K1: control.NewSwitchControl(client, name, "K1", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			K2: control.NewSwitchControl(client, name, "K2", control.Meta{
				Type: "switch",

				Order:    2,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			K3: control.NewSwitchControl(client, name, "K3", control.Meta{
				Type: "switch",

				Order:    3,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			K4: control.NewSwitchControl(client, name, "K4", control.Meta{
				Type: "switch",

				Order:    4,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			K5: control.NewSwitchControl(client, name, "K5", control.Meta{
				Type: "switch",

				Order:    5,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			K6: control.NewSwitchControl(client, name, "K6", control.Meta{
				Type: "switch",

				Order:    6,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			Serial: control.NewTextControl(client, name, "Serial", control.Meta{
				Type: "text",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Серийный номер`},
			}),
		}

		instanceWbMr6Cu145 = &WbMr6Cu145{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbMr6Cu145
}
