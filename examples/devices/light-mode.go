package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type LightModeControls struct {
	Enabled *controls.SwitchControl
	State   *controls.ValueControl
}

type LightMode struct {
	name     string
	device   string
	address  string
	Controls *LightModeControls
}

func (w *LightMode) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
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
		device := "light-mode"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &LightModeControls{
			Enabled: controls.NewSwitchControl(client, name, "enabled", controls.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			State: controls.NewValueControl(client, name, "state", controls.Meta{
				Type: "value",

				Order:    2,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `State`},
			}),
		}

		instanceLightMode = &LightMode{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceLightMode
}
