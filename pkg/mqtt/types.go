package mqtt

import mqtt "github.com/eclipse/paho.mqtt.golang"

// ClientInterface defines an interface for MQTT client functionality.
type ClientInterface interface {
	// Subscribe subscribes to an MQTT topic with a given message handler.
	Subscribe(topic string, handler mqtt.MessageHandler) error

	// Publish publishes a message to an MQTT topic with the specified payload.
	Publish(p PublishPayload) error

	// Unsubscribe unsubscribes from one or more MQTT topics.
	Unsubscribe(topics ...string) error

	// Disconnect disconnects from the MQTT broker with the specified quiesce timeout.
	Disconnect(quiesce uint)

	// GetClient returns the underlying MQTT client.
	GetClient() mqtt.Client
}

// Client is a structure for managing the MQTT connection
type Client struct {
	client  mqtt.Client
	options Options
}

type Options struct {
	Broker           string
	ClientId         string
	Username         string
	Password         string
	QoS              byte
	OnConnectionLost func(err error)
}

type PublishPayload struct {
	Topic    string
	Value    string
	QOS      byte
	Retained bool
}
