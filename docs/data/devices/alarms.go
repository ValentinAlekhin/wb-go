package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type AlarmsControls struct {
	Log *controls.TextControl
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
		deviceName := fmt.Sprintf("%s_%s", "alarms", "")
		controlList := &AlarmsControls{
			Log: controls.NewTextControl(client, deviceName, "log"),
		}

		instanceAlarms = &Alarms{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceAlarms
}
