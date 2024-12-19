package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type KnxControls struct {
	Data *TextControl
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
		name := "knx"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "")
		controls := &KnxControls{
			Data: NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "data")),
		}

		instanceKnx = &Knx{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceKnx
}
