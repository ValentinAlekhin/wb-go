package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
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
	name     string
	device   string
	address  string
	Controls *ScriptsleepControls
}

func (w *Scriptsleep) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
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
		device := "script"
		address := "sleep"
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &ScriptsleepControls{
			Current: controls.NewValueControl(client, name, "current", controls.Meta{
				Type: "temperature",

				Order:    16,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Current Temperature`, "ru": `Current Temperature`},
			}),
			Enable: controls.NewSwitchControl(client, name, "enable", controls.Meta{
				Type: "switch",

				Order:    10,
				ReadOnly: false,
				Title:    controls.MultilingualText{"en": `Enable`, "ru": `Enable`},
			}),
			State: controls.NewTextControl(client, name, "state", controls.Meta{
				Type: "text",

				Order:    15,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Mode State`, "ru": `State`},
			}),
			Target: controls.NewRangeControl(client, name, "target", controls.Meta{
				Type: "range",

				Max: 30,
				Min: 14,

				Order:    20,
				ReadOnly: false,
				Title:    controls.MultilingualText{"en": `Target Temperature`, "ru": `Target Temperature`},
			}),
			Zone1RelayStatus: controls.NewSwitchControl(client, name, "zone1_relay_status", controls.Meta{
				Type: "switch",

				Order:    50,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Zone Relay Status`, "ru": `Zone Relay Status`},
			}),
			Zone1Status: controls.NewValueControl(client, name, "zone1_status", controls.Meta{
				Type: "temperature",

				Order:    30,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Zone Temperature`, "ru": `Zone Temperature`},
			}),
		}

		instanceScriptsleep = &Scriptsleep{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceScriptsleep
}
