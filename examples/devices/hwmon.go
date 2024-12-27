package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type HwmonControls struct {
	BoardTemperature *controls.ValueControl
	CpuTemperature   *controls.ValueControl
}

type Hwmon struct {
	name     string
	device   string
	address  string
	Controls *HwmonControls
}

func (w *Hwmon) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Hwmon) GetControlsInfo() []controls.ControlInfo {
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
	onceHwmon     sync.Once
	instanceHwmon *Hwmon
)

func NewHwmon(client *mqtt.Client) *Hwmon {
	onceHwmon.Do(func() {
		device := "hwmon"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &HwmonControls{
			BoardTemperature: controls.NewValueControl(client, name, "Board Temperature", controls.Meta{
				Type: "temperature",

				Order:    1,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			CpuTemperature: controls.NewValueControl(client, name, "CPU Temperature", controls.Meta{
				Type: "temperature",

				Order:    2,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
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
