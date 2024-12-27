package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type PowerstatusControls struct {
	Vin              *controls.ValueControl
	WorkingOnBattery *controls.SwitchControl
}

type Powerstatus struct {
	name     string
	device   string
	address  string
	Controls *PowerstatusControls
}

func (w *Powerstatus) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Powerstatus) GetControlsInfo() []controls.ControlInfo {
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
	oncePowerstatus     sync.Once
	instancePowerstatus *Powerstatus
)

func NewPowerstatus(client *mqtt.Client) *Powerstatus {
	oncePowerstatus.Do(func() {
		device := "power"
		address := "status"
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &PowerstatusControls{
			Vin: controls.NewValueControl(client, name, "Vin", controls.Meta{
				Type: "voltage",

				Order:    1,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Input voltage`, "ru": `Входное напряжение`},
			}),
			WorkingOnBattery: controls.NewSwitchControl(client, name, "working on battery", controls.Meta{
				Type: "switch",

				Order:    2,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Working on battery`, "ru": `Работа от батареи`},
			}),
		}

		instancePowerstatus = &Powerstatus{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instancePowerstatus
}
