package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type LightModeControls struct {
	Enabled *controls.SwitchControl
	State   *controls.ValueControl
}

type LightMode struct {
	Name     string
	Address  string
	Controls *LightModeControls
}

func (w *LightMode) GetControlsInfo() []controls.ControlInfo {
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
	onceLightMode     sync.Once
	instanceLightMode *LightMode
)

func NewLightMode(client *mqtt.Client) *LightMode {
	onceLightMode.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "light-mode", "")
		controlList := &LightModeControls{
			Enabled: controls.NewSwitchControl(client, deviceName, "enabled"),
			State:   controls.NewValueControl(client, deviceName, "state"),
		}

		instanceLightMode = &LightMode{
			Name:     deviceName,
			Address:  "",
			Controls: controlList,
		}
	})

	return instanceLightMode
}
