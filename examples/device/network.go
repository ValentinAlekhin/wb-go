package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type NetworkControls struct {
	ActiveConnections            *control.TextControl
	DefaultInterface             *control.TextControl
	Ethernet2Ip                  *control.TextControl
	Ethernet2IpOnlineStatus      *control.SwitchControl
	EthernetIp                   *control.TextControl
	EthernetIpOnlineStatus       *control.SwitchControl
	GprsIp                       *control.TextControl
	GprsIpOnlineStatus           *control.SwitchControl
	InternetConnection           *control.TextControl
	WiFi2Ip                      *control.TextControl
	WiFi2IpOnlineStatus          *control.SwitchControl
	WiFiIp                       *control.TextControl
	WiFiIpOnlineStatus           *control.SwitchControl
	Ethernet2IpConnectionEnabled *control.SwitchControl
	EthernetIpConnectionEnabled  *control.SwitchControl
	GprsIpConnectionEnabled      *control.SwitchControl
	WiFi2IpConnectionEnabled     *control.SwitchControl
	WiFiIpConnectionEnabled      *control.SwitchControl
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

func NewNetwork(client *mqtt.Client) *Network {
	onceNetwork.Do(func() {
		name := "network"

		controlList := &NetworkControls{
			ActiveConnections: control.NewTextControl(client, name, "Active Connections", control.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Active Connections`, "ru": `Активные соединения`},
			}),
			DefaultInterface: control.NewTextControl(client, name, "Default Interface", control.Meta{
				Type: "text",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Default Interface`, "ru": `Интерфейс по умолчанию`},
			}),
			Ethernet2Ip: control.NewTextControl(client, name, "Ethernet 2 IP", control.Meta{
				Type: "text",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			Ethernet2IpOnlineStatus: control.NewSwitchControl(client, name, "Ethernet 2 IP Online Status", control.Meta{
				Type: "switch",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Ethernet 2 Internet Access`, "ru": `Ethernet 2 Доступ к интернету`},
			}),
			EthernetIp: control.NewTextControl(client, name, "Ethernet IP", control.Meta{
				Type: "text",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Ethernet 1 IP`, "ru": `Ethernet 1 IP`},
			}),
			EthernetIpOnlineStatus: control.NewSwitchControl(client, name, "Ethernet IP Online Status", control.Meta{
				Type: "switch",

				Order:    5,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Ethernet 1 Internet Access`, "ru": `Ethernet 1 Доступ к интернету`},
			}),
			GprsIp: control.NewTextControl(client, name, "GPRS IP", control.Meta{
				Type: "text",

				Order:    16,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			GprsIpOnlineStatus: control.NewSwitchControl(client, name, "GPRS IP Online Status", control.Meta{
				Type: "switch",

				Order:    17,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `GPRS IP Internet Access`, "ru": `GPRS IP Доступ к интернету`},
			}),
			InternetConnection: control.NewTextControl(client, name, "Internet Connection", control.Meta{
				Type: "text",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Internet Connection`, "ru": `Интернет соединение`},
			}),
			WiFi2Ip: control.NewTextControl(client, name, "Wi-Fi 2 IP", control.Meta{
				Type: "text",

				Order:    13,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Wi-Fi 2 IP`, "ru": `Wi-Fi 2 IP`},
			}),
			WiFi2IpOnlineStatus: control.NewSwitchControl(client, name, "Wi-Fi 2 IP Online Status", control.Meta{
				Type: "switch",

				Order:    14,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Wi-Fi 2 Internet Access`, "ru": `Wi-Fi 2 Доступ к интернету`},
			}),
			WiFiIp: control.NewTextControl(client, name, "Wi-Fi IP", control.Meta{
				Type: "text",

				Order:    10,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Wi-Fi 1 IP`, "ru": `Wi-Fi 1 IP`},
			}),
			WiFiIpOnlineStatus: control.NewSwitchControl(client, name, "Wi-Fi IP Online Status", control.Meta{
				Type: "switch",

				Order:    11,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Wi-Fi 1 Internet Access`, "ru": `Wi-Fi 1 Доступ к интернету`},
			}),
			Ethernet2IpConnectionEnabled: control.NewSwitchControl(client, name, "Ethernet 2 IP Connection Enabled", control.Meta{
				Type: "switch",

				Order:    9,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Ethernet 2 Enabled`, "ru": `Ethernet 2 Включен`},
			}),
			EthernetIpConnectionEnabled: control.NewSwitchControl(client, name, "Ethernet IP Connection Enabled", control.Meta{
				Type: "switch",

				Order:    6,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Ethernet 1 Enabled`, "ru": `Ethernet 1 Включен`},
			}),
			GprsIpConnectionEnabled: control.NewSwitchControl(client, name, "GPRS IP Connection Enabled", control.Meta{
				Type: "switch",

				Order:    18,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `GPRS IP Enabled`, "ru": `GPRS IP Включен`},
			}),
			WiFi2IpConnectionEnabled: control.NewSwitchControl(client, name, "Wi-Fi 2 IP Connection Enabled", control.Meta{
				Type: "switch",

				Order:    15,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Wi-Fi 2 Enabled`, "ru": `Wi-Fi 2 Включен`},
			}),
			WiFiIpConnectionEnabled: control.NewSwitchControl(client, name, "Wi-Fi IP Connection Enabled", control.Meta{
				Type: "switch",

				Order:    12,
				ReadOnly: true,
				Title:    control.MultilingualText{"en": `Wi-Fi 1 Enabled`, "ru": `Wi-Fi 1 Включен`},
			}),
		}

		instanceNetwork = &Network{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceNetwork
}
