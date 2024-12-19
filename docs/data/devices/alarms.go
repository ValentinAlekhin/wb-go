package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type AlarmsControls struct {
	Log *TextControl
}

type Alarms struct {
	Name     string
	Controls *AlarmsControls
}

var (
	onceAlarms     sync.Once
	instanceAlarms *Alarms
)

func NewAlarms(client *mqtt.Client) *Alarms {
	onceAlarms.Do(func() {
		name := "alarms"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "")
		controls := &AlarmsControls{
			Log: NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "log")),
		}

		instanceAlarms = &Alarms{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceAlarms
}
