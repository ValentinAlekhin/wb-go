package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type AdaptiveLightTestControls struct {
	Enabled        *control.SwitchControl
	MinTemperature *control.RangeControl
	MaxTemperature *control.RangeControl
	Temperature    *control.RangeControl
	MinBrightness  *control.RangeControl
	MaxBrightness  *control.RangeControl
	Brightness     *control.RangeControl
	SleepMode      *control.SwitchControl
	Sunrise        *control.TextControl
	Sunset         *control.TextControl
	SlipStart      *control.TextControl
	SlipEnd        *control.TextControl
}

type AdaptiveLightTest struct {
	name     string
	Controls *AdaptiveLightTestControls
}

func (w *AdaptiveLightTest) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		ControlsInfo: controlsInfo,
	}
}

func (w *AdaptiveLightTest) GetControlsInfo() []control.ControlInfo {
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
	onceAdaptiveLightTest     sync.Once
	instanceAdaptiveLightTest *AdaptiveLightTest
)

func NewAdaptiveLightTest(client *mqtt.Client) *AdaptiveLightTest {
	onceAdaptiveLightTest.Do(func() {
		name := "adaptive-light-test"

		controlList := &AdaptiveLightTestControls{
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
			Temperature: control.NewRangeControl(client, name, "Temperature", control.Meta{
				Type: "range",

				Max: 100,

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Температура`},
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
			Brightness: control.NewRangeControl(client, name, "Brightness", control.Meta{
				Type: "range",

				Max: 100,

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Яркость`},
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
		}

		instanceAdaptiveLightTest = &AdaptiveLightTest{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceAdaptiveLightTest
}
