package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type NetworkControls struct {
	ActiveConnections *control.TextControl
}

type Network struct {
	name     string
	Controls *NetworkControls
}

func (w *Network) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceNetwork     sync.Once
	instanceNetwork *Network
)

func NewNetwork(client mqtt.ClientInterface) *Network {
	onceNetwork.Do(func() {
		name := "network"

		controlList := &NetworkControls{
			ActiveConnections: control.NewTextControl(client, name, "Active Connections", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Active Connections`, "ru": `Активные соединения`},
			}),
		}

		instanceNetwork = &Network{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceNetwork
}
