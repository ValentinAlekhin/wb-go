package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type KnxControls struct {
	Data *controls.TextControl
}

type Knx struct {
	Name     string
	Controls *KnxControls
}

var (
	onceKnx     sync.Once
	instanceKnx *Knx
)

func NewKnx(client *mqtt.Client) *Knx {
	onceKnx.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "knx", "")
		controlList := &KnxControls{
			Data: controls.NewTextControl(client, deviceName, "data"),
		}

		instanceKnx = &Knx{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceKnx
}
