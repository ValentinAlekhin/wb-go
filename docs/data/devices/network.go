package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type NetworkControls struct {
	ActiveConnections            *TextControl
	DefaultInterface             *TextControl
	Ethernet2Ip                  *TextControl
	Ethernet2IpOnlineStatus      *SwitchControl
	EthernetIp                   *TextControl
	EthernetIpOnlineStatus       *SwitchControl
	GprsIp                       *TextControl
	GprsIpOnlineStatus           *SwitchControl
	InternetConnection           *TextControl
	WiFi2Ip                      *TextControl
	WiFi2IpOnlineStatus          *SwitchControl
	WiFiIp                       *TextControl
	WiFiIpOnlineStatus           *SwitchControl
	Ethernet2IpConnectionEnabled *SwitchControl
	EthernetIpConnectionEnabled  *SwitchControl
	GprsIpConnectionEnabled      *SwitchControl
	WiFi2IpConnectionEnabled     *SwitchControl
	WiFiIpConnectionEnabled      *SwitchControl
}

type Network struct {
	Name          string
	ModbusAddress int32
	Controls      *NetworkControls
}

var (
	onceNetwork     sync.Once
	instanceNetwork *Network
)

func NewNetwork(client *mqtt.Client) *Network {
	onceNetwork.Do(func() {
		name := "network"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "")
		controls := &NetworkControls{
			ActiveConnections:            NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Active Connections")),
			DefaultInterface:             NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Default Interface")),
			Ethernet2Ip:                  NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Ethernet 2 IP")),
			Ethernet2IpOnlineStatus:      NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Ethernet 2 IP Online Status")),
			EthernetIp:                   NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Ethernet IP")),
			EthernetIpOnlineStatus:       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Ethernet IP Online Status")),
			GprsIp:                       NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "GPRS IP")),
			GprsIpOnlineStatus:           NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "GPRS IP Online Status")),
			InternetConnection:           NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Internet Connection")),
			WiFi2Ip:                      NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Wi-Fi 2 IP")),
			WiFi2IpOnlineStatus:          NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Wi-Fi 2 IP Online Status")),
			WiFiIp:                       NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Wi-Fi IP")),
			WiFiIpOnlineStatus:           NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Wi-Fi IP Online Status")),
			Ethernet2IpConnectionEnabled: NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Ethernet 2 IP Connection Enabled")),
			EthernetIpConnectionEnabled:  NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Ethernet IP Connection Enabled")),
			GprsIpConnectionEnabled:      NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "GPRS IP Connection Enabled")),
			WiFi2IpConnectionEnabled:     NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Wi-Fi 2 IP Connection Enabled")),
			WiFiIpConnectionEnabled:      NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Wi-Fi IP Connection Enabled")),
		}

		instanceNetwork = &Network{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceNetwork
}
