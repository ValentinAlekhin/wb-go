package mqttmock

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMockMessage(t *testing.T) {
	t.Parallel()

	msg := &MockMessage{
		topic:     "test/topic",
		payload:   []byte("test payload"),
		duplicate: false,
		qos:       1,
		retained:  true,
		messageID: 1234,
	}

	assert.Equal(t, "test/topic", msg.Topic(), "Topic should match")
	assert.Equal(t, []byte("test payload"), msg.Payload(), "Payload should match")
	assert.True(t, msg.Retained(), "Message should be retained")
	assert.Equal(t, byte(1), msg.Qos(), "QoS should match")
	assert.Equal(t, uint16(1234), msg.MessageID(), "MessageID should match")
	assert.False(t, msg.Duplicate(), "Duplicate should be false by default")

	msg.Ack()
	assert.True(t, msg.WasAcked(), "Ack should have been called")
}

func TestMockClient(t *testing.T) {
	t.Parallel()

	client := NewMockClient()

	receivedMessages := []mqtt.Message{}

	// Подписываемся на тему
	err := client.Subscribe("test/topic", func(_ mqtt.Client, msg mqtt.Message) {
		receivedMessages = append(receivedMessages, msg)
	})
	assert.NoError(t, err, "Subscribe should not return an error")

	// Публикуем сообщение
	err = client.Publish(wb.PublishPayload{
		Topic:    "test/topic",
		Value:    "Hello, world!",
		QOS:      1,
		Retained: false,
	})
	assert.NoError(t, err, "Publish should not return an error")

	// Проверяем, что сообщение получено подписчиком
	assert.Len(t, receivedMessages, 1, "Subscriber should receive one message")

	msg := receivedMessages[0]

	assert.Equal(t, "test/topic", msg.Topic(), "Topic should match")
	assert.Equal(t, "Hello, world!", string(msg.Payload()), "Payload should match")
	assert.False(t, msg.Retained(), "Message should not be retained")

	// Тестируем вызов Ack()
	mockMsg, ok := msg.(*MockMessage)
	assert.True(t, ok, "Message should be of type MockMessage")

	mockMsg.Ack()
	assert.True(t, mockMsg.WasAcked(), "Message should have been acknowledged")
}

func TestMockClientRetained(t *testing.T) {
	t.Parallel()

	client := NewMockClient()

	// Публикуем сообщение с флагом Retained
	err := client.Publish(wb.PublishPayload{
		Topic:    "test/retained",
		Value:    "Retained message",
		QOS:      1,
		Retained: true,
	})
	assert.NoError(t, err, "Publish should not return an error")

	receivedMessages := []mqtt.Message{}

	// Новый подписчик должен получить Retained сообщение
	err = client.Subscribe("test/retained", func(_ mqtt.Client, msg mqtt.Message) {
		receivedMessages = append(receivedMessages, msg)
	})
	assert.NoError(t, err, "Subscribe should not return an error")

	// Проверяем, что Retained сообщение доставлено
	assert.Len(t, receivedMessages, 1, "Subscriber should receive one retained message")

	msg := receivedMessages[0]
	assert.Equal(t, "Retained message", string(msg.Payload()), "Retained message payload should match")
	assert.True(t, msg.Retained(), "Message should be retained")
}

func TestMockClientUnsubscribe(t *testing.T) {
	t.Parallel()

	client := NewMockClient()

	receivedMessages := []mqtt.Message{}

	// Подписываемся на две темы
	err := client.Subscribe("test/topic1", func(_ mqtt.Client, msg mqtt.Message) {
		receivedMessages = append(receivedMessages, msg)
	})
	assert.NoError(t, err, "Subscribe to topic1 should not return an error")

	err = client.Subscribe("test/topic2", func(_ mqtt.Client, msg mqtt.Message) {
		receivedMessages = append(receivedMessages, msg)
	})
	assert.NoError(t, err, "Subscribe to topic2 should not return an error")

	// Проверяем, что подписки существуют
	err = client.Publish(wb.PublishPayload{
		Topic:    "test/topic1",
		Value:    "Message for topic1",
		QOS:      1,
		Retained: false,
	})
	assert.NoError(t, err, "Publish to topic1 should not return an error")
	assert.Len(t, receivedMessages, 1, "Subscriber should receive message on topic1")

	// Отписываемся от темы `test/topic1`
	err = client.Unsubscribe("test/topic1")
	assert.NoError(t, err, "Unsubscribe from topic1 should not return an error")

	// Публикуем ещё одно сообщение на `test/topic1`
	err = client.Publish(wb.PublishPayload{
		Topic:    "test/topic1",
		Value:    "Another message for topic1",
		QOS:      1,
		Retained: false,
	})
	assert.NoError(t, err, "Publish to topic1 should not return an error")
	assert.Len(t, receivedMessages, 1, "Subscriber should not receive messages on unsubscribed topic")

	// Проверяем, что подписка на `test/topic2` всё ещё работает
	err = client.Publish(wb.PublishPayload{
		Topic:    "test/topic2",
		Value:    "Message for topic2",
		QOS:      1,
		Retained: false,
	})
	assert.NoError(t, err, "Publish to topic2 should not return an error")
	assert.Len(t, receivedMessages, 2, "Subscriber should receive messages on topic2")
}

