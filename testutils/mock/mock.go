package mock

import (
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/mock"
)

// MqttClientMock — мок для интерфейса ClientInterface.
type MqttClientMock struct {
	mock.Mock
}

// Subscribe — мок-реализация метода Subscribe.
func (m *MqttClientMock) Subscribe(topic string, handler mqtt.MessageHandler) error {
	args := m.Called(topic, handler)
	return args.Error(0)
}

// Publish — мок-реализация метода Publish.
func (m *MqttClientMock) Publish(p wb.PublishPayload) error {
	args := m.Called(p)
	return args.Error(0)
}

// Unsubscribe — мок-реализация метода Unsubscribe.
func (m *MqttClientMock) Unsubscribe(topics ...string) error {
	args := m.Called(topics)
	return args.Error(0)
}

// Disconnect — мок-реализация метода Disconnect.
func (m *MqttClientMock) Disconnect(quiesce uint) {
	m.Called(quiesce)
}

// GetClient — мок-реализация метода GetClient.
func (m *MqttClientMock) GetClient() mqtt.Client {
	args := m.Called()
	return args.Get(0).(mqtt.Client)
}
