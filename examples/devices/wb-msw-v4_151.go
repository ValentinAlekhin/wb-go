package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbMswV4151Controls struct {
	Temperature       *controls.ValueControl
	Humidity          *controls.ValueControl
	Co2               *controls.ValueControl
	AirQualityVoc     *controls.ValueControl
	SoundLevel        *controls.ValueControl
	Illuminance       *controls.ValueControl
	MaxMotion         *controls.ValueControl
	CurrentMotion     *controls.ValueControl
	Buzzer            *controls.SwitchControl
	RedLed            *controls.SwitchControl
	GreenLed          *controls.SwitchControl
	LedPeriods        *controls.RangeControl
	LedGlowDurationms *controls.RangeControl
	LearnToRam        *controls.SwitchControl
	PlayFromRam       *controls.PushbuttonControl
	LearnToRom1       *controls.SwitchControl
	LearnToRom2       *controls.SwitchControl
	LearnToRom3       *controls.SwitchControl
	LearnToRom4       *controls.SwitchControl
	LearnToRom5       *controls.SwitchControl
	LearnToRom6       *controls.SwitchControl
	LearnToRom7       *controls.SwitchControl
	PlayFromRom1      *controls.PushbuttonControl
	PlayFromRom2      *controls.PushbuttonControl
	PlayFromRom3      *controls.PushbuttonControl
	PlayFromRom4      *controls.PushbuttonControl
	PlayFromRom5      *controls.PushbuttonControl
	PlayFromRom6      *controls.PushbuttonControl
	PlayFromRom7      *controls.PushbuttonControl
	Serial            *controls.TextControl
}

type WbMswV4151 struct {
	name     string
	device   string
	address  string
	Controls *WbMswV4151Controls
}

func (w *WbMswV4151) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbMswV4151) GetControlsInfo() []controls.ControlInfo {
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
	onceWbMswV4151     sync.Once
	instanceWbMswV4151 *WbMswV4151
)

