package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
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
	name     string
	device   string
	address  string
	Controls *WbMrm2Mini61Controls
}

func (w *WbMrm2Mini61) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
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
		device := "wb-mrm2-mini"
		address := "61"
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &WbMrm2Mini61Controls{
			Input1: controls.NewSwitchControl(client, name, "Input 1", controls.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 1`},
			}),
			Input1Counter: controls.NewValueControl(client, name, "Input 1 counter", controls.Meta{
				Type: "value",

				Order:    2,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 1`},
			}),
			Input2: controls.NewSwitchControl(client, name, "Input 2", controls.Meta{
				Type: "switch",

				Order:    3,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 2`},
			}),
			Input2Counter: controls.NewValueControl(client, name, "Input 2 counter", controls.Meta{
				Type: "value",

				Order:    4,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 2`},
			}),
			K1: controls.NewSwitchControl(client, name, "K1", controls.Meta{
				Type: "switch",

				Order:    5,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			K2: controls.NewSwitchControl(client, name, "K2", controls.Meta{
				Type: "switch",

				Order:    6,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			Serial: controls.NewTextControl(client, name, "Serial", controls.Meta{
				Type: "text",

				Order:    7,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Серийный номер`},
			}),
		}

		instanceWbMrm2Mini61 = &WbMrm2Mini61{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceWbMrm2Mini61
}
