package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbMswV4151Controls struct {
	Illuminance       *controls.ValueControl
	PlayFromRam       *controls.PushbuttonControl
	LearnToRom4       *controls.SwitchControl
	LearnToRom5       *controls.SwitchControl
	LearnToRom6       *controls.SwitchControl
	LearnToRom7       *controls.SwitchControl
	PlayFromRom3      *controls.PushbuttonControl
	PlayFromRom4      *controls.PushbuttonControl
	PlayFromRom5      *controls.PushbuttonControl
	PlayFromRom6      *controls.PushbuttonControl
	Temperature       *controls.ValueControl
	Humidity          *controls.ValueControl
	Co2               *controls.ValueControl
	AirQualityVoc     *controls.ValueControl
	SoundLevel        *controls.ValueControl
	MaxMotion         *controls.ValueControl
	CurrentMotion     *controls.ValueControl
	Buzzer            *controls.SwitchControl
	RedLed            *controls.SwitchControl
	GreenLed          *controls.SwitchControl
	LedPeriods        *controls.RangeControl
	LedGlowDurationms *controls.RangeControl
	LearnToRam        *controls.SwitchControl
	LearnToRom1       *controls.SwitchControl
	LearnToRom2       *controls.SwitchControl
	LearnToRom3       *controls.SwitchControl
	PlayFromRom1      *controls.PushbuttonControl
	PlayFromRom2      *controls.PushbuttonControl
	PlayFromRom7      *controls.PushbuttonControl
	Serial            *controls.TextControl
}

type WbMswV4151 struct {
	Name     string
	Address  string
	Controls *WbMswV4151Controls
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
		deviceName := fmt.Sprintf("%s_%s", "wb-msw-v4", "151")
		controlList := &WbMswV4151Controls{
			Illuminance:       controls.NewValueControl(client, deviceName, "Illuminance"),
			PlayFromRam:       controls.NewPushbuttonControl(client, deviceName, "Play from RAM"),
			LearnToRom4:       controls.NewSwitchControl(client, deviceName, "Learn to ROM4"),
			LearnToRom5:       controls.NewSwitchControl(client, deviceName, "Learn to ROM5"),
			LearnToRom6:       controls.NewSwitchControl(client, deviceName, "Learn to ROM6"),
			LearnToRom7:       controls.NewSwitchControl(client, deviceName, "Learn to ROM7"),
			PlayFromRom3:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM3"),
			PlayFromRom4:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM4"),
			PlayFromRom5:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM5"),
			PlayFromRom6:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM6"),
			Temperature:       controls.NewValueControl(client, deviceName, "Temperature"),
			Humidity:          controls.NewValueControl(client, deviceName, "Humidity"),
			Co2:               controls.NewValueControl(client, deviceName, "CO2"),
			AirQualityVoc:     controls.NewValueControl(client, deviceName, "Air Quality (VOC)"),
			SoundLevel:        controls.NewValueControl(client, deviceName, "Sound Level"),
			MaxMotion:         controls.NewValueControl(client, deviceName, "Max Motion"),
			CurrentMotion:     controls.NewValueControl(client, deviceName, "Current Motion"),
			Buzzer:            controls.NewSwitchControl(client, deviceName, "Buzzer"),
			RedLed:            controls.NewSwitchControl(client, deviceName, "Red LED"),
			GreenLed:          controls.NewSwitchControl(client, deviceName, "Green LED"),
			LedPeriods:        controls.NewRangeControl(client, deviceName, "LED Period (s)"),
			LedGlowDurationms: controls.NewRangeControl(client, deviceName, "LED Glow Duration (ms)"),
			LearnToRam:        controls.NewSwitchControl(client, deviceName, "Learn to RAM"),
			LearnToRom1:       controls.NewSwitchControl(client, deviceName, "Learn to ROM1"),
			LearnToRom2:       controls.NewSwitchControl(client, deviceName, "Learn to ROM2"),
			LearnToRom3:       controls.NewSwitchControl(client, deviceName, "Learn to ROM3"),
			PlayFromRom1:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM1"),
			PlayFromRom2:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM2"),
			PlayFromRom7:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM7"),
			Serial:            controls.NewTextControl(client, deviceName, "Serial"),
		}

		instanceWbMswV4151 = &WbMswV4151{
			Name:     deviceName,
			Address:  "151",
			Controls: controlList,
		}
	})

	return instanceWbMswV4151
}
