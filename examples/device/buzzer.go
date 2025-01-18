package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type BuzzerControls struct {
	Enabled *control.SwitchControl
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
		}

		instanceBuzzer = &Buzzer{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceBuzzer
}
