package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type SystemControls struct {
	BatchNo            *TextControl
	CurrentUptime      *TextControl
	DtsVersion         *TextControl
	HwRevision         *TextControl
	ManufacturingDate  *TextControl
	Reboot             *PushbuttonControl
	ReleaseName        *TextControl
	ReleaseSuite       *TextControl
	ShortSn            *TextControl
	TemperatureGrade   *TextControl
	Status             *TextControl
	ActivationLink     *TextControl
	CloudBaseUrl       *TextControl
	Name               *TextControl
	Uuid               *TextControl
	Type               *TextControl
	Active             *SwitchControl
	Device             *TextControl
	State              *TextControl
	Address            *TextControl
	Connectivity       *SwitchControl
	UpDown             *PushbuttonControl
	Operator           *TextControl
	SignalQuality      *TextControl
	AccessTechnologies *TextControl
}

type System struct {
	Name     string
	Controls *SystemControls
}

var (
	onceSystem     sync.Once
	instanceSystem *System
)

func NewSystem(client *mqtt.Client) *System {
	onceSystem.Do(func() {
		name := "system"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "")
		controls := &SystemControls{
			BatchNo:            NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Batch No")),
			CurrentUptime:      NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Current uptime")),
			DtsVersion:         NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "DTS Version")),
			HwRevision:         NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "HW Revision")),
			ManufacturingDate:  NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Manufacturing Date")),
			Reboot:             NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Reboot")),
			ReleaseName:        NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Release name")),
			ReleaseSuite:       NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Release suite")),
			ShortSn:            NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Short SN")),
			TemperatureGrade:   NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Temperature Grade")),
			Status:             NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "status")),
			ActivationLink:     NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "activation_link")),
			CloudBaseUrl:       NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "cloud_base_url")),
			Name:               NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Name")),
			Uuid:               NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "UUID")),
			Type:               NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Type")),
			Active:             NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Active")),
			Device:             NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Device")),
			State:              NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "State")),
			Address:            NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Address")),
			Connectivity:       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Connectivity")),
			UpDown:             NewPushbuttonControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "UpDown")),
			Operator:           NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Operator")),
			SignalQuality:      NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "SignalQuality")),
			AccessTechnologies: NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "AccessTechnologies")),
		}

		instanceSystem = &System{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceSystem
}
