package mqttmock

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strings"
	"sync"
)

// MockClient симулирует работу MQTT клиента
type MockClient struct {
	subscriptions map[string]mqtt.MessageHandler // Хранение подписок
	messages      map[string]mqtt.Message        // Хранение Retained сообщений
	mu            sync.RWMutex                   // Мьютекс для конкурентного доступа
}

// NewMockClient создаёт новый MockClient
func NewMockClient() *MockClient {
	return &MockClient{
		subscriptions: make(map[string]mqtt.MessageHandler),
		messages:      make(map[string]mqtt.Message),
	}
}

// Subscribe подписывается на тему с обработчиком сообщений
func (m *MockClient) Subscribe(topic string, handler mqtt.MessageHandler) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscriptions[topic] = handler

	// Если есть Retained сообщение для этой темы, доставляем его
	for t, msg := range m.messages {
		if matchTopic(t, topic) {
			handler(nil, msg)
		}
	}
	return nil
}

// Publish публикует сообщение в указанную тему
func (m *MockClient) Publish(p wb.PublishPayload) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	msg := &MockMessage{
		topic:    p.Topic,
		payload:  []byte(p.Value),
		qos:      p.QOS,
		retained: p.Retained,
	}

	// Сохраняем Retained сообщение, если требуется
	if p.Retained {
		m.messages[p.Topic] = msg
	}

	// Доставляем сообщение всем подписчикам, темы которых совпадают с опубликованной
	for subTopic, handler := range m.subscriptions {
		if matchTopic(p.Topic, subTopic) {
			handler(nil, msg)
		}
	}
	return nil
}

// Unsubscribe отписывается от указанных тем
func (m *MockClient) Unsubscribe(topics ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, topic := range topics {
		delete(m.subscriptions, topic)
	}
	return nil
}

// Disconnect отключает клиента (сбрасывает подписки и сообщения)
func (m *MockClient) Disconnect(_ uint) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscriptions = make(map[string]mqtt.MessageHandler)
	m.messages = make(map[string]mqtt.Message)
}

// GetClient возвращает nil (MockClient не имеет реального клиента)
func (m *MockClient) GetClient() mqtt.Client {
	return nil
}

func (m *MockClient) GetPublishedMessages() []mqtt.Message {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var messages []mqtt.Message
	for _, msg := range m.messages {
		messages = append(messages, msg)
	}
	return messages
}

// matchTopic проверяет, совпадает ли опубликованная тема с подпиской
func matchTopic(publishTopic, subscriptionTopic string) bool {
	publishLevels := strings.Split(publishTopic, "/")
	subscriptionLevels := strings.Split(subscriptionTopic, "/")

	for i, subLevel := range subscriptionLevels {
		if subLevel == "#" {
			return true // # соответствует всем оставшимся уровням
		}
		if subLevel == "+" {
			continue // + соответствует любому одноуровневому элементу
		}
		if i >= len(publishLevels) || subLevel != publishLevels[i] {
			return false
		}
	}
	return len(publishLevels) == len(subscriptionLevels)
}

func AddOnHandler(client wb.ClientInterface) {
	var handler mqtt.MessageHandler = func(_ mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		err := client.Publish(wb.PublishPayload{
			Topic:    strings.Replace(topic, "/on", "", 1),
			Value:    string(msg.Payload()),
			QOS:      0,
			Retained: true,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	subTopic := fmt.Sprintf(conventions.CONV_CONTROL_ON_VALUE_FMT, "+", "+")
	_ = client.Subscribe(subTopic, handler)
}

// MockMessage имитирует сообщение MQTT и реализует интерфейс Message
type MockMessage struct {
	topic     string
	payload   []byte
	duplicate bool
	qos       byte
	retained  bool
	messageID uint16
	acked     bool // для отслеживания вызова Ack()
}

// Duplicate возвращает флаг дубликата сообщения
func (m *MockMessage) Duplicate() bool {
	return m.duplicate
}

// Qos возвращает уровень QoS сообщения
func (m *MockMessage) Qos() byte {
	return m.qos
}

// Retained возвращает флаг Retained сообщения
func (m *MockMessage) Retained() bool {
	return m.retained
}

// Topic возвращает тему сообщения
func (m *MockMessage) Topic() string {
	return m.topic
}

// MessageID возвращает идентификатор сообщения
func (m *MockMessage) MessageID() uint16 {
	return m.messageID
}

// Payload возвращает содержимое сообщения
func (m *MockMessage) Payload() []byte {
	return m.payload
}

// Ack симулирует подтверждение сообщения
func (m *MockMessage) Ack() {
	m.acked = true
}

// WasAcked проверяет, был ли вызван Ack()
func (m *MockMessage) WasAcked() bool {
	return m.acked
}
