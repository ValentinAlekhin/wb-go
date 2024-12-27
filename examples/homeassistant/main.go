package main

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/examples/devices"
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
	wbMsw := devices.NewWbMswV4151(client)
	rgbLed := devices.NewWbLed106(client)
	cctLed := devices.NewWbLed150(client)

	// Создание конфигурации Home Assistant
	discovery := homeassistant.NewDiscovery("homeassistant", client)

	// Добавление устройств в Home Assistant
	discovery.AddDevice(wbMsw.GetInfo())
	discovery.AddDevice(rgbLed.GetInfo())
	discovery.AddDevice(cctLed.GetInfo())

	<-stop

	// Отключениие от брокера и завершение программы
	client.Disconnect(500)

	fmt.Println("Программа завершена!")
}
