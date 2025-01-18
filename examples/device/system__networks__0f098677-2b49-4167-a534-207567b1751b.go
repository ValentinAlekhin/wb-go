package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemNetworks0F0986772B494167A534207567B1751BControls struct {
	Name *control.TextControl
}

type SystemNetworks0F0986772B494167A534207567B1751B struct {
	name     string
	Controls *SystemNetworks0F0986772B494167A534207567B1751BControls
}

func (w *SystemNetworks0F0986772B494167A534207567B1751B) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceSystemNetworks0F0986772B494167A534207567B1751B     sync.Once
	instanceSystemNetworks0F0986772B494167A534207567B1751B *SystemNetworks0F0986772B494167A534207567B1751B
)

func NewSystemNetworks0F0986772B494167A534207567B1751B(client mqtt.ClientInterface) *SystemNetworks0F0986772B494167A534207567B1751B {
	onceSystemNetworks0F0986772B494167A534207567B1751B.Do(func() {
		name := "system__networks__0f098677-2b49-4167-a534-207567b1751b"

		controlList := &SystemNetworks0F0986772B494167A534207567B1751BControls{
			Name: control.NewTextControl(client, name, "Name", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
		}

		instanceSystemNetworks0F0986772B494167A534207567B1751B = &SystemNetworks0F0986772B494167A534207567B1751B{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceSystemNetworks0F0986772B494167A534207567B1751B
}
