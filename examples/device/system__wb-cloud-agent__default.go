package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemWbCloudAgentDefaultControls struct {
	Status         *control.TextControl
	ActivationLink *control.TextControl
	CloudBaseUrl   *control.TextControl
}

type SystemWbCloudAgentDefault struct {
	name     string
	Controls *SystemWbCloudAgentDefaultControls
}

func (w *SystemWbCloudAgentDefault) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemWbCloudAgentDefault     sync.Once
	instanceSystemWbCloudAgentDefault *SystemWbCloudAgentDefault
)

func NewSystemWbCloudAgentDefault(client *mqtt.Client) *SystemWbCloudAgentDefault {
	onceSystemWbCloudAgentDefault.Do(func() {
		name := "system__wb-cloud-agent__default"

		controlList := &SystemWbCloudAgentDefaultControls{
			Status: control.NewTextControl(client, name, "status", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Status`},
			}),
			ActivationLink: control.NewTextControl(client, name, "activation_link", control.Meta{
				Type: "text",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Link`},
			}),
			CloudBaseUrl: control.NewTextControl(client, name, "cloud_base_url", control.Meta{
				Type: "text",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `URL`},
			}),
		}

		instanceSystemWbCloudAgentDefault = &SystemWbCloudAgentDefault{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemWbCloudAgentDefault
}
