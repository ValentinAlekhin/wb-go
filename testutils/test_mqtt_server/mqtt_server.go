package test_mqtt_server

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"log"
	"strings"
	"time"

	mochi "github.com/mochi-mqtt/server/v2"
)

// StartMQTTBroker запускает локальный Mochi-MQTT брокер
func StartMQTTBroker(deviceOnMock bool) *mochi.Server {
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

	if deviceOnMock {
		callbackFn := func(cl *mochi.Client, sub packets.Subscription, pk packets.Packet) {
			server.Log.Info("inline client received message from subscription", "client", cl.ID, "subscriptionId", sub.Identifier, "topic", pk.TopicName, "payload", string(pk.Payload))
			topic := strings.Replace(pk.TopicName, "/on", "", 1)
			_ = server.Publish(topic, pk.Payload, true, 1)
		}
		server.Log.Info("inline client subscribing")
		subTopic := fmt.Sprintf(conventions.CONV_CONTROL_ON_VALUE_FMT, "+", "+")
		_ = server.Subscribe(subTopic, 1, callbackFn)
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
