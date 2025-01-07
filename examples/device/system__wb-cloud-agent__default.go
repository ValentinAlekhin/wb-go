package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type SystemWbCloudAgentDefaultControls struct {
	Status         *control.TextControl
	ActivationLink *control.TextControl
	CloudBaseUrl   *control.TextControl
}

type SystemWbCloudAgentDefault struct {
	name     string
	Controls *SystemWbCloudAgentDefaultControls
}

func (w *SystemWbCloudAgentDefault) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		ControlsInfo: controlsInfo,
	}
}

func (w *SystemWbCloudAgentDefault) GetControlsInfo() []control.ControlInfo {
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
	onceSystemWbCloudAgentDefault     sync.Once
	instanceSystemWbCloudAgentDefault *SystemWbCloudAgentDefault
)

func NewSystemWbCloudAgentDefault(client *mqtt.Client) *SystemWbCloudAgentDefault {
	onceSystemWbCloudAgentDefault.Do(func() {
		name := "system__wb-cloud-agent__default"

		controlList := &SystemWbCloudAgentDefaultControls{
			Status: control.NewTextControl(client, name, "status", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Status`},
			}),
			ActivationLink: control.NewTextControl(client, name, "activation_link", control.Meta{
				Type: "text",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Link`},
			}),
			CloudBaseUrl: control.NewTextControl(client, name, "cloud_base_url", control.Meta{
				Type: "text",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `URL`},
			}),
		}

		instanceSystemWbCloudAgentDefault = &SystemWbCloudAgentDefault{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemWbCloudAgentDefault
}
