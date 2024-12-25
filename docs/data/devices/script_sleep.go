package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type ScriptsleepControls struct {
	Current          *controls.ValueControl
	Enable           *controls.SwitchControl
	State            *controls.TextControl
	Target           *controls.RangeControl
	Zone1RelayStatus *controls.SwitchControl
	Zone1Status      *controls.ValueControl
}

type Scriptsleep struct {
	Name     string
	Address  string
	Controls *ScriptsleepControls
}

func (w *Scriptsleep) GetControlsInfo() []controls.ControlInfo {
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
	onceScriptsleep     sync.Once
	instanceScriptsleep *Scriptsleep
)

func NewScriptsleep(client *mqtt.Client) *Scriptsleep {
	onceScriptsleep.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "script", "sleep")
		controlList := &ScriptsleepControls{
			Current:          controls.NewValueControl(client, deviceName, "current"),
			Enable:           controls.NewSwitchControl(client, deviceName, "enable"),
			State:            controls.NewTextControl(client, deviceName, "state"),
			Target:           controls.NewRangeControl(client, deviceName, "target"),
			Zone1RelayStatus: controls.NewSwitchControl(client, deviceName, "zone1_relay_status"),
			Zone1Status:      controls.NewValueControl(client, deviceName, "zone1_status"),
		}

		instanceScriptsleep = &Scriptsleep{
			Name:     deviceName,
			Address:  "sleep",
			Controls: controlList,
		}
	})

	return instanceScriptsleep
}
