package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbMr6Cu145Controls struct {
	K1     *controls.SwitchControl
	K2     *controls.SwitchControl
	K3     *controls.SwitchControl
	K4     *controls.SwitchControl
	K5     *controls.SwitchControl
	K6     *controls.SwitchControl
	Serial *controls.TextControl
}

type WbMr6Cu145 struct {
	name     string
	device   string
	address  string
	Controls *WbMr6Cu145Controls
}

func (w *WbMr6Cu145) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbMr6Cu145) GetControlsInfo() []controls.ControlInfo {
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
	onceWbMr6Cu145     sync.Once
	instanceWbMr6Cu145 *WbMr6Cu145
)

func NewWbMr6Cu145(client *mqtt.Client) *WbMr6Cu145 {
	onceWbMr6Cu145.Do(func() {
		device := "wb-mr6cu"
		address := "145"
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &WbMr6Cu145Controls{
			K1: controls.NewSwitchControl(client, name, "K1", controls.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			K2: controls.NewSwitchControl(client, name, "K2", controls.Meta{
				Type: "switch",

				Order:    2,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			K3: controls.NewSwitchControl(client, name, "K3", controls.Meta{
				Type: "switch",

				Order:    3,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			K4: controls.NewSwitchControl(client, name, "K4", controls.Meta{
				Type: "switch",

				Order:    4,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			K5: controls.NewSwitchControl(client, name, "K5", controls.Meta{
				Type: "switch",

				Order:    5,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			K6: controls.NewSwitchControl(client, name, "K6", controls.Meta{
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

		instanceWbMr6Cu145 = &WbMr6Cu145{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceWbMr6Cu145
}
