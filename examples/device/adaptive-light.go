package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type AdaptiveLightcontrols struct {
	Enabled        *control.SwitchControl
	MinTemperature *control.RangeControl
	MaxTemperature *control.RangeControl
	MinBrightness  *control.RangeControl
	MaxBrightness  *control.RangeControl
	SleepMode      *control.SwitchControl
	Sunrise        *control.TextControl
	Sunset         *control.TextControl
	SlipStart      *control.TextControl
	SlipEnd        *control.TextControl
	Temperature    *control.RangeControl
	Brightness     *control.RangeControl
}

type AdaptiveLight struct {
	name     string
	device   string
	address  string
	Controls *AdaptiveLightcontrols
}

func (w *AdaptiveLight) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *AdaptiveLight) GetControlsInfo() []control.ControlInfo {
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
	onceAdaptiveLight     sync.Once
	instanceAdaptiveLight *AdaptiveLight
)

func NewAdaptiveLight(client *mqtt.Client) *AdaptiveLight {
	onceAdaptiveLight.Do(func() {
		device := "adaptive-light"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &AdaptiveLightcontrols{
			Enabled: control.NewSwitchControl(client, name, "Enabled", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Включено`},
			}),
			MinTemperature: control.NewRangeControl(client, name, "Min Temperature", control.Meta{
				Type: "range",

				Max: 100,

				Order:    2,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Минимальная температура`},
			}),
			MaxTemperature: control.NewRangeControl(client, name, "Max Temperature", control.Meta{
				Type: "range",

				Max: 100,

				Order:    3,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Максимальная температура`},
			}),
			MinBrightness: control.NewRangeControl(client, name, "Min Brightness", control.Meta{
				Type: "range",

				Max: 100,

				Order:    5,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Минимальная яркость`},
			}),
			MaxBrightness: control.NewRangeControl(client, name, "Max Brightness", control.Meta{
				Type: "range",

				Max: 100,

				Order:    6,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Максимальная яркость`},
			}),
			SleepMode: control.NewSwitchControl(client, name, "Sleep Mode", control.Meta{
				Type: "switch",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Режим сна`},
			}),
			Sunrise: control.NewTextControl(client, name, "Sunrise", control.Meta{
				Type: "text",

				Order:    9,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Рассвет`},
			}),
			Sunset: control.NewTextControl(client, name, "Sunset", control.Meta{
				Type: "text",

				Order:    10,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Закат`},
			}),
			SlipStart: control.NewTextControl(client, name, "Slip Start", control.Meta{
				Type: "text",

				Order:    11,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Начало сна`},
			}),
			SlipEnd: control.NewTextControl(client, name, "Slip End", control.Meta{
				Type: "text",

				Order:    12,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Конец сна`},
			}),
			Temperature: control.NewRangeControl(client, name, "Temperature", control.Meta{
				Type: "range",

				Max: 100,

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Температура`},
			}),
			Brightness: control.NewRangeControl(client, name, "Brightness", control.Meta{
				Type: "range",

				Max: 100,

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Яркость`},
			}),
		}

		instanceAdaptiveLight = &AdaptiveLight{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceAdaptiveLight
}
