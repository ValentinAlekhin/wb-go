package mqtt

import (
	"github.com/ValentinAlekhin/wb-go/internal/testutils/test_mqtt_server"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var broker *mochi.Server

const testTopic = "test/topic"
const testValue = "test_value"

func TestMain(m *testing.M) {
	broker = test_mqtt_server.StartMQTTBroker()

	m.Run()

	err := broker.Close()
	if err != nil {
		return
	}
}

func TestClient(t *testing.T) {
	options := Options{
		Broker:   "localhost:1883",
		ClientId: "test-client",
		QoS:      1,
	}

	client, err := NewClient(options)
	if err != nil {
		t.Error(err)
	}

	defer client.Disconnect(100)

	payload := PublishPayload{
		Topic: testTopic,
		Value: testValue,
		QOS:   1,
	}

	var receivedMessage string

	handler := func(client mqtt.Client, msg mqtt.Message) {
		receivedMessage = string(msg.Payload())
	}

	err = client.Subscribe(testTopic, handler)
	if err != nil {
		t.Error(err)
	}
	err = client.Publish(payload)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, testValue, receivedMessage, "Полученное сообщение не совпадает с отправленным")
}
