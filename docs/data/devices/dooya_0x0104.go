package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type Dooya0X0104Controls struct {
	Position       *RangeControl
	Open           *PushbuttonControl
	Close          *PushbuttonControl
	Stop           *PushbuttonControl
	FactoryDefault *PushbuttonControl
}

type Dooya0X0104 struct {
	Name     string
	Controls *Dooya0X0104Controls
}

var (
	onceDooya0X0104     sync.Once
	instanceDooya0X0104 *Dooya0X0104
)

func NewDooya0X0104(client *mqtt.Client) *Dooya0X0104 {
	onceDooya0X0104.Do(func() {
		name := "dooya"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "0x0104")
		controls := &Dooya0X0104Controls{
			Position:       NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Position")),
			Open:           NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Open")),
			Close:          NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Close")),
			Stop:           NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Stop")),
			FactoryDefault: NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Factory Default")),
		}

		instanceDooya0X0104 = &Dooya0X0104{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceDooya0X0104
}
