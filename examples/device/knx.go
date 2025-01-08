package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type KnxControls struct {
	Data *control.TextControl
}

type Knx struct {
	name     string
	Controls *KnxControls
}

func (w *Knx) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceKnx     sync.Once
	instanceKnx *Knx
)

func NewKnx(client mqtt.ClientInterface) *Knx {
	onceKnx.Do(func() {
		name := "knx"

		controlList := &KnxControls{
			Data: control.NewTextControl(client, name, "data", control.Meta{
				Type: "text",

				Order:    0,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Message`, "ru": `Сообщение`},
			}),
		}

		instanceKnx = &Knx{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceKnx
}
