package testutils

import (
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"log"
)

func GetMqttClient() wb.ClientInterface {
	options := wb.Options{
		Broker:   "localhost:1883", // Используем URL брокера, который запустили
		ClientId: "test-client",
		QoS:      1,
	}

	client, err := wb.NewClient(options)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
