package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbrulesControls struct {
	RuleDebugging *controls.SwitchControl
}

type Wbrules struct {
	name     string
	device   string
	address  string
	Controls *WbrulesControls
}

func (w *Wbrules) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Wbrules) GetControlsInfo() []controls.ControlInfo {
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
	onceWbrules     sync.Once
	instanceWbrules *Wbrules
)

func NewWbrules(client *mqtt.Client) *Wbrules {
	onceWbrules.Do(func() {
		device := "wbrules"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &WbrulesControls{
			RuleDebugging: controls.NewSwitchControl(client, name, "Rule debugging", controls.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    controls.MultilingualText{"en": `Rule debugging`, "ru": `Отладка правил`},
			}),
		}

		instanceWbrules = &Wbrules{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceWbrules
}
