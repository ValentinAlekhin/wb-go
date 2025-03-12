package control

import (
	"context"
	"testing"
	"time"

	"github.com/ValentinAlekhin/wb-go/internal/mqttmock"
	"github.com/ValentinAlekhin/wb-go/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestSetAndGetValue(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := mqttmock.NewMockClient()
	mqttmock.AddOnHandler(client)

	meta := Meta{Type: "switch"}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	control := NewControl(ctx, client, device, controlName, meta)

	control.SetValue("on")
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, "on", control.GetValue())

	control.SetValue("off")
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, "off", control.GetValue())
}

func TestControl_AddWatcher(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := mqttmock.NewMockClient()
	mqttmock.AddOnHandler(client)

	meta := Meta{Type: "switch"}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	control := NewControl(ctx, client, device, controlName, meta)

	var newValue, oldValue string

	control.AddWatcher(func(payload WatcherPayloadString) {
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

func TestControl_ContextCancellation(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	mqttmock.AddOnHandler(client)

	meta := Meta{
		Type: "switch",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	ctx, cancel := context.WithCancel(context.Background())
	control := NewControl(ctx, client, device, controlName, meta)

	control.SetValue("on")
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, "on", control.GetValue())

	// Отмена контекста
	cancel()
	time.Sleep(50 * time.Millisecond)

	assert.NotPanics(t, func() {
		control.SetValue("off")
	})

	assert.NotPanics(t, func() {
		control.AddWatcher(func(payload WatcherPayloadString) {
		})
	})

	// Проверяем, что значение не изменилось после отмены контекста
	assert.Equal(t, "on", control.GetValue())
}
