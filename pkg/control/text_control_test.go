package control

import (
	"github.com/ValentinAlekhin/wb-go/internal/mqttmock"
	"github.com/ValentinAlekhin/wb-go/internal/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTextControl_SetAndGetValue(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	mqttmock.AddOnHandler(client)

	meta := Meta{
		Type: "text",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем TextControl
	textControl := NewTextControl(client, device, controlName, meta)

	// Генерируем случайные строки для теста
	newValue := testutils.RandString(10)
	updatedValue := testutils.RandString(10)

	// Устанавливаем значение
	textControl.SetValue(newValue)
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значение установлено
	assert.Equal(t, newValue, textControl.GetValue())

	// Устанавливаем другое значение
	textControl.SetValue(updatedValue)
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что новое значение установлено
	assert.Equal(t, updatedValue, textControl.GetValue())
}

func TestTextControl_AddWatcher(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	mqttmock.AddOnHandler(client)

	meta := Meta{
		Type: "text",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем TextControl
	textControl := NewTextControl(client, device, controlName, meta)

	var newValue, oldValue string

	// Добавляем наблюдателя
	textControl.AddWatcher(func(payload WatcherPayload) {
		newValue = payload.NewValue
		oldValue = payload.OldValue
	})

	// Генерируем случайные строки для теста
	initialValue := testutils.RandString(10)
	updatedValue := testutils.RandString(10)

	// Устанавливаем значение
	textControl.SetValue(initialValue)
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что новый и старый значения обновились
	assert.Equal(t, initialValue, newValue)
	assert.Equal(t, "", oldValue)

	// Устанавливаем новое значение
	textControl.SetValue(updatedValue)
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что новое и старое значения обновились
	assert.Equal(t, updatedValue, newValue)
	assert.Equal(t, initialValue, oldValue)
}
