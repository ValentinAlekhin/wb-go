package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMswV4151Controls struct {
	Temperature       *ValueControl
	Humidity          *ValueControl
	Co2               *ValueControl
	AirQualityVoc     *ValueControl
	SoundLevel        *ValueControl
	Illuminance       *ValueControl
	MaxMotion         *ValueControl
	CurrentMotion     *ValueControl
	Buzzer            *SwitchControl
	RedLed            *SwitchControl
	GreenLed          *SwitchControl
	LedPeriods        *RangeControl
	LedGlowDurationms *RangeControl
	LearnToRam        *SwitchControl
	PlayFromRam       *PushbuttonControl
	LearnToRom1       *SwitchControl
	LearnToRom2       *SwitchControl
	LearnToRom3       *SwitchControl
	LearnToRom4       *SwitchControl
	LearnToRom5       *SwitchControl
	LearnToRom6       *SwitchControl
	LearnToRom7       *SwitchControl
	PlayFromRom1      *PushbuttonControl
	PlayFromRom2      *PushbuttonControl
	PlayFromRom3      *PushbuttonControl
	PlayFromRom4      *PushbuttonControl
	PlayFromRom5      *PushbuttonControl
	PlayFromRom6      *PushbuttonControl
	PlayFromRom7      *PushbuttonControl
	Serial            *TextControl
}

type WbMswV4151 struct {
	Name          string
	ModbusAddress int32
	Controls      *WbMswV4151Controls
}

var (
	onceWbMswV4151     sync.Once
	instanceWbMswV4151 *WbMswV4151
)

func NewWbMswV4151(client *mqtt.Client) *WbMswV4151 {
	onceWbMswV4151.Do(func() {
		name := "wb-msw-v4"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "151")
		controls := &WbMswV4151Controls{
			Temperature:       NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Temperature")),
			Humidity:          NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Humidity")),
			Co2:               NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "CO2")),
			AirQualityVoc:     NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Air Quality (VOC)")),
			SoundLevel:        NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Sound Level")),
			Illuminance:       NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Illuminance")),
			MaxMotion:         NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Max Motion")),
			CurrentMotion:     NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Current Motion")),
			Buzzer:            NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Buzzer")),
			RedLed:            NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Red LED")),
			GreenLed:          NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Green LED")),
			LedPeriods:        NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "LED Period (s)")),
			LedGlowDurationms: NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "LED Glow Duration (ms)")),
			LearnToRam:        NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Learn to RAM")),
			PlayFromRam:       NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Play from RAM")),
			LearnToRom1:       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Learn to ROM1")),
			LearnToRom2:       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Learn to ROM2")),
			LearnToRom3:       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Learn to ROM3")),
			LearnToRom4:       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Learn to ROM4")),
			LearnToRom5:       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Learn to ROM5")),
			LearnToRom6:       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Learn to ROM6")),
			LearnToRom7:       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Learn to ROM7")),
			PlayFromRom1:      NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Play from ROM1")),
			PlayFromRom2:      NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Play from ROM2")),
			PlayFromRom3:      NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Play from ROM3")),
			PlayFromRom4:      NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Play from ROM4")),
			PlayFromRom5:      NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Play from ROM5")),
			PlayFromRom6:      NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Play from ROM6")),
			PlayFromRom7:      NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Play from ROM7")),
			Serial:            NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Serial")),
		}

		instanceWbMswV4151 = &WbMswV4151{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceWbMswV4151
}
