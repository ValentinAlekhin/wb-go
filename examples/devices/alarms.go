package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type AlarmsControls struct {
	Log *controls.TextControl
}

type Alarms struct {
	name     string
	device   string
	address  string
	Controls *AlarmsControls
}

func (w *Alarms) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Alarms) GetControlsInfo() []controls.ControlInfo {
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
	onceAlarms     sync.Once
	instanceAlarms *Alarms
)

func NewAlarms(client *mqtt.Client) *Alarms {
	onceAlarms.Do(func() {
		device := "alarms"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &AlarmsControls{
			Log: controls.NewTextControl(client, name, "log", controls.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Log`, "ru": `Лог`},
			}),
		}

		instanceAlarms = &Alarms{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceAlarms
}
