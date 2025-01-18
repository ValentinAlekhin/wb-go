package test_mqtt_server

import (
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"log"
	"time"

	mochi "github.com/mochi-mqtt/server/v2"
)

// StartMQTTBroker запускает локальный Mochi-MQTT брокер
func StartMQTTBroker() *mochi.Server {
	server := mochi.New(&mochi.Options{
		InlineClient: true,
	})

	_ = server.AddHook(new(auth.AllowHook), nil)
	tcp := listeners.NewTCP(listeners.Config{
		ID:      "t1",
		Address: ":1883",
	})
	err := server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	return server
}
