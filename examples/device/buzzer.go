package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type BuzzerControls struct {
	Enabled   *control.SwitchControl
	Frequency *control.RangeControl
	Volume    *control.RangeControl
}

type Buzzer struct {
	name     string
	Controls *BuzzerControls
}

func (w *Buzzer) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceBuzzer     sync.Once
	instanceBuzzer *Buzzer
)

func NewBuzzer(client mqtt.ClientInterface) *Buzzer {
	onceBuzzer.Do(func() {
		name := "buzzer"

		controlList := &BuzzerControls{
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
			Controls: controlList,
		}
	})

	return instanceBuzzer
}
