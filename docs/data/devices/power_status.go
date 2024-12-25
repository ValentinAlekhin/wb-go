package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type PowerstatusControls struct {
	Vin              *controls.ValueControl
	WorkingOnBattery *controls.SwitchControl
}

type Powerstatus struct {
	Name     string
	Address  string
	Controls *PowerstatusControls
}

func (w *Powerstatus) GetControlsInfo() []controls.ControlInfo {
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
	oncePowerstatus     sync.Once
	instancePowerstatus *Powerstatus
)

func NewPowerstatus(client *mqtt.Client) *Powerstatus {
	oncePowerstatus.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "power", "status")
		controlList := &PowerstatusControls{
			Vin:              controls.NewValueControl(client, deviceName, "Vin"),
			WorkingOnBattery: controls.NewSwitchControl(client, deviceName, "working on battery"),
		}

		instancePowerstatus = &Powerstatus{
			Name:     deviceName,
			Address:  "status",
			Controls: controlList,
		}
	})

	return instancePowerstatus
}
