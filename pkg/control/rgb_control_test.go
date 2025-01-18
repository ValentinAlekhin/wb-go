package control

import (
	"github.com/ValentinAlekhin/wb-go/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRgbControl_SetAndGetValue(t *testing.T) {
	t.Parallel()

	client, server, destroy := testutils.GetClientWithBroker()
	testutils.AddOnHandler(server)
	defer destroy()

	meta := Meta{
		Type: "rgb",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем RgbControl
	rgbControl := NewRgbControl(client, device, controlName, meta)

	// Устанавливаем значение
	rgbControl.SetValue(RgbValue{Red: 255, Green: 0, Blue: 0})
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значение получено корректно
	assert.Equal(t, RgbValue{Red: 255, Green: 0, Blue: 0}, rgbControl.GetValue())

	// Устанавливаем другое значение
	rgbControl.SetValue(RgbValue{Red: 0, Green: 255, Blue: 0})
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что новое значение получено корректно
	assert.Equal(t, RgbValue{Red: 0, Green: 255, Blue: 0}, rgbControl.GetValue())

	// Устанавливаем другое значение
	rgbControl.SetValue(RgbValue{Red: 0, Green: 0, Blue: 255})
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что новое значение получено корректно
	assert.Equal(t, RgbValue{Red: 0, Green: 0, Blue: 255}, rgbControl.GetValue())
}

func TestRgbControl_AddWatcher(t *testing.T) {
	t.Parallel()

	client, server, destroy := testutils.GetClientWithBroker()
	testutils.AddOnHandler(server)
	defer destroy()

	meta := Meta{
		Type: "rgb",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем RgbControl
	rgbControl := NewRgbControl(client, device, controlName, meta)

	var newValue, oldValue RgbValue

	rgbControl.AddWatcher(func(payload RgbControlWatcherPayload) {
		newValue = payload.NewValue
		oldValue = payload.OldValue
	})

	// Устанавливаем значение
	rgbControl.SetValue(RgbValue{Red: 255, Green: 0, Blue: 0})
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значения обновились корректно
	assert.Equal(t, RgbValue{Red: 255, Green: 0, Blue: 0}, newValue)
	assert.Equal(t, RgbValue{Red: 0, Green: 0, Blue: 0}, oldValue) // начальное значение должно быть 0

	// Устанавливаем другое значение
	rgbControl.SetValue(RgbValue{Red: 0, Green: 255, Blue: 0})
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значения обновились корректно
	assert.Equal(t, RgbValue{Red: 0, Green: 255, Blue: 0}, newValue)
	assert.Equal(t, RgbValue{Red: 255, Green: 0, Blue: 0}, oldValue)

	// Устанавливаем другое значение
	rgbControl.SetValue(RgbValue{Red: 0, Green: 0, Blue: 255})
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значения обновились корректно
	assert.Equal(t, RgbValue{Red: 0, Green: 0, Blue: 255}, newValue)
	assert.Equal(t, RgbValue{Red: 0, Green: 255, Blue: 0}, oldValue)
}