func TestMockClientDisconnect(t *testing.T) {
	t.Parallel()

	client := NewMockClient()

	receivedMessages := []mqtt.Message{}

	// Подписываемся на тему
	err := client.Subscribe("test/topic", func(_ mqtt.Client, msg mqtt.Message) {
		receivedMessages = append(receivedMessages, msg)
	})
	assert.NoError(t, err, "Subscribe should not return an error")

	// Публикуем сообщение
	err = client.Publish(wb.PublishPayload{
		Topic:    "test/topic",
		Value:    "Message before disconnect",
		QOS:      1,
		Retained: false,
	})
	assert.NoError(t, err, "Publish before disconnect should not return an error")
	assert.Len(t, receivedMessages, 1, "Subscriber should receive message before disconnect")

	// Вызываем Disconnect
	client.Disconnect(0)

	// Пытаемся публиковать сообщение после Disconnect
	err = client.Publish(wb.PublishPayload{
		Topic:    "test/topic",
		Value:    "Message after disconnect",
		QOS:      1,
		Retained: false,
	})
	assert.NoError(t, err, "Publish after disconnect should not return an error")

	// Проверяем, что сообщение не доставлено
	assert.Len(t, receivedMessages, 1, "Subscriber should not receive messages after disconnect")
}

func TestMockClientGetClient(t *testing.T) {
	t.Parallel()

	client := NewMockClient()

	// GetClient всегда возвращает nil для MockClient
	actualClient := client.GetClient()
	assert.Nil(t, actualClient, "MockClient.GetClient should return nil")
}

func TestMockClientWildcardSubscription(t *testing.T) {
	t.Parallel()

	client := NewMockClient()

	receivedMessages := []mqtt.Message{}

	// Подписываемся с использованием +
	err := client.Subscribe("test/+/sub", func(_ mqtt.Client, msg mqtt.Message) {
		receivedMessages = append(receivedMessages, msg)
	})
	assert.NoError(t, err, "Subscribe with '+' should not return an error")

	// Публикуем сообщение
	err = client.Publish(wb.PublishPayload{
		Topic:    "test/foo/sub",
		Value:    "Message with +",
		QOS:      1,
		Retained: false,
	})
	assert.NoError(t, err, "Publish to matching topic should not return an error")

	assert.Len(t, receivedMessages, 1, "Subscriber should receive message for topic matching '+'")
	assert.Equal(t, "Message with +", string(receivedMessages[0].Payload()), "Payload should match")

	// Подписываемся с использованием #
	err = client.Subscribe("test/#", func(_ mqtt.Client, msg mqtt.Message) {
		receivedMessages = append(receivedMessages, msg)
	})
	assert.NoError(t, err, "Subscribe with '#' should not return an error")

	// Публикуем сообщение
	err = client.Publish(wb.PublishPayload{
		Topic:    "test/foo/bar",
		Value:    "Message with #",
		QOS:      1,
		Retained: false,
	})
	assert.NoError(t, err, "Publish to matching topic should not return an error")

	assert.Len(t, receivedMessages, 2, "Subscriber should receive message for topic matching '#'")
	assert.Equal(t, "Message with #", string(receivedMessages[1].Payload()), "Payload should match")
}

// AddOnHandler проверяет преобразование и публикацию сообщений
func TestAddOnHandler(t *testing.T) {
	t.Parallel()

	client := NewMockClient()

	// Вызываем AddOnHandler для мока
	AddOnHandler(client)

	// Публикуем сообщение в топик, соответствующий шаблону
	deviceName := "device1"
	controlName := "control1"
	err := client.Publish(wb.PublishPayload{
		Topic:    fmt.Sprintf(conventions.CONV_CONTROL_ON_VALUE_FMT, deviceName, controlName),
		Value:    "test-value",
		QOS:      0,
		Retained: false,
	})
	assert.NoError(t, err, "Publish to /on topic should not return an error")

	// Проверяем, что сообщение было опубликовано с корректным преобразованием топика
	expectedTopic := fmt.Sprintf(conventions.CONV_CONTROL_VALUE_FMT, deviceName, controlName)
	publishedMessages := client.GetPublishedMessages()

	assert.Len(t, publishedMessages, 1, "There should be one published message")
	assert.Equal(t, expectedTopic, publishedMessages[0].Topic(), "Topic should be transformed correctly")
	assert.Equal(t, "test-value", string(publishedMessages[0].Payload()), "Payload should match")
	assert.Equal(t, byte(0), publishedMessages[0].Qos(), "QoS should match")
	assert.True(t, publishedMessages[0].Retained(), "Retained flag should be true")
}
