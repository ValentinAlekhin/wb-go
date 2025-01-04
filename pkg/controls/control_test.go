package controls

import (
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/testutils"
	"github.com/ValentinAlekhin/wb-go/testutils/test_mqtt_server"
	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var broker *mochi.Server
var client *wb.Client

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

func TestSetAndGetValue(t *testing.T) {
	meta := Meta{
		Type: "switch",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	control := NewControl(client, device, controlName, meta)

	control.SetValue("on")
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, "on", control.GetValue())

	control.SetValue("off")
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, "off", control.GetValue())
}

func TestControl_AddWatcher(t *testing.T) {
	meta := Meta{
		Type: "switch",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	control := NewControl(client, device, controlName, meta)

	var newValue, oldValue string

	control.AddWatcher(func(payload ControlWatcherPayload) {
		newValue = payload.NewValue
		oldValue = payload.OldValue
	})

	control.SetValue("on")
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, "on", newValue)
	assert.Equal(t, "", oldValue)

	control.SetValue("off")
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, "off", newValue)
	assert.Equal(t, "on", oldValue)
}
