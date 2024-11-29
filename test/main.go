package main

import (
	"fmt"
	"wb-go/pkg/mqtt"
	"wb-go/test/devices"
)

func main() {
	opt := mqtt.Options{
		Broker:   "192.168.1.150:1883",
		ClientId: "wb-g0",
	}
	client := mqtt.NewClient(opt)

	rele := devices.NewWbMr6Cu145(client, "", "")
	WbMswV4151 := devices.NewWbMswV4151(client, "", "")

	WbMswV4151.Controls.CurrentMotion.AddWatcher(func(payload devices.ValueControlWatcherPayload) {
		fmt.Println(payload.NewValue)

		if payload.NewValue > 100 {
			rele.Controls.K1.SetValue(true)
		} else {
			rele.Controls.K1.SetValue(false)
		}
	})

	select {}

}
