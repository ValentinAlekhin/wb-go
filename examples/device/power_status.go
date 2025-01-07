package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type PowerStatusControls struct {
	Vin              *control.ValueControl
	WorkingOnBattery *control.SwitchControl
}

type PowerStatus struct {
	name     string
	Controls *PowerStatusControls
}

func (w *PowerStatus) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		ControlsInfo: controlsInfo,
	}
}

func (w *PowerStatus) GetControlsInfo() []control.ControlInfo {
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
	oncePowerStatus     sync.Once
	instancePowerStatus *PowerStatus
)

func NewPowerStatus(client *mqtt.Client) *PowerStatus {
	oncePowerStatus.Do(func() {
		name := "power_status"

		controlList := &PowerStatusControls{
			Vin: control.NewValueControl(client, name, "Vin", control.Meta{
				Type: "voltage",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Input voltage`, "ru": `Входное напряжение`},
			}),
			WorkingOnBattery: control.NewSwitchControl(client, name, "working on battery", control.Meta{
				Type: "switch",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Working on battery`, "ru": `Работа от батареи`},
			}),
		}

		instancePowerStatus = &PowerStatus{
			name:     name,
			Controls: controlList,
		}
	})

	return instancePowerStatus
}
