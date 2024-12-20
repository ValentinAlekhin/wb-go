package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
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
	Name     string
	Controls *WbMswV4151Controls
}

var (
	onceWbMswV4151     sync.Once
	instanceWbMswV4151 *WbMswV4151
)

func NewWbMswV4151(client *mqtt.Client) *WbMswV4151 {
	onceWbMswV4151.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "wb-msw-v4", "151")
		controlList := &WbMswV4151Controls{
			Temperature:       controls.NewValueControl(client, deviceName, "Temperature"),
			Humidity:          controls.NewValueControl(client, deviceName, "Humidity"),
			Co2:               controls.NewValueControl(client, deviceName, "CO2"),
			AirQualityVoc:     controls.NewValueControl(client, deviceName, "Air Quality (VOC)"),
			SoundLevel:        controls.NewValueControl(client, deviceName, "Sound Level"),
			Illuminance:       controls.NewValueControl(client, deviceName, "Illuminance"),
			MaxMotion:         controls.NewValueControl(client, deviceName, "Max Motion"),
			CurrentMotion:     controls.NewValueControl(client, deviceName, "Current Motion"),
			Buzzer:            controls.NewSwitchControl(client, deviceName, "Buzzer"),
			RedLed:            controls.NewSwitchControl(client, deviceName, "Red LED"),
			GreenLed:          controls.NewSwitchControl(client, deviceName, "Green LED"),
			LedPeriods:        controls.NewRangeControl(client, deviceName, "LED Period (s)"),
			LedGlowDurationms: controls.NewRangeControl(client, deviceName, "LED Glow Duration (ms)"),
			LearnToRam:        controls.NewSwitchControl(client, deviceName, "Learn to RAM"),
			PlayFromRam:       controls.NewPushbuttonControl(client, deviceName, "Play from RAM"),
			LearnToRom1:       controls.NewSwitchControl(client, deviceName, "Learn to ROM1"),
			LearnToRom2:       controls.NewSwitchControl(client, deviceName, "Learn to ROM2"),
			LearnToRom3:       controls.NewSwitchControl(client, deviceName, "Learn to ROM3"),
			LearnToRom4:       controls.NewSwitchControl(client, deviceName, "Learn to ROM4"),
			LearnToRom5:       controls.NewSwitchControl(client, deviceName, "Learn to ROM5"),
			LearnToRom6:       controls.NewSwitchControl(client, deviceName, "Learn to ROM6"),
			LearnToRom7:       controls.NewSwitchControl(client, deviceName, "Learn to ROM7"),
			PlayFromRom1:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM1"),
			PlayFromRom2:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM2"),
			PlayFromRom3:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM3"),
			PlayFromRom4:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM4"),
			PlayFromRom5:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM5"),
			PlayFromRom6:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM6"),
			PlayFromRom7:      controls.NewPushbuttonControl(client, deviceName, "Play from ROM7"),
			Serial:            controls.NewTextControl(client, deviceName, "Serial"),
		}

		instanceWbMswV4151 = &WbMswV4151{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceWbMswV4151
}
