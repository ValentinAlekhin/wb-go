package virualcontrol

import (
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestVirtualTimeControlGetValue(t *testing.T) {
	controlName := testutils.RandString(10)
	defaultValue := time.Now()

	opt := TimeOptions{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: testClient,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualTimeControl(opt)

	// Проверяем, что значение по умолчанию корректно возвращается
	assert.Equal(t, defaultValue.Format("15:04:05"), vc.GetValue().Format("15:04:05"))
}

func TestVirtualTimeControlSetValue(t *testing.T) {
	controlName := testutils.RandString(10)
	defaultValue := time.Now()

	opt := TimeOptions{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: testClient,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualTimeControl(opt)

	// Устанавливаем новое значение
	newValue := time.Date(2025, 1, 1, 10, 30, 45, 0, time.UTC)
	vc.SetValue(newValue)

	// Проверяем, что значение установилось
	assert.Equal(t, newValue.Format("15:04:05"), vc.GetValue().Format("15:04:05"))
}

func TestVirtualTimeControlOnHandler(t *testing.T) {
	controlName := testutils.RandString(10)
	defaultValue := time.Now()

	handlerCalled := false

	opt := TimeOptions{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: testClient,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
		OnHandler: func(payload OnTimeHandlerPayload) {
			handlerCalled = true
			assert.Equal(t, "10:30:45", payload.Value.Format("15:04:05")) // Проверяем, что передано правильное значение
		},
	}

	vc := NewVirtualTimeControl(opt)

	err := testClient.Publish(wb.PublishPayload{
		Value: "10:30:45",
		QOS:   0,
		Topic: vc.GetInfo().CommandTopic,
	})
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Проверяем, что обработчик был вызван
	assert.True(t, handlerCalled)
}

func TestVirtualTimeControlAddWatcher(t *testing.T) {
	controlName := testutils.RandString(10)
	defaultValue := time.Now()

	vc := NewVirtualTimeControl(TimeOptions{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: testClient,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	})

	var watcherCalled bool
	// Добавляем watcher для контроля изменений
	vc.AddWatcher(func(payload TimeControlWatcherPayload) {
		watcherCalled = true
		assert.Equal(t, "10:30:45", payload.NewValue.Format("15:04:05"))                      // Проверяем, что новое значение корректное
		assert.Equal(t, defaultValue.Format("15:04:05"), payload.OldValue.Format("15:04:05")) // Проверяем, что старое значение корректное
	})

	// Устанавливаем новое значение, что должно вызвать срабатывание watcher
	newValue := time.Date(2025, 1, 1, 10, 30, 45, 0, time.UTC)
	vc.SetValue(newValue)

	// Проверяем, что watcher был вызван
	assert.True(t, watcherCalled)
}

func TestVirtualTimeControlMetaType(t *testing.T) {
	controlName := testutils.RandString(10)
	defaultValue := time.Now()

	opt := TimeOptions{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: testClient,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualTimeControl(opt)

	// Проверяем, что Meta.Type корректно установлен в "text"
	assert.Equal(t, "text", vc.control.meta.Type)
}