func NewWbMswV4151(client *mqtt.Client) *WbMswV4151 {
	onceWbMswV4151.Do(func() {
		device := "wb-msw-v4"
		address := "151"
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &WbMswV4151Controls{
			Temperature: controls.NewValueControl(client, name, "Temperature", controls.Meta{
				Type:  "value",
				Units: "deg C",

				Order:    1,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Температура`},
			}),
			Humidity: controls.NewValueControl(client, name, "Humidity", controls.Meta{
				Type:  "value",
				Units: "%, RH",

				Order:    2,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Влажность`},
			}),
			Co2: controls.NewValueControl(client, name, "CO2", controls.Meta{
				Type: "concentration",

				Order:    3,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `CO₂`, "ru": `Уровень CO₂`},
			}),
			AirQualityVoc: controls.NewValueControl(client, name, "Air Quality (VOC)", controls.Meta{
				Type:  "value",
				Units: "ppb",

				Order:    4,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Качество воздуха (VOC)`},
			}),
			SoundLevel: controls.NewValueControl(client, name, "Sound Level", controls.Meta{
				Type: "sound_level",

				Order:    5,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Уровень шума`},
			}),
			Illuminance: controls.NewValueControl(client, name, "Illuminance", controls.Meta{
				Type: "lux",

				Order:    6,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Освещенность`},
			}),
			MaxMotion: controls.NewValueControl(client, name, "Max Motion", controls.Meta{
				Type: "value",

				Order:    7,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Максимальное движение`},
			}),
			CurrentMotion: controls.NewValueControl(client, name, "Current Motion", controls.Meta{
				Type: "value",

				Order:    8,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Текущее движение`},
			}),
			Buzzer: controls.NewSwitchControl(client, name, "Buzzer", controls.Meta{
				Type: "switch",

				Order:    9,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Зуммер`},
			}),
			RedLed: controls.NewSwitchControl(client, name, "Red LED", controls.Meta{
				Type: "switch",

				Order:    10,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Красный светодиод`},
			}),
			GreenLed: controls.NewSwitchControl(client, name, "Green LED", controls.Meta{
				Type: "switch",

				Order:    11,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Зеленый светодиод`},
			}),
			LedPeriods: controls.NewRangeControl(client, name, "LED Period (s)", controls.Meta{
				Type: "range",

				Max: 10,

				Order:    12,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Период включения светодиодов (с)`},
			}),
			LedGlowDurationms: controls.NewRangeControl(client, name, "LED Glow Duration (ms)", controls.Meta{
				Type: "range",

				Max: 50,

				Order:    13,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Длительность включения светодиодов (мс)`},
			}),
			LearnToRam: controls.NewSwitchControl(client, name, "Learn to RAM", controls.Meta{
				Type: "switch",

				Order:    14,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Записать команду в RAM`},
			}),
			PlayFromRam: controls.NewPushbuttonControl(client, name, "Play from RAM", controls.Meta{
				Type: "pushbutton",

				Order:    15,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Воспроизвести команду из RAM`},
			}),
			LearnToRom1: controls.NewSwitchControl(client, name, "Learn to ROM1", controls.Meta{
				Type: "switch",

				Order:    16,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Записать команду в ROM1`},
			}),
			LearnToRom2: controls.NewSwitchControl(client, name, "Learn to ROM2", controls.Meta{
				Type: "switch",

				Order:    17,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Записать команду в ROM2`},
			}),
			LearnToRom3: controls.NewSwitchControl(client, name, "Learn to ROM3", controls.Meta{
				Type: "switch",

				Order:    18,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Записать команду в ROM3`},
			}),
			LearnToRom4: controls.NewSwitchControl(client, name, "Learn to ROM4", controls.Meta{
				Type: "switch",

				Order:    19,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Записать команду в ROM4`},
			}),
			LearnToRom5: controls.NewSwitchControl(client, name, "Learn to ROM5", controls.Meta{
				Type: "switch",

				Order:    20,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Записать команду в ROM5`},
			}),
			LearnToRom6: controls.NewSwitchControl(client, name, "Learn to ROM6", controls.Meta{
				Type: "switch",

				Order:    21,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Записать команду в ROM6`},
			}),
			LearnToRom7: controls.NewSwitchControl(client, name, "Learn to ROM7", controls.Meta{
				Type: "switch",

				Order:    22,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Записать команду в ROM7`},
			}),
			PlayFromRom1: controls.NewPushbuttonControl(client, name, "Play from ROM1", controls.Meta{
				Type: "pushbutton",

				Order:    23,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Воспроизвести команду из ROM1`},
			}),
			PlayFromRom2: controls.NewPushbuttonControl(client, name, "Play from ROM2", controls.Meta{
				Type: "pushbutton",

				Order:    24,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Воспроизвести команду из ROM2`},
			}),
			PlayFromRom3: controls.NewPushbuttonControl(client, name, "Play from ROM3", controls.Meta{
				Type: "pushbutton",

				Order:    25,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Воспроизвести команду из ROM3`},
			}),
			PlayFromRom4: controls.NewPushbuttonControl(client, name, "Play from ROM4", controls.Meta{
				Type: "pushbutton",

				Order:    26,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Воспроизвести команду из ROM4`},
			}),
			PlayFromRom5: controls.NewPushbuttonControl(client, name, "Play from ROM5", controls.Meta{
				Type: "pushbutton",

				Order:    27,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Воспроизвести команду из ROM5`},
			}),
			PlayFromRom6: controls.NewPushbuttonControl(client, name, "Play from ROM6", controls.Meta{
				Type: "pushbutton",

				Order:    28,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Воспроизвести команду из ROM6`},
			}),
			PlayFromRom7: controls.NewPushbuttonControl(client, name, "Play from ROM7", controls.Meta{
				Type: "pushbutton",

				Order:    29,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Воспроизвести команду из ROM7`},
			}),
			Serial: controls.NewTextControl(client, name, "Serial", controls.Meta{
				Type: "text",

				Order:    30,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Серийный номер`},
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
