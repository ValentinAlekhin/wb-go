package control

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestValueControl_SetAndGetValue(t *testing.T) {
	t.Parallel()

	client, server, destroy := testutils.GetClientWithBroker()
	testutils.AddOnHandler(server)
	defer destroy()

	meta := Meta{
		Type: "value",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем ValueControl
	valueControl := NewValueControl(client, device, controlName, meta)

	// Числа для теста
	newValue := 11.111111
	updatedValue := 12.122222

	// Устанавливаем значение
	valueControl.control.SetValue(fmt.Sprintf("%f", newValue))
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значение установлено
	assert.Equal(t, newValue, valueControl.GetValue())

	// Устанавливаем другое значение
	valueControl.control.SetValue(fmt.Sprintf("%f", updatedValue))
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что новое значение установлено
	assert.Equal(t, updatedValue, valueControl.GetValue())
}

func TestValueControl_AddWatcher(t *testing.T) {
	t.Parallel()

	client, server, destroy := testutils.GetClientWithBroker()
	testutils.AddOnHandler(server)
	defer destroy()

	meta := Meta{
		Type: "value",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем ValueControl
	valueControl := NewValueControl(client, device, controlName, meta)

	var newValue, oldValue float64

	// Добавляем наблюдателя
	valueControl.AddWatcher(func(payload ValueControlWatcherPayload) {
		newValue = payload.NewValue
		oldValue = payload.OldValue
	})

	// Генерируем случайные числа для теста
	initialValue := 11.111111
	updatedValue := 12.122222

	// Устанавливаем значение
	valueControl.control.SetValue(fmt.Sprintf("%f", initialValue))
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что новый и старый значения обновились
	assert.Equal(t, initialValue, newValue)
	assert.Equal(t, 0.0, oldValue)

	// Устанавливаем новое значение
	valueControl.control.SetValue(fmt.Sprintf("%f", updatedValue))
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что новое и старое значения обновились
	assert.Equal(t, updatedValue, newValue)
	assert.Equal(t, initialValue, oldValue)
}
