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

func TestVirtualRangeControlGetValue(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := 42

	opt := RangeOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualRangeControl(opt)

	// Проверяем, что значение по умолчанию корректно возвращается
	assert.Equal(t, defaultValue, vc.GetValue())
}

func TestVirtualRangeControlSetValue(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := 0

	opt := RangeOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualRangeControl(opt)

	// Устанавливаем новое значение
	vc.SetValue(25)

	// Проверяем, что значение установилось
	assert.Equal(t, 25, vc.GetValue())
}

func TestVirtualRangeControlOnHandler(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := 0

	handlerCalled := false

	opt := RangeOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
		OnHandler: func(payload OnRangeHandlerPayload) {
			handlerCalled = true
			assert.Equal(t, 25, payload.Value) // Проверяем, что передано правильное значение
		},
	}

	vc := NewVirtualRangeControl(opt)

	err := client.Publish(wb.PublishPayload{
		Value: "25",
		QOS:   0,
		Topic: vc.GetInfo().CommandTopic,
	})
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Проверяем, что обработчик был вызван
	assert.True(t, handlerCalled)
}

func TestVirtualRangeControlAddWatcher(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := 0

	vc := NewVirtualRangeControl(RangeOptions{
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
	vc.AddWatcher(func(payload control.RangeControlWatcherPayload) {
		watcherCalled = true
		assert.Equal(t, 25, payload.NewValue) // Проверяем, что новое значение корректное
		assert.Equal(t, 0, payload.OldValue)  // Проверяем, что старое значение корректное
	})

	// Устанавливаем новое значение, что должно вызвать срабатывание watcher
	vc.SetValue(25)

	time.Sleep(100 * time.Millisecond)

	// Проверяем, что watcher был вызван
	assert.True(t, watcherCalled)
}

func TestVirtualRangeControlMetaType(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := 0

	opt := RangeOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualRangeControl(opt)

	// Проверяем, что Meta.Type корректно установлен в "range"
	assert.Equal(t, "range", vc.control.meta.Type)
}
