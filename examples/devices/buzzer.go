package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type BuzzerControls struct {
	Enabled   *controls.SwitchControl
	Frequency *controls.RangeControl
	Volume    *controls.RangeControl
}

type Buzzer struct {
	name     string
	device   string
	address  string
	Controls *BuzzerControls
}

func (w *Buzzer) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Buzzer) GetControlsInfo() []controls.ControlInfo {
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
	onceBuzzer     sync.Once
	instanceBuzzer *Buzzer
)

func NewBuzzer(client *mqtt.Client) *Buzzer {
	onceBuzzer.Do(func() {
		device := "buzzer"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &BuzzerControls{
			Enabled: controls.NewSwitchControl(client, name, "enabled", controls.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: false,
				Title:    controls.MultilingualText{"en": `Enabled`, "ru": `Включен`},
			}),
			Frequency: controls.NewRangeControl(client, name, "frequency", controls.Meta{
				Type: "range",

				Max: 7000,

				Order:    2,
				ReadOnly: false,
				Title:    controls.MultilingualText{"en": `Frequency`, "ru": `Частота`},
			}),
			Volume: controls.NewRangeControl(client, name, "volume", controls.Meta{
				Type: "range",

				Max: 100,

				Order:    3,
				ReadOnly: false,
				Title:    controls.MultilingualText{"en": `Volume`, "ru": `Громкость`},
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
