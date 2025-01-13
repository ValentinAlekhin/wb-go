package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

// NewClient creates a new MQTT client
func NewClient(opt Options) (ClientInterface, error) {
	fmt.Printf("Connecting to broker %s\n", opt.Broker)

	client := &Client{
		options: opt,
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(opt.Broker)
	opts.SetClientID(opt.ClientId)
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("Successfully connected to the MQTT broker!")
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("Connection lost: %v\n", err)
		if opt.OnConnectionLost != nil {
			opt.OnConnectionLost(err)
		}
	}

	if opt.Username != "" {
		opts.SetUsername(opt.Username)
	}

	if opt.Password != "" {
		opts.SetPassword(opt.Password)
	}

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("error connecting to MQTT broker: %v", token.Error())
	}

	client.client = mqttClient
	return client, nil
}

// Subscribe subscribes to an MQTT topic
func (c *Client) Subscribe(topic string, handler mqtt.MessageHandler) error {
	if token := c.client.Subscribe(topic, c.options.QoS, handler); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error subscribing to topic %s: %v", topic, token.Error())
	}
	return nil
}

// Publish publishes a message to an MQTT topic
func (c *Client) Publish(p PublishPayload) error {
	if token := c.client.Publish(p.Topic, p.QOS, p.Retained, p.Value); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error publishing to topic %s: %v", p.Topic, token.Error())
	}
	return nil
}

// Unsubscribe unsubscribes from MQTT topics
func (c *Client) Unsubscribe(topics ...string) error {
	if token := c.client.Unsubscribe(topics...); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error unsubscribing from topics: %v", token.Error())
	}
	return nil
}

// GetClient retrieves the underlying MQTT client
func (c *Client) GetClient() mqtt.Client {
	return c.client
}

// Disconnect disconnects from the MQTT broker
func (c *Client) Disconnect(quiesce uint) {
	c.client.Disconnect(quiesce)
}
