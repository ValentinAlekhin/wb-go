package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type TermostatControls struct {
	R01Ts161Lock     *controls.SwitchControl
	R01Ts161Mode     *controls.SwitchControl
	R01Ts161Setpoint *controls.RangeControl
}

type Termostat struct {
	name     string
	device   string
	address  string
	Controls *TermostatControls
}

func (w *Termostat) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Termostat) GetControlsInfo() []controls.ControlInfo {
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
	onceTermostat     sync.Once
	instanceTermostat *Termostat
)

func NewTermostat(client *mqtt.Client) *Termostat {
	onceTermostat.Do(func() {
		device := "Termostat"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &TermostatControls{
			R01Ts161Lock: controls.NewSwitchControl(client, name, "R01-TS16-1-lock", controls.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			R01Ts161Mode: controls.NewSwitchControl(client, name, "R01-TS16-1-mode", controls.Meta{
				Type: "switch",

				Order:    2,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			R01Ts161Setpoint: controls.NewRangeControl(client, name, "R01-TS16-1-setpoint", controls.Meta{
				Type: "range",

				Max: 30,

				Order:    3,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
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
