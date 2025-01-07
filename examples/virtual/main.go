package main

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/examples/device"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/pkg/virtuladevice"
	"log"
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

	wbMsw := device.NewWbMswV4151(client)

	thermostat, err := virtuladevice.NewThermostat(virtuladevice.ThermostatConfig{
		Client:              client,
		Device:              "thermostat",
		TargetTemperature:   21,
		Hysteresis:          1.5,
		TemperatureControls: []*control.ValueControl{wbMsw.Controls.Temperature},
	})

	if err != nil {
		log.Fatal(err)
	}

	thermostat.AddHeaterWatcher(func(p control.SwitchControlWatcherPayload) {
		fmt.Println("Heater: ", p.NewValue)
	})

	_, err = virtuladevice.NewAdaptiveLight(virtuladevice.AdaptiveLightConfig{
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
