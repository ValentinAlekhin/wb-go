package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type Powerstatuscontrols struct {
	Vin              *control.ValueControl
	WorkingOnBattery *control.SwitchControl
}

type Powerstatus struct {
	name     string
	device   string
	address  string
	Controls *Powerstatuscontrols
}

func (w *Powerstatus) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Powerstatus) GetControlsInfo() []control.ControlInfo {
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
	oncePowerstatus     sync.Once
	instancePowerstatus *Powerstatus
)

func NewPowerstatus(client *mqtt.Client) *Powerstatus {
	oncePowerstatus.Do(func() {
		device := "power"
		address := "status"
		name := fmt.Sprintf("%s_%s", device, address)

		controlList := &Powerstatuscontrols{
			Vin: control.NewValueControl(client, name, "Vin", control.Meta{
				Type: "voltage",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Input voltage`, "ru": `Входное напряжение`},
			}),
			WorkingOnBattery: control.NewSwitchControl(client, name, "working on battery", control.Meta{
				Type: "switch",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Working on battery`, "ru": `Работа от батареи`},
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
