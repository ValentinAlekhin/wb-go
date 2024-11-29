package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type Options struct {
	Broker   string
	ClientId string
	Username string
	Password string
	QoS      byte
}

// Client структура для управления MQTT-соединением
type Client struct {
	client  mqtt.Client
	options Options
}

// NewClient создает новый MQTT клиент
func NewClient(opt Options) *Client {
	fmt.Printf("Подключение к брокеру %s\n", opt.Broker)

	client := &Client{
		options: opt,
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(opt.Broker)
	opts.SetClientID(opt.ClientId)
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("Подключение к MQTT-брокеру успешно!")
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Fatalf("Соединение потеряно: %v\n", err)
	}

	if opt.Username != "" {
		opts.SetUsername(opt.Username)
	}

	if opt.Password != "" {
		opts.SetPassword(opt.Password)
	}

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Ошибка подключения к MQTT брокеру: %v", token.Error())
	}

	client.client = mqttClient

	return client
}

// Subscribe подписывается на топик MQTT
func (c *Client) Subscribe(topic string, handler mqtt.MessageHandler) {
	if token := c.client.Subscribe(topic, c.options.QoS, handler); token.Wait() && token.Error() != nil {
		log.Fatalf("Ошибка подписки на топик %s: %v", topic, token.Error())
	}
}

// Publish публикует сообщение в топик MQTT
func (c *Client) Publish(topic string, value string) {
	if token := c.client.Publish(topic, 0, false, value); token.Wait() && token.Error() != nil {
		log.Printf("Ошибка публикации в топик %s: %v", topic, token.Error())
	}
}

func (c *Client) GetClient() mqtt.Client {
	return c.client
}

func (c *Client) Disconnect(quiesce uint) {
	c.client.Disconnect(quiesce)
}
