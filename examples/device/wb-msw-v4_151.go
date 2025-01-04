package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbMswV4151Controls struct {
	Temperature       *control.ValueControl
	Humidity          *control.ValueControl
	Co2               *control.ValueControl
	AirQualityVoc     *control.ValueControl
	SoundLevel        *control.ValueControl
	Illuminance       *control.ValueControl
	MaxMotion         *control.ValueControl
	CurrentMotion     *control.ValueControl
	Buzzer            *control.SwitchControl
	RedLed            *control.SwitchControl
	GreenLed          *control.SwitchControl
	LedPeriods        *control.RangeControl
	LedGlowDurationms *control.RangeControl
	LearnToRam        *control.SwitchControl
	PlayFromRam       *control.PushbuttonControl
	LearnToRom1       *control.SwitchControl
	LearnToRom2       *control.SwitchControl
	LearnToRom3       *control.SwitchControl
	LearnToRom4       *control.SwitchControl
	LearnToRom5       *control.SwitchControl
	LearnToRom6       *control.SwitchControl
	LearnToRom7       *control.SwitchControl
	PlayFromRom1      *control.PushbuttonControl
	PlayFromRom2      *control.PushbuttonControl
	PlayFromRom3      *control.PushbuttonControl
	PlayFromRom4      *control.PushbuttonControl
	PlayFromRom5      *control.PushbuttonControl
	PlayFromRom6      *control.PushbuttonControl
	PlayFromRom7      *control.PushbuttonControl
	Serial            *control.TextControl
}

type WbMswV4151 struct {
	name     string
	device   string
	address  string
	Controls *WbMswV4151Controls
}

func (w *WbMswV4151) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbMswV4151) GetControlsInfo() []control.ControlInfo {
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
	onceWbMswV4151     sync.Once
	instanceWbMswV4151 *WbMswV4151
)

func NewWbMswV4151(client *mqtt.Client) *WbMswV4151 {
	onceWbMswV4151.Do(func() {
		device := "wb-msw-v4"
		address := "151"
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &WbMswV4151Controls{
			Temperature: control.NewValueControl(client, name, "Temperature", control.Meta{
				Type:  "value",
				Units: "deg C",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Температура`},
			}),
			Humidity: control.NewValueControl(client, name, "Humidity", control.Meta{
				Type:  "value",
				Units: "%, RH",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Влажность`},
			}),
			Co2: control.NewValueControl(client, name, "CO2", control.Meta{
				Type: "concentration",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `CO₂`, "ru": `Уровень CO₂`},
			}),
			AirQualityVoc: control.NewValueControl(client, name, "Air Quality (VOC)", control.Meta{
				Type:  "value",
				Units: "ppb",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Качество воздуха (VOC)`},
			}),
			SoundLevel: control.NewValueControl(client, name, "Sound Level", control.Meta{
				Type: "sound_level",

				Order:    5,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Уровень шума`},
			}),
			Illuminance: control.NewValueControl(client, name, "Illuminance", control.Meta{
				Type: "lux",

				Order:    6,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Освещенность`},
			}),
			MaxMotion: control.NewValueControl(client, name, "Max Motion", control.Meta{
				Type: "value",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Максимальное движение`},
			}),
			CurrentMotion: control.NewValueControl(client, name, "Current Motion", control.Meta{
				Type: "value",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Текущее движение`},
			}),
			Buzzer: control.NewSwitchControl(client, name, "Buzzer", control.Meta{
				Type: "switch",

				Order:    9,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Зуммер`},
			}),
			RedLed: control.NewSwitchControl(client, name, "Red LED", control.Meta{
				Type: "switch",

				Order:    10,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Красный светодиод`},
			}),
			GreenLed: control.NewSwitchControl(client, name, "Green LED", control.Meta{
				Type: "switch",

				Order:    11,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Зеленый светодиод`},
			}),
			LedPeriods: control.NewRangeControl(client, name, "LED Period (s)", control.Meta{
				Type: "range",

				Max: 10,

				Order:    12,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Период включения светодиодов (с)`},
			}),
			LedGlowDurationms: control.NewRangeControl(client, name, "LED Glow Duration (ms)", control.Meta{
				Type: "range",

				Max: 50,

				Order:    13,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Длительность включения светодиодов (мс)`},
			}),
			LearnToRam: control.NewSwitchControl(client, name, "Learn to RAM", control.Meta{
				Type: "switch",

				Order:    14,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Записать команду в RAM`},
			}),
			PlayFromRam: control.NewPushbuttonControl(client, name, "Play from RAM", control.Meta{
				Type: "pushbutton",

				Order:    15,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Воспроизвести команду из RAM`},
			}),
			LearnToRom1: control.NewSwitchControl(client, name, "Learn to ROM1", control.Meta{
				Type: "switch",

				Order:    16,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Записать команду в ROM1`},
			}),
			LearnToRom2: control.NewSwitchControl(client, name, "Learn to ROM2", control.Meta{
				Type: "switch",

				Order:    17,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Записать команду в ROM2`},
			}),
			LearnToRom3: control.NewSwitchControl(client, name, "Learn to ROM3", control.Meta{
				Type: "switch",

				Order:    18,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Записать команду в ROM3`},
			}),
			LearnToRom4: control.NewSwitchControl(client, name, "Learn to ROM4", control.Meta{
				Type: "switch",

				Order:    19,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Записать команду в ROM4`},
			}),
			LearnToRom5: control.NewSwitchControl(client, name, "Learn to ROM5", control.Meta{
				Type: "switch",

				Order:    20,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Записать команду в ROM5`},
			}),
			LearnToRom6: control.NewSwitchControl(client, name, "Learn to ROM6", control.Meta{
				Type: "switch",

				Order:    21,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Записать команду в ROM6`},
			}),
			LearnToRom7: control.NewSwitchControl(client, name, "Learn to ROM7", control.Meta{
				Type: "switch",

				Order:    22,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Записать команду в ROM7`},
			}),
			PlayFromRom1: control.NewPushbuttonControl(client, name, "Play from ROM1", control.Meta{
				Type: "pushbutton",

				Order:    23,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Воспроизвести команду из ROM1`},
			}),
			PlayFromRom2: control.NewPushbuttonControl(client, name, "Play from ROM2", control.Meta{
				Type: "pushbutton",

				Order:    24,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Воспроизвести команду из ROM2`},
			}),
			PlayFromRom3: control.NewPushbuttonControl(client, name, "Play from ROM3", control.Meta{
				Type: "pushbutton",

				Order:    25,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Воспроизвести команду из ROM3`},
			}),
			PlayFromRom4: control.NewPushbuttonControl(client, name, "Play from ROM4", control.Meta{
				Type: "pushbutton",

				Order:    26,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Воспроизвести команду из ROM4`},
			}),
			PlayFromRom5: control.NewPushbuttonControl(client, name, "Play from ROM5", control.Meta{
				Type: "pushbutton",

				Order:    27,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Воспроизвести команду из ROM5`},
			}),
			PlayFromRom6: control.NewPushbuttonControl(client, name, "Play from ROM6", control.Meta{
				Type: "pushbutton",

				Order:    28,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Воспроизвести команду из ROM6`},
			}),
			PlayFromRom7: control.NewPushbuttonControl(client, name, "Play from ROM7", control.Meta{
				Type: "pushbutton",

				Order:    29,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Воспроизвести команду из ROM7`},
			}),
			Serial: control.NewTextControl(client, name, "Serial", control.Meta{
				Type: "text",

				Order:    30,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Серийный номер`},
			}),
		}

		instanceWbMswV4151 = &WbMswV4151{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceWbMswV4151
}
