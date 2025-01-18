package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9Controls struct {
	Name *control.TextControl
}

type SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9 struct {
	name     string
	Controls *SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9Controls
}

func (w *SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9     sync.Once
	instanceSystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9 *SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9
)

func NewSystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9(client mqtt.ClientInterface) *SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9 {
	onceSystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9.Do(func() {
		name := "system__networks__8b9964d4-b8dd-34d3-a3ed-481840bcf8c9"

		controlList := &SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9Controls{
			Name: control.NewTextControl(client, name, "Name", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
		}

		instanceSystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9 = &SystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworks8B9964D4B8Dd34D3A3Ed481840Bcf8C9
}
