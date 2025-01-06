package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type Buzzercontrols struct {
	Enabled   *control.SwitchControl
	Frequency *control.RangeControl
	Volume    *control.RangeControl
}

type Buzzer struct {
	name     string
	device   string
	address  string
	Controls *Buzzercontrols
}

func (w *Buzzer) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Buzzer) GetControlsInfo() []control.ControlInfo {
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
	onceBuzzer     sync.Once
	instanceBuzzer *Buzzer
)

func NewBuzzer(client *mqtt.Client) *Buzzer {
	onceBuzzer.Do(func() {
		device := "buzzer"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &Buzzercontrols{
			Enabled: control.NewSwitchControl(client, name, "enabled", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Enabled`, "ru": `Включен`},
			}),
			Frequency: control.NewRangeControl(client, name, "frequency", control.Meta{
				Type: "range",

				Max: 7000,

				Order:    2,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Frequency`, "ru": `Частота`},
			}),
			Volume: control.NewRangeControl(client, name, "volume", control.Meta{
				Type: "range",

				Max: 100,

				Order:    3,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Volume`, "ru": `Громкость`},
			}),
		}

		instanceBuzzer = &Buzzer{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceBuzzer
}
