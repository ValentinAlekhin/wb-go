package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type SystemControls struct {
	BatchNo            *controls.TextControl
	CurrentUptime      *controls.TextControl
	DtsVersion         *controls.TextControl
	HwRevision         *controls.TextControl
	ManufacturingDate  *controls.TextControl
	Reboot             *controls.PushbuttonControl
	ReleaseName        *controls.TextControl
	ReleaseSuite       *controls.TextControl
	ShortSn            *controls.TextControl
	TemperatureGrade   *controls.TextControl
	Status             *controls.TextControl
	ActivationLink     *controls.TextControl
	CloudBaseUrl       *controls.TextControl
	Name               *controls.TextControl
	Uuid               *controls.TextControl
	Type               *controls.TextControl
	Active             *controls.SwitchControl
	Device             *controls.TextControl
	State              *controls.TextControl
	Address            *controls.TextControl
	Connectivity       *controls.SwitchControl
	UpDown             *controls.PushbuttonControl
	Operator           *controls.TextControl
	SignalQuality      *controls.TextControl
	AccessTechnologies *controls.TextControl
}

type System struct {
	name     string
	device   string
	address  string
	Controls *SystemControls
}

func (w *System) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *System) GetControlsInfo() []controls.ControlInfo {
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
	onceSystem     sync.Once
	instanceSystem *System
)

func NewSystem(client *mqtt.Client) *System {
	onceSystem.Do(func() {
		device := "system"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &SystemControls{
			BatchNo: controls.NewTextControl(client, name, "Batch No", controls.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Batch No`, "ru": `Номер партии`},
			}),
			CurrentUptime: controls.NewTextControl(client, name, "Current uptime", controls.Meta{
				Type: "text",

				Order:    2,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Current uptime`, "ru": `Время работы`},
			}),
			DtsVersion: controls.NewTextControl(client, name, "DTS Version", controls.Meta{
				Type: "text",

				Order:    3,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `DTS Version`, "ru": `Версия DTS`},
			}),
			HwRevision: controls.NewTextControl(client, name, "HW Revision", controls.Meta{
				Type: "text",

				Order:    4,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `HW Revision`, "ru": `Версия контроллера`},
			}),
			ManufacturingDate: controls.NewTextControl(client, name, "Manufacturing Date", controls.Meta{
				Type: "text",

				Order:    5,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Manufacturing Date`, "ru": `Дата производства`},
			}),
			Reboot: controls.NewPushbuttonControl(client, name, "Reboot", controls.Meta{
				Type: "pushbutton",

				Order:    6,
				ReadOnly: false,
				Title:    controls.MultilingualText{"en": `Reboot`, "ru": `Перезагрузить`},
			}),
			ReleaseName: controls.NewTextControl(client, name, "Release name", controls.Meta{
				Type: "text",

				Order:    7,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Release name`, "ru": `Название релиза`},
			}),
			ReleaseSuite: controls.NewTextControl(client, name, "Release suite", controls.Meta{
				Type: "text",

				Order:    8,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Release suite`, "ru": `Тип релиза`},
			}),
			ShortSn: controls.NewTextControl(client, name, "Short SN", controls.Meta{
				Type: "text",

				Order:    9,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Short SN`, "ru": `Серийный номер`},
			}),
			TemperatureGrade: controls.NewTextControl(client, name, "Temperature Grade", controls.Meta{
				Type: "text",

				Order:    10,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Temperature Grade`, "ru": `Температурный диапазон`},
			}),
			Status: controls.NewTextControl(client, name, "status", controls.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Status`},
			}),
			ActivationLink: controls.NewTextControl(client, name, "activation_link", controls.Meta{
				Type: "text",

				Order:    2,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Link`},
			}),
			CloudBaseUrl: controls.NewTextControl(client, name, "cloud_base_url", controls.Meta{
				Type: "text",

				Order:    3,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `URL`},
			}),
			Name: controls.NewTextControl(client, name, "Name", controls.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			Uuid: controls.NewTextControl(client, name, "UUID", controls.Meta{
				Type: "text",

				Order:    2,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			Type: controls.NewTextControl(client, name, "Type", controls.Meta{
				Type: "text",

				Order:    3,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			Active: controls.NewSwitchControl(client, name, "Active", controls.Meta{
				Type: "switch",

				Order:    4,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			Device: controls.NewTextControl(client, name, "Device", controls.Meta{
				Type: "text",

				Order:    5,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			State: controls.NewTextControl(client, name, "State", controls.Meta{
				Type: "text",

				Order:    6,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			Address: controls.NewTextControl(client, name, "Address", controls.Meta{
				Type: "text",

				Order:    7,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			Connectivity: controls.NewSwitchControl(client, name, "Connectivity", controls.Meta{
				Type: "switch",

				Order:    8,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			UpDown: controls.NewPushbuttonControl(client, name, "UpDown", controls.Meta{
				Type: "pushbutton",

				Order:    12,
				ReadOnly: false,
				Title:    controls.MultilingualText{"en": `Down`},
			}),
			Operator: controls.NewTextControl(client, name, "Operator", controls.Meta{
				Type: "text",

				Order:    9,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			SignalQuality: controls.NewTextControl(client, name, "SignalQuality", controls.Meta{
				Type: "text",

				Order:    10,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Signal Quality`},
			}),
			AccessTechnologies: controls.NewTextControl(client, name, "AccessTechnologies", controls.Meta{
				Type: "text",

				Order:    11,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Access Technologies`},
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
