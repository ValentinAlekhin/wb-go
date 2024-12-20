package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type Dooya0X0104Controls struct {
	Position       *controls.RangeControl
	Open           *controls.PushbuttonControl
	Close          *controls.PushbuttonControl
	Stop           *controls.PushbuttonControl
	FactoryDefault *controls.PushbuttonControl
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
		deviceName := fmt.Sprintf("%s_%s", "dooya", "0x0104")
		controlList := &Dooya0X0104Controls{
			Position:       controls.NewRangeControl(client, deviceName, "Position"),
			Open:           controls.NewPushbuttonControl(client, deviceName, "Open"),
			Close:          controls.NewPushbuttonControl(client, deviceName, "Close"),
			Stop:           controls.NewPushbuttonControl(client, deviceName, "Stop"),
			FactoryDefault: controls.NewPushbuttonControl(client, deviceName, "Factory Default"),
		}

		instanceDooya0X0104 = &Dooya0X0104{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceDooya0X0104
}
