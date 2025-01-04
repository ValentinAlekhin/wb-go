package testutils

import (
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

func GetMqttClient() *wb.Client {
	options := wb.Options{
		Broker:   "localhost:1883", // Используем URL брокера, который запустили
		ClientId: "test-client",
		QoS:      1,
	}

	client := wb.NewClient(options)

	return client
}
