package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type ScriptSleepControls struct {
	Current          *control.ValueControl
	Enable           *control.SwitchControl
	State            *control.TextControl
	Target           *control.RangeControl
	Zone1RelayStatus *control.SwitchControl
	Zone1Status      *control.ValueControl
}

type ScriptSleep struct {
	name     string
	Controls *ScriptSleepControls
}

func (w *ScriptSleep) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		ControlsInfo: controlsInfo,
	}
}

func (w *ScriptSleep) GetControlsInfo() []control.ControlInfo {
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
	onceScriptSleep     sync.Once
	instanceScriptSleep *ScriptSleep
)

func NewScriptSleep(client *mqtt.Client) *ScriptSleep {
	onceScriptSleep.Do(func() {
		name := "script_sleep"

		controlList := &ScriptSleepControls{
			Current: control.NewValueControl(client, name, "current", control.Meta{
				Type: "temperature",

				Order:    16,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Current Temperature`, "ru": `Current Temperature`},
			}),
			Enable: control.NewSwitchControl(client, name, "enable", control.Meta{
				Type: "switch",

				Order:    10,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Enable`, "ru": `Enable`},
			}),
			State: control.NewTextControl(client, name, "state", control.Meta{
				Type: "text",

				Order:    15,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Mode State`, "ru": `State`},
			}),
			Target: control.NewRangeControl(client, name, "target", control.Meta{
				Type: "range",

				Max: 30,
				Min: 14,

				Order:    20,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Target Temperature`, "ru": `Target Temperature`},
			}),
			Zone1RelayStatus: control.NewSwitchControl(client, name, "zone1_relay_status", control.Meta{
				Type: "switch",

				Order:    50,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Zone Relay Status`, "ru": `Zone Relay Status`},
			}),
			Zone1Status: control.NewValueControl(client, name, "zone1_status", control.Meta{
				Type: "temperature",

				Order:    30,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Zone Temperature`, "ru": `Zone Temperature`},
			}),
		}

		instanceScriptSleep = &ScriptSleep{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceScriptSleep
}
