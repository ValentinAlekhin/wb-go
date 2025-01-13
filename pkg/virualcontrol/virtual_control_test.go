package virualcontrol

import (
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/testutils"
	mock2 "github.com/ValentinAlekhin/wb-go/testutils/mock"
	"github.com/ValentinAlekhin/wb-go/testutils/test_mqtt_server"
	//mqtt "github.com/eclipse/paho.mqtt.golang"
	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var broker *mochi.Server
var client wb.ClientInterface

func TestMain(m *testing.M) {
	broker = test_mqtt_server.StartMQTTBroker(true)
	client = testutils.GetMqttClient()

	m.Run()

	client.Disconnect(100)
	err := broker.Close()
	if err != nil {
		return
	}
}

func TestVirtualControl_SetAndGetValue(t *testing.T) {
	meta := control.Meta{}

	vc := NewVirtualControl(client, "testDevice", "testControl", meta, nil)

	// Set a value
	vc.SetValue("testValue")

	// Verify the value
	assert.Equal(t, "testValue", vc.GetValue())
}

func TestVirtualControl_EventChan(t *testing.T) {
	meta := control.Meta{}

	vc := NewVirtualControl(client, "testDevice", "testControl", meta, nil)

	eventReceived := make(chan bool)

	// Add a watcher
	vc.AddWatcher(func(payload control.WatcherPayload) {
		assert.Equal(t, "newValue", payload.NewValue)
		assert.Equal(t, "oldValue", payload.OldValue)
		eventReceived <- true
	})

	// Set the old value
	vc.SetValue("oldValue")

	// Set a new value
	vc.SetValue("newValue")

	// Wait for the watcher to trigger
	select {
	case <-eventReceived:
	case <-time.After(1 * time.Second):
		t.Fatal("Watcher did not trigger in time")
	}
}

//func TestVirtualControl_SubscribeToOnTopic(t *testing.T) {
//	mClient := new(mock2.MqttClientMock)
//	meta := control.Meta{}
//
//	vc := NewVirtualControl(mClient, "testDevice", "testControl", meta, nil)
//
//	mClient.On("Subscribe", vc.GetInfo().CommandTopic, mock.Anything).Return(nil)
//
//	// Subscribe to the "on" topic
//	vc.SetValue("initialValue")
//
//	// Simulate receiving a message
//	handler := mClient.Calls[0].Arguments.Get(1).(mqtt.MessageHandler)
//	handler(nil, mqtt.Message{Payload: []byte("newValue")})
//
//	// Verify that the value is updated
//	assert.Equal(t, "newValue", vc.GetValue())
//}

func TestVirtualControl_SetMeta(t *testing.T) {
	mClient := new(mock2.MqttClientMock)
	meta := control.Meta{}

	vc := NewVirtualControl(mClient, "testDevice", "testControl", meta, nil)

	mClient.On("Publish", mock.Anything).Return(nil)

	// Check if meta is set correctly
	assert.NotNil(t, vc.GetInfo().Meta)
}
