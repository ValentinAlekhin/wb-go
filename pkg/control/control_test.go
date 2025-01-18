package control

import (
	"github.com/ValentinAlekhin/wb-go/internal/mqttmock"
	"github.com/ValentinAlekhin/wb-go/internal/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSetAndGetValue(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	mqttmock.AddOnHandler(client)

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
	t.Parallel()

	client := mqttmock.NewMockClient()
	mqttmock.AddOnHandler(client)

	meta := Meta{
		Type: "switch",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	control := NewControl(client, device, controlName, meta)

	var newValue, oldValue string

	control.AddWatcher(func(payload WatcherPayload) {
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
