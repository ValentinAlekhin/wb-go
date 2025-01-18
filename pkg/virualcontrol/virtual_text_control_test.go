package virualcontrol

import (
	"github.com/ValentinAlekhin/wb-go/internal/dbmock"
	"github.com/ValentinAlekhin/wb-go/internal/mqttmock"
	"github.com/ValentinAlekhin/wb-go/internal/testutils"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestVirtualTextControlGetValue(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := "initial_value"

	opt := TextOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualTextControl(opt)

	// Проверяем, что значение по умолчанию корректно возвращается
	assert.Equal(t, defaultValue, vc.GetValue())
}

func TestVirtualTextControlSetValue(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := "default_value"

	opt := TextOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualTextControl(opt)

	// Устанавливаем новое значение
	newValue := "new_value"
	vc.SetValue(newValue)

	// Проверяем, что значение установилось
	assert.Equal(t, newValue, vc.GetValue())
}

func TestVirtualTextControlOnHandler(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := "default_value"

	handlerCalled := false

	opt := TextOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
		OnHandler: func(payload OnTextHandlerPayload) {
			handlerCalled = true
			assert.Equal(t, "new_value", payload.Value) // Проверяем, что передано правильное значение
		},
	}

	vc := NewVirtualTextControl(opt)

	err := client.Publish(wb.PublishPayload{
		Value: "new_value",
		QOS:   0,
		Topic: vc.GetInfo().CommandTopic,
	})
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Проверяем, что обработчик был вызван
	assert.True(t, handlerCalled)
}

func TestVirtualTextControlAddWatcher(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := "default_value"

	vc := NewVirtualTextControl(TextOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	})

	var watcherCalled bool
	// Добавляем watcher для контроля изменений
	vc.AddWatcher(func(payload control.WatcherPayload) {
		watcherCalled = true
		assert.Equal(t, "new_value", payload.NewValue)     // Проверяем, что новое значение корректное
		assert.Equal(t, "default_value", payload.OldValue) // Проверяем, что старое значение корректное
	})

	// Устанавливаем новое значение, что должно вызвать срабатывание watcher
	vc.SetValue("new_value")

	time.Sleep(50 * time.Millisecond)

	// Проверяем, что watcher был вызван
	assert.True(t, watcherCalled)
}

func TestVirtualTextControlMetaType(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := "default_value"

	opt := TextOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualTextControl(opt)

	// Проверяем, что Meta.Type корректно установлен в "text"
	assert.Equal(t, "text", vc.control.meta.Type)
}
