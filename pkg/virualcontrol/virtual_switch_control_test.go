package virualcontrol

import (
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestVirtualSwitchControlGetValue(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)

	opt := SwitchOptions{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: true,
	}

	vc := NewVirtualSwitchControl(opt)

	// Проверяем, что значение по умолчанию корректно возвращается
	assert.True(t, vc.GetValue())
}

func TestVirtualSwitchControlSetValue(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)

	opt := SwitchOptions{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
	}

	vc := NewVirtualSwitchControl(opt)

	// Устанавливаем новое значение
	vc.SetValue(true)

	// Проверяем, что значение установилось
	assert.Equal(t, true, vc.GetValue())
}

func TestVirtualSwitchControlOnHandler(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)
	handlerCalled := false

	opt := SwitchOptions{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		OnHandler: func(payload OnSwitchHandlerPayload) {
			handlerCalled = true
			assert.Equal(t, true, payload.Value) // Проверяем, что передано правильное значение
		},
	}

	vc := NewVirtualSwitchControl(opt)

	err := client.Publish(wb.PublishPayload{
		Value: conventions.CONV_META_BOOL_TRUE,
		QOS:   0,
		Topic: vc.GetInfo().CommandTopic,
	})
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Проверяем, что обработчик был вызван
	assert.True(t, handlerCalled)
}

func TestVirtualSwitchControlAddWatcher(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)
	vc := NewVirtualSwitchControl(SwitchOptions{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
	})

	var watcherCalled bool
	// Добавляем watcher для контроля изменений
	vc.AddWatcher(func(payload control.SwitchControlWatcherPayload) {
		watcherCalled = true
		assert.Equal(t, true, payload.NewValue)  // Проверяем, что новое значение корректное
		assert.Equal(t, false, payload.OldValue) // Проверяем, что старое значение корректное
	})

	// Устанавливаем новое значение, что должно вызвать срабатывание watcher
	vc.SetValue(true)

	time.Sleep(100 * time.Millisecond)

	// Проверяем, что watcher был вызван
	assert.True(t, watcherCalled)
}

func TestVirtualSwitchControlMetaType(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)
	opt := SwitchOptions{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
	}

	vc := NewVirtualSwitchControl(opt)

	// Проверяем, что Meta.Type корректно установлен в "switch"
	assert.Equal(t, "switch", vc.control.meta.Type)
}

func TestVirtualSwitchControlToggle(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)

	opt := SwitchOptions{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
	}

	vc := NewVirtualSwitchControl(opt)

	// Проверяем начальное состояние
	assert.False(t, vc.GetValue())

	// Тогглим значение
	vc.Toggle()
	assert.True(t, vc.GetValue())

	// Тогглим снова
	vc.Toggle()
	assert.False(t, vc.GetValue())
}

func TestVirtualSwitchControlTurnOff(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)

	opt := SwitchOptions{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: true,
	}

	vc := NewVirtualSwitchControl(opt)

	// Проверяем начальное состояние
	assert.True(t, vc.GetValue())

	// Выключаем
	vc.TurnOff()
	assert.False(t, vc.GetValue())
}

func TestVirtualSwitchControlTurnOn(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)

	opt := SwitchOptions{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
	}

	vc := NewVirtualSwitchControl(opt)

	// Проверяем начальное состояние
	assert.False(t, vc.GetValue())

	// Включаем
	vc.TurnOn()
	assert.True(t, vc.GetValue())
}
