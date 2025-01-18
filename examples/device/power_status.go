package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type PowerStatusControls struct {
	Vin *control.ValueControl
}

type PowerStatus struct {
	name     string
	Controls *PowerStatusControls
}

func (w *PowerStatus) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	oncePowerStatus     sync.Once
	instancePowerStatus *PowerStatus
)

func NewPowerStatus(client mqtt.ClientInterface) *PowerStatus {
	oncePowerStatus.Do(func() {
		name := "power_status"

		controlList := &PowerStatusControls{
			Vin: control.NewValueControl(client, name, "Vin", control.Meta{
				Type: "voltage",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Input voltage`, "ru": `Входное напряжение`},
			}),
		}

		instancePowerStatus = &PowerStatus{
			name:     name,
			Controls: controlList,
		}
	})

	return instancePowerStatus
}
