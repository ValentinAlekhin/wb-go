package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type SystemControls struct {
	BatchNo           *control.TextControl
	CurrentUptime     *control.TextControl
	DtsVersion        *control.TextControl
	HwRevision        *control.TextControl
	ManufacturingDate *control.TextControl
	Reboot            *control.PushbuttonControl
	ReleaseName       *control.TextControl
	ReleaseSuite      *control.TextControl
	ShortSn           *control.TextControl
	TemperatureGrade  *control.TextControl
}

type System struct {
	name     string
	Controls *SystemControls
}

func (w *System) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		ControlsInfo: controlsInfo,
	}
}

func (w *System) GetControlsInfo() []control.ControlInfo {
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
	onceSystem     sync.Once
	instanceSystem *System
)

func NewSystem(client *mqtt.Client) *System {
	onceSystem.Do(func() {
		name := "system"

		controlList := &SystemControls{
			BatchNo: control.NewTextControl(client, name, "Batch No", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Batch No`, "ru": `Номер партии`},
			}),
			CurrentUptime: control.NewTextControl(client, name, "Current uptime", control.Meta{
				Type: "text",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Current uptime`, "ru": `Время работы`},
			}),
			DtsVersion: control.NewTextControl(client, name, "DTS Version", control.Meta{
				Type: "text",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `DTS Version`, "ru": `Версия DTS`},
			}),
			HwRevision: control.NewTextControl(client, name, "HW Revision", control.Meta{
				Type: "text",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `HW Revision`, "ru": `Версия контроллера`},
			}),
			ManufacturingDate: control.NewTextControl(client, name, "Manufacturing Date", control.Meta{
				Type: "text",

				Order:    5,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Manufacturing Date`, "ru": `Дата производства`},
			}),
			Reboot: control.NewPushbuttonControl(client, name, "Reboot", control.Meta{
				Type: "pushbutton",

				Order:    6,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Reboot`, "ru": `Перезагрузить`},
			}),
			ReleaseName: control.NewTextControl(client, name, "Release name", control.Meta{
				Type: "text",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Release name`, "ru": `Название релиза`},
			}),
			ReleaseSuite: control.NewTextControl(client, name, "Release suite", control.Meta{
				Type: "text",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Release suite`, "ru": `Тип релиза`},
			}),
			ShortSn: control.NewTextControl(client, name, "Short SN", control.Meta{
				Type: "text",

				Order:    9,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Short SN`, "ru": `Серийный номер`},
			}),
			TemperatureGrade: control.NewTextControl(client, name, "Temperature Grade", control.Meta{
				Type: "text",

				Order:    10,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Temperature Grade`, "ru": `Температурный диапазон`},
			}),
		}

		instanceSystem = &System{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystem
}
