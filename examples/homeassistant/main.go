package main

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/examples/device"
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/homeassistant"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"log"
)

func main() {
	// Подключение к брокеру
	opt := wb.Options{
		Broker:   "192.168.1.150:1883",
		ClientId: "client-wb-go-test",
	}
	client, err := wb.NewClient(opt)
	if err != nil {
		log.Fatal(err)
	}

	// Создание устройств
	WbMswV4151 := device.NewWbMswV4151(client)
	rgbLed := device.NewWbLed106(client)
	cctLed := device.NewWbLed150(client)

	// Создание конфигурации Home Assistant
	discoveryOpt := homeassistant.DiscoveryOptions{
		Client: client,
		Prefix: "homeassistant",
	}
	discovery := homeassistant.NewDiscovery(discoveryOpt)

	// Удаление всех устройств из Home Assistant, добавленных с помощью wb-go.
	// Нужно вызывать, если есть необходимость удалить из Home Assistant ранее добавленные устройства,
	// но которые уже не используются
	discovery.Clear()

	// Добавление устройств в Home Assistant
	discovery.AddDevice(WbMswV4151.GetInfo())
	discovery.AddDevice(rgbLed.GetInfo())

	// Создаем карту, где ключ - название контрола, а значение - имя в Home Assistant
	nameMap := map[string]string{"CCT1": "Свет в гостиной", "CCT2": "Свет в спальне"}
	var configMiddleware homeassistant.ConfigMiddleware = func(config *homeassistant.MqttDiscoveryConfig, device basedevice.Info, control control.Info) {
		name, ok := nameMap[control.Name]
		if !ok {
			return
		}

		config.Name = name
	}
	// Добавляем устройство с использованием промежуточного обработчика конфигурации
	discovery.AddDeviceWithMiddleware(cctLed.GetInfo(), configMiddleware)

	// Отключениие от брокера и завершение программы
	client.Disconnect(500)

	fmt.Println("Программа завершена!")
}
