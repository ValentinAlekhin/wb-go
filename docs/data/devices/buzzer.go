package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type BuzzerControls struct {
	Enabled   *controls.SwitchControl
	Frequency *controls.RangeControl
	Volume    *controls.RangeControl
}

type Buzzer struct {
	Name     string
	Controls *BuzzerControls
}

var (
	onceBuzzer     sync.Once
	instanceBuzzer *Buzzer
)

func NewBuzzer(client *mqtt.Client) *Buzzer {
	onceBuzzer.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "buzzer", "")
		controlList := &BuzzerControls{
			Enabled:   controls.NewSwitchControl(client, deviceName, "enabled"),
			Frequency: controls.NewRangeControl(client, deviceName, "frequency"),
			Volume:    controls.NewRangeControl(client, deviceName, "volume"),
		}

		instanceBuzzer = &Buzzer{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceBuzzer
}
