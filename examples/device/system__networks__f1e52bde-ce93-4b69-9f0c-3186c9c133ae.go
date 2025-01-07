package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type SystemNetworksF1E52BdeCe934B699F0C3186C9C133AeControls struct {
	Name         *control.TextControl
	Uuid         *control.TextControl
	Type         *control.TextControl
	Active       *control.SwitchControl
	Device       *control.TextControl
	State        *control.TextControl
	Address      *control.TextControl
	Connectivity *control.SwitchControl
	UpDown       *control.PushbuttonControl
}

type SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae struct {
	name     string
	Controls *SystemNetworksF1E52BdeCe934B699F0C3186C9C133AeControls
}

func (w *SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		ControlsInfo: controlsInfo,
	}
}

func (w *SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae) GetControlsInfo() []control.ControlInfo {
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
	onceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae     sync.Once
	instanceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae *SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae
)

func NewSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae(client *mqtt.Client) *SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae {
	onceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae.Do(func() {
		name := "system__networks__f1e52bde-ce93-4b69-9f0c-3186c9c133ae"

		controlList := &SystemNetworksF1E52BdeCe934B699F0C3186C9C133AeControls{
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
				Title:    control.MultilingualText{"en": `Down`},
			}),
		}

		instanceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae = &SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae
}
