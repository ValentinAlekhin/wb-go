package virualcontrol

import (
	"github.com/ValentinAlekhin/wb-go/internal/dbmock"
	"github.com/ValentinAlekhin/wb-go/internal/mqttmock"
	"github.com/ValentinAlekhin/wb-go/internal/testutils"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/pkg/timeonly"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestVirtualTimeControlGetValue(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := timeonly.NewTime(14, 30, 0)

	opt := TimeOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualTimeControl(opt)

	// Verify that the default value is returned correctly
	assert.Equal(t, defaultValue.String(), vc.GetValue().String())
}

func TestVirtualTimeControlSetValue(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := timeonly.NewTime(8, 15, 0)

	opt := TimeOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualTimeControl(opt)

	// Set a new value
	newValue := timeonly.NewTime(10, 30, 45)
	vc.SetValue(newValue)

	// Verify that the value was set correctly
	assert.Equal(t, newValue.String(), vc.GetValue().String())
}

func TestVirtualTimeControlOnHandler(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := timeonly.NewTime(6, 45, 0)

	handlerCalled := false

	opt := TimeOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
		OnHandler: func(payload OnTimeHandlerPayload) {
			handlerCalled = true
			assert.Equal(t, "10:30:45", payload.Value.String()) // Verify the correct value is passed
		},
	}

	vc := NewVirtualTimeControl(opt)

	err := client.Publish(wb.PublishPayload{
		Value: "10:30:45",
		QOS:   0,
		Topic: vc.GetInfo().CommandTopic,
	})
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Verify that the handler was called
	assert.True(t, handlerCalled)
}

func TestVirtualTimeControlAddWatcher(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := timeonly.NewTime(7, 0, 0)

	vc := NewVirtualTimeControl(TimeOptions{
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
	// Add a watcher to monitor changes
	vc.AddWatcher(func(payload TimeControlWatcherPayload) {
		watcherCalled = true
		assert.Equal(t, "10:30:45", payload.NewValue.String())            // Verify the new value is correct
		assert.Equal(t, defaultValue.String(), payload.OldValue.String()) // Verify the old value is correct
	})

	// Set a new value, which should trigger the watcher
	newValue := timeonly.NewTime(10, 30, 45)
	vc.SetValue(newValue)

	time.Sleep(50 * time.Millisecond)

	// Verify that the watcher was called
	assert.True(t, watcherCalled)
}

func TestVirtualTimeControlMetaType(t *testing.T) {
	t.Parallel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := timeonly.NewTime(12, 0, 0)

	opt := TimeOptions{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualTimeControl(opt)

	// Verify that Meta.Type is correctly set to "text"
	assert.Equal(t, "text", vc.control.meta.Type)
}
