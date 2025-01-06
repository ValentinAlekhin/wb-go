package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type Systemcontrols struct {
	BatchNo            *control.TextControl
	CurrentUptime      *control.TextControl
	DtsVersion         *control.TextControl
	HwRevision         *control.TextControl
	ManufacturingDate  *control.TextControl
	Reboot             *control.PushbuttonControl
	ReleaseName        *control.TextControl
	ReleaseSuite       *control.TextControl
	ShortSn            *control.TextControl
	TemperatureGrade   *control.TextControl
	Status             *control.TextControl
	ActivationLink     *control.TextControl
	CloudBaseUrl       *control.TextControl
	Name               *control.TextControl
	Uuid               *control.TextControl
	Type               *control.TextControl
	Active             *control.SwitchControl
	Device             *control.TextControl
	State              *control.TextControl
	Address            *control.TextControl
	Connectivity       *control.SwitchControl
	UpDown             *control.PushbuttonControl
	Operator           *control.TextControl
	SignalQuality      *control.TextControl
	AccessTechnologies *control.TextControl
}

type System struct {
	name     string
	device   string
	address  string
	Controls *Systemcontrols
}

func (w *System) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
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
		device := "system"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &Systemcontrols{
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
			Status: control.NewTextControl(client, name, "status", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Status`},
			}),
			ActivationLink: control.NewTextControl(client, name, "activation_link", control.Meta{
				Type: "text",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Link`},
			}),
			CloudBaseUrl: control.NewTextControl(client, name, "cloud_base_url", control.Meta{
				Type: "text",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `URL`},
			}),
			Name: control.NewTextControl(client, name, "Name", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Uuid: control.NewTextControl(client, name, "UUID", control.Meta{
				Type: "text",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Type: control.NewTextControl(client, name, "Type", control.Meta{
				Type: "text",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Active: control.NewSwitchControl(client, name, "Active", control.Meta{
				Type: "switch",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Device: control.NewTextControl(client, name, "Device", control.Meta{
				Type: "text",

				Order:    5,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			State: control.NewTextControl(client, name, "State", control.Meta{
				Type: "text",

				Order:    6,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Address: control.NewTextControl(client, name, "Address", control.Meta{
				Type: "text",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Connectivity: control.NewSwitchControl(client, name, "Connectivity", control.Meta{
				Type: "switch",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			UpDown: control.NewPushbuttonControl(client, name, "UpDown", control.Meta{
				Type: "pushbutton",

				Order:    12,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Down`},
			}),
			Operator: control.NewTextControl(client, name, "Operator", control.Meta{
				Type: "text",

				Order:    9,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			SignalQuality: control.NewTextControl(client, name, "SignalQuality", control.Meta{
				Type: "text",

				Order:    10,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Signal Quality`},
			}),
			AccessTechnologies: control.NewTextControl(client, name, "AccessTechnologies", control.Meta{
				Type: "text",

				Order:    11,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Access Technologies`},
			}),
		}

		instanceSystem = &System{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceSystem
}
