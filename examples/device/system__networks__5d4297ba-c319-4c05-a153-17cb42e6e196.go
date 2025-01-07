package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type SystemNetworks5D4297BaC3194C05A15317Cb42E6E196Controls struct {
	Name               *control.TextControl
	Uuid               *control.TextControl
	Type               *control.TextControl
	Active             *control.SwitchControl
	Device             *control.TextControl
	State              *control.TextControl
	Address            *control.TextControl
	Connectivity       *control.SwitchControl
	UpDown             *control.PushbuttonControl
	Operator           *control.TextControl
	SignalQuality      *control.TextControl
	AccessTechnologies *control.TextControl
}

type SystemNetworks5D4297BaC3194C05A15317Cb42E6E196 struct {
	name     string
	Controls *SystemNetworks5D4297BaC3194C05A15317Cb42E6E196Controls
}

func (w *SystemNetworks5D4297BaC3194C05A15317Cb42E6E196) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		ControlsInfo: controlsInfo,
	}
}

func (w *SystemNetworks5D4297BaC3194C05A15317Cb42E6E196) GetControlsInfo() []control.ControlInfo {
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
	onceSystemNetworks5D4297BaC3194C05A15317Cb42E6E196     sync.Once
	instanceSystemNetworks5D4297BaC3194C05A15317Cb42E6E196 *SystemNetworks5D4297BaC3194C05A15317Cb42E6E196
)

func NewSystemNetworks5D4297BaC3194C05A15317Cb42E6E196(client *mqtt.Client) *SystemNetworks5D4297BaC3194C05A15317Cb42E6E196 {
	onceSystemNetworks5D4297BaC3194C05A15317Cb42E6E196.Do(func() {
		name := "system__networks__5d4297ba-c319-4c05-a153-17cb42e6e196"

		controlList := &SystemNetworks5D4297BaC3194C05A15317Cb42E6E196Controls{
			Name: control.NewTextControl(client, name, "Name", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Uuid: control.NewTextControl(client, name, "UUID", control.Meta{
				Type: "text",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Type: control.NewTextControl(client, name, "Type", control.Meta{
				Type: "text",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Active: control.NewSwitchControl(client, name, "Active", control.Meta{
				Type: "switch",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Device: control.NewTextControl(client, name, "Device", control.Meta{
				Type: "text",

				Order:    5,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			State: control.NewTextControl(client, name, "State", control.Meta{
				Type: "text",

				Order:    6,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Address: control.NewTextControl(client, name, "Address", control.Meta{
				Type: "text",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Connectivity: control.NewSwitchControl(client, name, "Connectivity", control.Meta{
				Type: "switch",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			UpDown: control.NewPushbuttonControl(client, name, "UpDown", control.Meta{
				Type: "pushbutton",

				Order:    12,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Up`},
			}),
			Operator: control.NewTextControl(client, name, "Operator", control.Meta{
				Type: "text",

				Order:    9,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			SignalQuality: control.NewTextControl(client, name, "SignalQuality", control.Meta{
				Type: "text",

				Order:    10,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Signal Quality`},
			}),
			AccessTechnologies: control.NewTextControl(client, name, "AccessTechnologies", control.Meta{
				Type: "text",

				Order:    11,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Access Technologies`},
			}),
		}

		instanceSystemNetworks5D4297BaC3194C05A15317Cb42E6E196 = &SystemNetworks5D4297BaC3194C05A15317Cb42E6E196{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworks5D4297BaC3194C05A15317Cb42E6E196
}
