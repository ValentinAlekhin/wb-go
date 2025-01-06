package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type Hwmoncontrols struct {
	BoardTemperature *control.ValueControl
	CpuTemperature   *control.ValueControl
}

type Hwmon struct {
	name     string
	device   string
	address  string
	Controls *Hwmoncontrols
}

func (w *Hwmon) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Hwmon) GetControlsInfo() []control.ControlInfo {
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
	onceHwmon     sync.Once
	instanceHwmon *Hwmon
)

func NewHwmon(client *mqtt.Client) *Hwmon {
	onceHwmon.Do(func() {
		device := "hwmon"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &Hwmoncontrols{
			BoardTemperature: control.NewValueControl(client, name, "Board Temperature", control.Meta{
				Type: "temperature",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			CpuTemperature: control.NewValueControl(client, name, "CPU Temperature", control.Meta{
				Type: "temperature",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
		}

		instanceHwmon = &Hwmon{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceHwmon
}
