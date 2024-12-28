package main

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/examples/devices"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/homeassistant"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Подключение к брокеру
	opt := wb.Options{
		Broker:   "192.168.1.150:1883",
		ClientId: "client-wb-go-test",
	}
	client := wb.NewClient(opt)

	// Создание устройств
	WbMswV4151 := devices.NewWbMswV4151(client)
	WbMr6Cu145 := devices.NewWbMr6Cu145(client)
	rgbLed := devices.NewWbLed106(client)
	cctLed := devices.NewWbLed150(client)
	wbMdm := devices.NewWbMdm381(client)

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
	discovery.AddDevice(cctLed.GetInfo())
	discovery.AddDevice(wbMdm.GetInfo())

	// Добавление скрипта
	WbMswV4151.Controls.CurrentMotion.AddWatcher(func(payload controls.ValueControlWatcherPayload) {
		fmt.Printf("Получено новое сообщение: %f\n", payload.NewValue)

		if payload.NewValue > 100 {
			WbMr6Cu145.Controls.K1.TurnOff()
		} else {
			WbMr6Cu145.Controls.K1.TurnOn()
		}
	})

	<-stop

	// Отключениие от брокера и завершение программы
	client.Disconnect(500)

	fmt.Println("Программа завершена!")
}
