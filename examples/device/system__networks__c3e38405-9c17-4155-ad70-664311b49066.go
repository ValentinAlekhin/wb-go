package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworksC3E384059C174155Ad70664311B49066Controls struct {
	Name *control.TextControl
}

type SystemNetworksC3E384059C174155Ad70664311B49066 struct {
	name     string
	Controls *SystemNetworksC3E384059C174155Ad70664311B49066Controls
}

func (w *SystemNetworksC3E384059C174155Ad70664311B49066) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworksC3E384059C174155Ad70664311B49066     sync.Once
	instanceSystemNetworksC3E384059C174155Ad70664311B49066 *SystemNetworksC3E384059C174155Ad70664311B49066
)

func NewSystemNetworksC3E384059C174155Ad70664311B49066(client mqtt.ClientInterface) *SystemNetworksC3E384059C174155Ad70664311B49066 {
	onceSystemNetworksC3E384059C174155Ad70664311B49066.Do(func() {
		name := "system__networks__c3e38405-9c17-4155-ad70-664311b49066"

		controlList := &SystemNetworksC3E384059C174155Ad70664311B49066Controls{
			Name: control.NewTextControl(client, name, "Name", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
		}

		instanceSystemNetworksC3E384059C174155Ad70664311B49066 = &SystemNetworksC3E384059C174155Ad70664311B49066{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworksC3E384059C174155Ad70664311B49066
}
