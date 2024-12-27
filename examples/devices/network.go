package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type NetworkControls struct {
	ActiveConnections            *controls.TextControl
	DefaultInterface             *controls.TextControl
	Ethernet2Ip                  *controls.TextControl
	Ethernet2IpOnlineStatus      *controls.SwitchControl
	EthernetIp                   *controls.TextControl
	EthernetIpOnlineStatus       *controls.SwitchControl
	GprsIp                       *controls.TextControl
	GprsIpOnlineStatus           *controls.SwitchControl
	InternetConnection           *controls.TextControl
	WiFi2Ip                      *controls.TextControl
	WiFi2IpOnlineStatus          *controls.SwitchControl
	WiFiIp                       *controls.TextControl
	WiFiIpOnlineStatus           *controls.SwitchControl
	Ethernet2IpConnectionEnabled *controls.SwitchControl
	EthernetIpConnectionEnabled  *controls.SwitchControl
	GprsIpConnectionEnabled      *controls.SwitchControl
	WiFi2IpConnectionEnabled     *controls.SwitchControl
	WiFiIpConnectionEnabled      *controls.SwitchControl
}

type Network struct {
	name     string
	device   string
	address  string
	Controls *NetworkControls
}

func (w *Network) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Network) GetControlsInfo() []controls.ControlInfo {
	var infoList []controls.ControlInfo

	// Получаем значение и тип структуры Controls
	controlsValue := reflect.ValueOf(w.Controls).Elem()
	controlsType := controlsValue.Type()

	// Проходимся по всем полям структуры Controls
	for i := 0; i < controlsValue.NumField(); i++ {
		field := controlsValue.Field(i)

		// Проверяем, что поле является указателем и не nil
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			// Проверяем, реализует ли поле метод GetInfo
			method := field.MethodByName("GetInfo")
			if method.IsValid() {
				// Вызываем метод GetInfo
				info := method.Call(nil)[0].Interface().(controls.ControlInfo)
				infoList = append(infoList, info)
			} else {
				fmt.Printf("Field %s does not implement GetInfo\n", controlsType.Field(i).Name)
			}
		}
	}

	return infoList
}

var (
	onceNetwork     sync.Once
	instanceNetwork *Network
)

func NewNetwork(client *mqtt.Client) *Network {
	onceNetwork.Do(func() {
		device := "network"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &NetworkControls{
			ActiveConnections: controls.NewTextControl(client, name, "Active Connections", controls.Meta{
				Type: "text",

				Order:    1,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Active Connections`, "ru": `Активные соединения`},
			}),
			DefaultInterface: controls.NewTextControl(client, name, "Default Interface", controls.Meta{
				Type: "text",

				Order:    2,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Default Interface`, "ru": `Интерфейс по умолчанию`},
			}),
			Ethernet2Ip: controls.NewTextControl(client, name, "Ethernet 2 IP", controls.Meta{
				Type: "text",

				Order:    7,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			Ethernet2IpOnlineStatus: controls.NewSwitchControl(client, name, "Ethernet 2 IP Online Status", controls.Meta{
				Type: "switch",

				Order:    8,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Ethernet 2 Internet Access`, "ru": `Ethernet 2 Доступ к интернету`},
			}),
			EthernetIp: controls.NewTextControl(client, name, "Ethernet IP", controls.Meta{
				Type: "text",

				Order:    4,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Ethernet 1 IP`, "ru": `Ethernet 1 IP`},
			}),
			EthernetIpOnlineStatus: controls.NewSwitchControl(client, name, "Ethernet IP Online Status", controls.Meta{
				Type: "switch",

				Order:    5,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Ethernet 1 Internet Access`, "ru": `Ethernet 1 Доступ к интернету`},
			}),
			GprsIp: controls.NewTextControl(client, name, "GPRS IP", controls.Meta{
				Type: "text",

				Order:    16,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			GprsIpOnlineStatus: controls.NewSwitchControl(client, name, "GPRS IP Online Status", controls.Meta{
				Type: "switch",

				Order:    17,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `GPRS IP Internet Access`, "ru": `GPRS IP Доступ к интернету`},
			}),
			InternetConnection: controls.NewTextControl(client, name, "Internet Connection", controls.Meta{
				Type: "text",

				Order:    3,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Internet Connection`, "ru": `Интернет соединение`},
			}),
			WiFi2Ip: controls.NewTextControl(client, name, "Wi-Fi 2 IP", controls.Meta{
				Type: "text",

				Order:    13,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Wi-Fi 2 IP`, "ru": `Wi-Fi 2 IP`},
			}),
			WiFi2IpOnlineStatus: controls.NewSwitchControl(client, name, "Wi-Fi 2 IP Online Status", controls.Meta{
				Type: "switch",

				Order:    14,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Wi-Fi 2 Internet Access`, "ru": `Wi-Fi 2 Доступ к интернету`},
			}),
			WiFiIp: controls.NewTextControl(client, name, "Wi-Fi IP", controls.Meta{
				Type: "text",

				Order:    10,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Wi-Fi 1 IP`, "ru": `Wi-Fi 1 IP`},
			}),
			WiFiIpOnlineStatus: controls.NewSwitchControl(client, name, "Wi-Fi IP Online Status", controls.Meta{
				Type: "switch",

				Order:    11,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Wi-Fi 1 Internet Access`, "ru": `Wi-Fi 1 Доступ к интернету`},
			}),
			Ethernet2IpConnectionEnabled: controls.NewSwitchControl(client, name, "Ethernet 2 IP Connection Enabled", controls.Meta{
				Type: "switch",

				Order:    9,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Ethernet 2 Enabled`, "ru": `Ethernet 2 Включен`},
			}),
			EthernetIpConnectionEnabled: controls.NewSwitchControl(client, name, "Ethernet IP Connection Enabled", controls.Meta{
				Type: "switch",

				Order:    6,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Ethernet 1 Enabled`, "ru": `Ethernet 1 Включен`},
			}),
			GprsIpConnectionEnabled: controls.NewSwitchControl(client, name, "GPRS IP Connection Enabled", controls.Meta{
				Type: "switch",

				Order:    18,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `GPRS IP Enabled`, "ru": `GPRS IP Включен`},
			}),
			WiFi2IpConnectionEnabled: controls.NewSwitchControl(client, name, "Wi-Fi 2 IP Connection Enabled", controls.Meta{
				Type: "switch",

				Order:    15,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Wi-Fi 2 Enabled`, "ru": `Wi-Fi 2 Включен`},
			}),
			WiFiIpConnectionEnabled: controls.NewSwitchControl(client, name, "Wi-Fi IP Connection Enabled", controls.Meta{
				Type: "switch",

				Order:    12,
				ReadOnly: true,
				Title:    controls.MultilingualText{"en": `Wi-Fi 1 Enabled`, "ru": `Wi-Fi 1 Включен`},
			}),
		}

		instanceNetwork = &Network{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceNetwork
}
