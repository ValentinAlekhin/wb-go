package controls

import (
	"github.com/ValentinAlekhin/wb-go/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRangeControl_SetAndGetValue(t *testing.T) {
	meta := Meta{}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем RangeControl
	rangeControl := NewRangeControl(client, device, controlName, meta)

	// Устанавливаем значение
	rangeControl.SetValue(42)
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значение получено корректно
	assert.Equal(t, 42, rangeControl.GetValue())

	// Устанавливаем другое значение
	rangeControl.SetValue(10)
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что новое значение получено корректно
	assert.Equal(t, 10, rangeControl.GetValue())
}

func TestRangeControl_AddWatcher(t *testing.T) {
	meta := Meta{}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем RangeControl
	rangeControl := NewRangeControl(client, device, controlName, meta)

	var newValue, oldValue int

	rangeControl.AddWatcher(func(payload RangeControlWatcherPayload) {
		newValue = payload.NewValue
		oldValue = payload.OldValue
	})

	// Устанавливаем значение
	rangeControl.SetValue(100)
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значения обновились корректно
	assert.Equal(t, 100, newValue)
	assert.Equal(t, 0, oldValue) // начальное значение должно быть 0

	// Устанавливаем другое значение
	rangeControl.SetValue(200)
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значения обновились корректно
	assert.Equal(t, 200, newValue)
	assert.Equal(t, 100, oldValue)
}
