package main

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/examples/device"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/pkg/virtualdevice"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Подключение к базе данных
	db, err := gorm.Open(sqlite.Open("./db/test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Подключение к брокеру
	opt := wb.Options{
		Broker:   "192.168.1.150:1883",
		ClientId: "client-wb-go-test",
		OnConnectionLost: func(err error) {
			log.Fatal(err)
		},
	}
	client, err := wb.NewClient(opt)
	if err != nil {
		log.Fatal(err)
	}

	wbMsw := device.NewWbMswV4151(client)

	thermostat, err := virtualdevice.NewThermostat(virtualdevice.ThermostatConfig{
		DB:                  db,
		Client:              client,
		Device:              "thermostat",
		TargetTemperature:   21,
		Hysteresis:          1.5,
		TemperatureControls: []*control.ValueControl{wbMsw.Controls.Temperature},
	})

	if err != nil {
		log.Fatal(err)
	}

	thermostat.Controls.Relay.AddWatcher(func(p control.SwitchControlWatcherPayload) {
		fmt.Println("Relay: ", p.NewValue)
	})

	_, err = virtualdevice.NewAdaptiveLight(virtualdevice.AdaptiveLightConfig{
		DB:     db,
		Client: client,
		Device: "adaptive-light",
	})

	if err != nil {
		log.Fatal(err)
	}

	<-stop

	// Отключениие от брокера и завершение программы
	client.Disconnect(500)

	fmt.Println("Программа завершена!")
}
