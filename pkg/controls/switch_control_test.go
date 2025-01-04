package controls

import (
	"github.com/ValentinAlekhin/wb-go/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSwitchControl_SetAndGetValue(t *testing.T) {
	meta := Meta{
		Type: "switch",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем SwitchControl
	switchControl := NewSwitchControl(client, device, controlName, meta)

	// Устанавливаем значение
	switchControl.SetValue(true)
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значение получено корректно
	assert.Equal(t, true, switchControl.GetValue())

	// Устанавливаем другое значение
	switchControl.SetValue(false)
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что новое значение получено корректно
	assert.Equal(t, false, switchControl.GetValue())
}

func TestSwitchControl_Toggle(t *testing.T) {
	meta := Meta{
		Type: "switch",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем SwitchControl
	switchControl := NewSwitchControl(client, device, controlName, meta)

	// Устанавливаем значение в false
	switchControl.SetValue(false)
	time.Sleep(50 * time.Millisecond)

	// Переключаем в true
	switchControl.Toggle()
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значение стало true
	assert.Equal(t, true, switchControl.GetValue())

	// Переключаем в false
	switchControl.Toggle()
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значение снова стало false
	assert.Equal(t, false, switchControl.GetValue())
}

func TestSwitchControl_AddWatcher(t *testing.T) {
	meta := Meta{
		Type: "switch",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем SwitchControl
	switchControl := NewSwitchControl(client, device, controlName, meta)

	var newValue, oldValue bool

	switchControl.AddWatcher(func(payload SwitchControlWatcherPayload) {
		newValue = payload.NewValue
		oldValue = payload.OldValue
	})

	// Устанавливаем значение в true
	switchControl.SetValue(true)
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значения обновились корректно
	assert.Equal(t, true, newValue)
	assert.Equal(t, false, oldValue)

	// Устанавливаем значение в false
	switchControl.SetValue(false)
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значения обновились корректно
	assert.Equal(t, false, newValue)
	assert.Equal(t, true, oldValue)
}

func TestSwitchControl_TurnOn(t *testing.T) {
	meta := Meta{
		Type: "switch",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем SwitchControl
	switchControl := NewSwitchControl(client, device, controlName, meta)

	// Устанавливаем значение в false
	switchControl.SetValue(false)
	time.Sleep(50 * time.Millisecond)

	// Включаем
	switchControl.TurnOn()
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значение стало true
	assert.Equal(t, true, switchControl.GetValue())
}

func TestSwitchControl_TurnOff(t *testing.T) {
	meta := Meta{
		Type: "switch",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем SwitchControl
	switchControl := NewSwitchControl(client, device, controlName, meta)

	// Устанавливаем значение в true
	switchControl.SetValue(true)
	time.Sleep(50 * time.Millisecond)

	// Выключаем
	switchControl.TurnOff()
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значение стало false
	assert.Equal(t, false, switchControl.GetValue())
}

func TestSwitchControl_TurnOnAndOff(t *testing.T) {
	meta := Meta{
		Type: "switch",
	}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем SwitchControl
	switchControl := NewSwitchControl(client, device, controlName, meta)

	// Устанавливаем значение в false
	switchControl.SetValue(false)
	time.Sleep(50 * time.Millisecond)

	// Включаем
	switchControl.TurnOn()
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, true, switchControl.GetValue())

	// Выключаем
	switchControl.TurnOff()
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, false, switchControl.GetValue())
}
