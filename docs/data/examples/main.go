package main

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/docs/data/devices"
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
		ClientId: "client-wb-go",
	}
	client := wb.NewClient(opt)

	//Создание устройств
	WbMswV4151 := devices.NewWbMswV4151(client)
	WbMr6Cu145 := devices.NewWbMr6Cu145(client)

	// Добавление скрипта
	WbMswV4151.Controls.CurrentMotion.AddWatcher(func(payload devices.ValueControlWatcherPayload) {
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
