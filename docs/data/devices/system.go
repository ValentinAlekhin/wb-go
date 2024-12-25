package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
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
	Name     string
	Address  string
	Controls *SystemControls
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
		deviceName := fmt.Sprintf("%s_%s", "system", "")
		controlList := &SystemControls{
			BatchNo:            controls.NewTextControl(client, deviceName, "Batch No"),
			CurrentUptime:      controls.NewTextControl(client, deviceName, "Current uptime"),
			DtsVersion:         controls.NewTextControl(client, deviceName, "DTS Version"),
			HwRevision:         controls.NewTextControl(client, deviceName, "HW Revision"),
			ManufacturingDate:  controls.NewTextControl(client, deviceName, "Manufacturing Date"),
			Reboot:             controls.NewPushbuttonControl(client, deviceName, "Reboot"),
			ReleaseName:        controls.NewTextControl(client, deviceName, "Release name"),
			ReleaseSuite:       controls.NewTextControl(client, deviceName, "Release suite"),
			ShortSn:            controls.NewTextControl(client, deviceName, "Short SN"),
			TemperatureGrade:   controls.NewTextControl(client, deviceName, "Temperature Grade"),
			Status:             controls.NewTextControl(client, deviceName, "status"),
			ActivationLink:     controls.NewTextControl(client, deviceName, "activation_link"),
			CloudBaseUrl:       controls.NewTextControl(client, deviceName, "cloud_base_url"),
			Name:               controls.NewTextControl(client, deviceName, "Name"),
			Uuid:               controls.NewTextControl(client, deviceName, "UUID"),
			Type:               controls.NewTextControl(client, deviceName, "Type"),
			Active:             controls.NewSwitchControl(client, deviceName, "Active"),
			Device:             controls.NewTextControl(client, deviceName, "Device"),
			State:              controls.NewTextControl(client, deviceName, "State"),
			Address:            controls.NewTextControl(client, deviceName, "Address"),
			Connectivity:       controls.NewSwitchControl(client, deviceName, "Connectivity"),
			UpDown:             controls.NewPushbuttonControl(client, deviceName, "UpDown"),
			Operator:           controls.NewTextControl(client, deviceName, "Operator"),
			SignalQuality:      controls.NewTextControl(client, deviceName, "SignalQuality"),
			AccessTechnologies: controls.NewTextControl(client, deviceName, "AccessTechnologies"),
		}

		instanceSystem = &System{
			Name:     deviceName,
			Address:  "",
			Controls: controlList,
		}
	})

	return instanceSystem
}
