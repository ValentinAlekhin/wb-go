package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type Termostatcontrols struct {
	R01Ts161Lock     *control.SwitchControl
	R01Ts161Mode     *control.SwitchControl
	R01Ts161Setpoint *control.RangeControl
}

type Termostat struct {
	name     string
	device   string
	address  string
	Controls *Termostatcontrols
}

func (w *Termostat) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Termostat) GetControlsInfo() []control.ControlInfo {
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
	onceTermostat     sync.Once
	instanceTermostat *Termostat
)

func NewTermostat(client *mqtt.Client) *Termostat {
	onceTermostat.Do(func() {
		device := "Termostat"
		address := ""
		name := device

		controlList := &Termostatcontrols{
			R01Ts161Lock: control.NewSwitchControl(client, name, "R01-TS16-1-lock", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			R01Ts161Mode: control.NewSwitchControl(client, name, "R01-TS16-1-mode", control.Meta{
				Type: "switch",

				Order:    2,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			R01Ts161Setpoint: control.NewRangeControl(client, name, "R01-TS16-1-setpoint", control.Meta{
				Type: "range",

				Max: 30,

				Order:    3,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
		}

		instanceTermostat = &Termostat{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceTermostat
}
