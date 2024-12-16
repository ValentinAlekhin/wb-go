package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type BuzzerControls struct {
	Enabled   *SwitchControl
	Frequency *RangeControl
	Volume    *RangeControl
}

type Buzzer struct {
	Name          string
	ModbusAddress int32
	Controls      *BuzzerControls
}

var (
	onceBuzzer     sync.Once
	instanceBuzzer *Buzzer
)

func NewBuzzer(client *mqtt.Client) *Buzzer {
	onceBuzzer.Do(func() {
		name := "buzzer"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "")
		controls := &BuzzerControls{
			Enabled:   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "enabled")),
			Frequency: NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "frequency")),
			Volume:    NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "volume")),
		}

		instanceBuzzer = &Buzzer{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceBuzzer
}
