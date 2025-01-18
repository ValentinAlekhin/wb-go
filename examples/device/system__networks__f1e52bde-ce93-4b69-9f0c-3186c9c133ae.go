package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworksF1E52BdeCe934B699F0C3186C9C133AeControls struct {
	Name *control.TextControl
}

type SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae struct {
	name     string
	Controls *SystemNetworksF1E52BdeCe934B699F0C3186C9C133AeControls
}

func (w *SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae     sync.Once
	instanceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae *SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae
)

func NewSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae(client mqtt.ClientInterface) *SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae {
	onceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae.Do(func() {
		name := "system__networks__f1e52bde-ce93-4b69-9f0c-3186c9c133ae"

		controlList := &SystemNetworksF1E52BdeCe934B699F0C3186C9C133AeControls{
			Name: control.NewTextControl(client, name, "Name", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
		}

		instanceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae = &SystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworksF1E52BdeCe934B699F0C3186C9C133Ae
}
