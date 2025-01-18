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

func TestVirtualValueControlGetValue(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := 42.42

	opt := ValueOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualValueControl(opt)

	// Проверяем, что значение по умолчанию корректно возвращается
	assert.Equal(t, defaultValue, vc.GetValue())
}

func TestVirtualValueControlSetValue(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := 0.0

	opt := ValueOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualValueControl(opt)

	// Устанавливаем новое значение
	newValue := 25.75
	vc.SetValue(newValue)

	// Проверяем, что значение установилось
	assert.Equal(t, newValue, vc.GetValue())
}

func TestVirtualValueControlOnHandler(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := 0.0

	handlerCalled := false

	opt := ValueOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
		OnHandler: func(payload OnValueHandlerPayload) {
			handlerCalled = true
			assert.Equal(t, 25.75, payload.Value) // Проверяем, что передано правильное значение
		},
	}

	vc := NewVirtualValueControl(opt)

	err := client.Publish(wb.PublishPayload{
		Value: "25.75",
		QOS:   0,
		Topic: vc.GetInfo().CommandTopic,
	})
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Проверяем, что обработчик был вызван
	assert.True(t, handlerCalled)
}

func TestVirtualValueControlAddWatcher(t *testing.T) {
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := 0.0

	vc := NewVirtualValueControl(ValueOptions{
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
	vc.AddWatcher(func(payload control.ValueControlWatcherPayload) {
		watcherCalled = true
		assert.Equal(t, 25.75, payload.NewValue)        // Проверяем, что новое значение корректное
		assert.Equal(t, defaultValue, payload.OldValue) // Проверяем, что старое значение корректное
	})

	// Устанавливаем новое значение, что должно вызвать срабатывание watcher
	newValue := 25.75
	vc.SetValue(newValue)

	time.Sleep(100 * time.Millisecond)

	// Проверяем, что watcher был вызван
	assert.True(t, watcherCalled)
}

func TestVirtualValueControlMetaType(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := 0.0

	opt := ValueOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualValueControl(opt)

	// Проверяем, что Meta.Type корректно установлен в "value"
	assert.Equal(t, "value", vc.control.meta.Type)
}
