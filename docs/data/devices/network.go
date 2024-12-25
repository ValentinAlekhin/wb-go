package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
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
	Name     string
	Address  string
	Controls *NetworkControls
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
		deviceName := fmt.Sprintf("%s_%s", "network", "")
		controlList := &NetworkControls{
			ActiveConnections:            controls.NewTextControl(client, deviceName, "Active Connections"),
			DefaultInterface:             controls.NewTextControl(client, deviceName, "Default Interface"),
			Ethernet2Ip:                  controls.NewTextControl(client, deviceName, "Ethernet 2 IP"),
			Ethernet2IpOnlineStatus:      controls.NewSwitchControl(client, deviceName, "Ethernet 2 IP Online Status"),
			EthernetIp:                   controls.NewTextControl(client, deviceName, "Ethernet IP"),
			EthernetIpOnlineStatus:       controls.NewSwitchControl(client, deviceName, "Ethernet IP Online Status"),
			GprsIp:                       controls.NewTextControl(client, deviceName, "GPRS IP"),
			GprsIpOnlineStatus:           controls.NewSwitchControl(client, deviceName, "GPRS IP Online Status"),
			InternetConnection:           controls.NewTextControl(client, deviceName, "Internet Connection"),
			WiFi2Ip:                      controls.NewTextControl(client, deviceName, "Wi-Fi 2 IP"),
			WiFi2IpOnlineStatus:          controls.NewSwitchControl(client, deviceName, "Wi-Fi 2 IP Online Status"),
			WiFiIp:                       controls.NewTextControl(client, deviceName, "Wi-Fi IP"),
			WiFiIpOnlineStatus:           controls.NewSwitchControl(client, deviceName, "Wi-Fi IP Online Status"),
			Ethernet2IpConnectionEnabled: controls.NewSwitchControl(client, deviceName, "Ethernet 2 IP Connection Enabled"),
			EthernetIpConnectionEnabled:  controls.NewSwitchControl(client, deviceName, "Ethernet IP Connection Enabled"),
			GprsIpConnectionEnabled:      controls.NewSwitchControl(client, deviceName, "GPRS IP Connection Enabled"),
			WiFi2IpConnectionEnabled:     controls.NewSwitchControl(client, deviceName, "Wi-Fi 2 IP Connection Enabled"),
			WiFiIpConnectionEnabled:      controls.NewSwitchControl(client, deviceName, "Wi-Fi IP Connection Enabled"),
		}

		instanceNetwork = &Network{
			Name:     deviceName,
			Address:  "",
			Controls: controlList,
		}
	})

	return instanceNetwork
}
