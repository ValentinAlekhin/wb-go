package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type Thermostatcontrols struct {
	TargetTemperature  *control.RangeControl
	CurrentTemperature *control.ValueControl
	Enabled            *control.SwitchControl
	On                 *control.SwitchControl
}

type Thermostat struct {
	name     string
	device   string
	address  string
	Controls *Thermostatcontrols
}

func (w *Thermostat) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Thermostat) GetControlsInfo() []control.ControlInfo {
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
	onceThermostat     sync.Once
	instanceThermostat *Thermostat
)

func NewThermostat(client *mqtt.Client) *Thermostat {
	onceThermostat.Do(func() {
		device := "thermostat"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &Thermostatcontrols{
			TargetTemperature: control.NewRangeControl(client, name, "Target Temperature", control.Meta{
				Type:  "range",
				Units: "°C",
				Max:   100,

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Целевая температура`},
			}),
			CurrentTemperature: control.NewValueControl(client, name, "Current Temperature", control.Meta{
				Type:  "value",
				Units: "°C",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Текущая температура`},
			}),
			Enabled: control.NewSwitchControl(client, name, "Enabled", control.Meta{
				Type: "switch",

				Order:    3,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Термостат включен`},
			}),
			On: control.NewSwitchControl(client, name, "On", control.Meta{
				Type: "switch",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Нагрев`},
			}),
		}

		instanceThermostat = &Thermostat{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceThermostat
}
