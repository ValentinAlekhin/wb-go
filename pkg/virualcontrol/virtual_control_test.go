package virualcontrol

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ValentinAlekhin/wb-go/internal/db"
	"github.com/ValentinAlekhin/wb-go/internal/dbmock"
	"github.com/ValentinAlekhin/wb-go/internal/mqttmock"
	"github.com/ValentinAlekhin/wb-go/internal/testutils"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const device = "test-device" // Константа для устройства

func TestVirtualControlInitialization(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta: control.Meta{
				Type:  "value",
				Units: "°C",
			},
		},
		DefaultValue: "0",
	}

	vc := NewVirtualControl(ctx, opt)

	assert.Equal(t, controlName, vc.GetInfo().Name)
	assert.Equal(t, "0", vc.GetValue())
}

func TestVirtualControlSetValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: "0",
	}

	vc := NewVirtualControl(ctx, opt)

	vc.SetValue("25")
	assert.Equal(t, "25", vc.GetValue())

	var model db.ControlModel
	err := database.First(&model, "topic = ?", vc.GetInfo().ValueTopic).Error
	require.NoError(t, err)
	assert.Equal(t, "25", model.Value)
}

func TestVirtualControlWatchers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: "0",
	}

	vc := NewVirtualControl(ctx, opt)

	var payloads []control.WatcherPayloadString
	vc.AddWatcher(func(payload control.WatcherPayloadString) {
		payloads = append(payloads, payload)
	})

	vc.SetValue("42")
	time.Sleep(100 * time.Millisecond) // Подождем, пока обработчик выполнится

	require.Len(t, payloads, 1)
	assert.Equal(t, "0", payloads[0].OldValue)
	assert.Equal(t, "42", payloads[0].NewValue)
	assert.Equal(t, vc.GetInfo().ValueTopic, payloads[0].Topic)
}

func TestVirtualControlMQTTIntegration(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: "0",
	}

	vc := NewVirtualControl(ctx, opt)

	// Подписываемся на MQTT-топик и проверяем сообщения
	messageChan := make(chan string, 2)
	err := client.Subscribe(vc.GetInfo().ValueTopic, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println(msg.Topic(), string(msg.Payload()))
		messageChan <- string(msg.Payload())
	})
	require.NoError(t, err)

	<-messageChan

	// Устанавливаем новое значение, которое должно быть отправлено в MQTT
	vc.SetValue("50")

	// Проверяем, что сообщение с новым значением пришло в канал
	select {
	case msg := <-messageChan:
		assert.Equal(t, "50", msg)
	case <-time.After(1 * time.Second):
		t.Fatal("Не дождались сообщения в MQTT-топике")
	}
}

func TestVirtualControlDefaultValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: "42",
	}

	vc := NewVirtualControl(ctx, opt)

	assert.Equal(t, "42", vc.GetValue())
}

func TestVirtualControlMetaData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)

	meta := control.Meta{
		Type:      "value",
		Units:     "°C",
		Max:       100,
		Min:       0,
		Precision: 0.1,
		Order:     1,
		ReadOnly:  false,
		Title:     control.MultilingualText{"en": "Temperature", "ru": "Температура"},
	}

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   meta,
		},
		DefaultValue: "0",
	}

	vc := NewVirtualControl(ctx, opt)

	assert.Equal(t, meta.Type, vc.GetInfo().Meta.Type)
	assert.Equal(t, meta.Units, vc.GetInfo().Meta.Units)
	assert.Equal(t, meta.Max, vc.GetInfo().Meta.Max)
	assert.Equal(t, meta.Min, vc.GetInfo().Meta.Min)
	assert.Equal(t, meta.Precision, vc.GetInfo().Meta.Precision)
	assert.Equal(t, meta.Order, vc.GetInfo().Meta.Order)
	assert.Equal(t, meta.ReadOnly, vc.GetInfo().Meta.ReadOnly)
	assert.Equal(t, meta.Title["en"], vc.GetInfo().Meta.Title["en"])
	assert.Equal(t, meta.Title["ru"], vc.GetInfo().Meta.Title["ru"])
}

func TestVirtualControlOnHandler(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)

	var handlerCalled bool
	var lastSetValue string

	onHandler := func(payload OnHandlerPayload[string]) {
		handlerCalled = true
		lastSetValue = payload.Value
	}

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: "0",
		OnHandler:    onHandler,
	}

	vc := NewVirtualControl(ctx, opt)

	assert.False(t, handlerCalled)

	err := client.Publish(wb.PublishPayload{
		Value: "99",
		QOS:   0,
		Topic: vc.GetInfo().CommandTopic,
	})
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	assert.True(t, handlerCalled)
	assert.Equal(t, "99", lastSetValue)
}

func TestVirtualControlDefaultValueInTopic(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := "42"

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualControl(ctx, opt)

	messageChan := make(chan string, 1)
	err := client.Subscribe(vc.GetInfo().ValueTopic, func(client mqtt.Client, msg mqtt.Message) {
		messageChan <- string(msg.Payload())
	})
	require.NoError(t, err)

	select {
	case msg := <-messageChan:
		assert.Equal(t, defaultValue, msg)
	case <-time.After(1 * time.Second):
		t.Fatal("Не дождались сообщения с дефолтным значением в MQTT-топике")
	}
}

func TestVirtualControlMetaInTopic(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)

	meta := control.Meta{
		Type:      "value",
		Units:     "°C",
		Max:       100,
		Min:       11,
		Precision: 0.1,
		Order:     1,
		ReadOnly:  false,
		Title:     control.MultilingualText{"en": "Temperature", "ru": "Температура"},
	}

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   meta,
		},
		DefaultValue: "0",
	}

	vc := NewVirtualControl(ctx, opt)

	metaTopic := fmt.Sprintf("%s/meta", vc.GetInfo().ValueTopic)
	messageChan := make(chan string, 1)
	err := client.Subscribe(metaTopic, func(client mqtt.Client, msg mqtt.Message) {
		messageChan <- string(msg.Payload())
	})
	require.NoError(t, err)

	select {
	case msg := <-messageChan:
		expectedMeta := `{"type":"value","units":"°C","max":100,"min":11,"precision":0.1,"order":1,"readonly":false,"title":{"en":"Temperature","ru":"Температура"}}`
		assert.JSONEq(t, expectedMeta, msg)
	case <-time.After(1 * time.Second):
		t.Fatal("Не дождались сообщения с метаданными в MQTT-топике")
	}
}

func TestVirtualControlNoDuplicatePushesWithMqtt(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)
	defaultValue := "42"
	topic := fmt.Sprintf(conventions.CONV_CONTROL_VALUE_FMT, device, controlName)

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualControl(ctx, opt)

	messageChan := make(chan string, 10)

	err := client.Subscribe(topic, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println(string(msg.Payload()))
		messageChan <- string(msg.Payload())
	})
	require.NoError(t, err)

	vc.SetValue(defaultValue)
	vc.SetValue(defaultValue)
	vc.SetValue(defaultValue)
	vc.SetValue("50")

	select {
	case msg := <-messageChan:
		fmt.Println("read: ", msg)
		assert.Equal(t, defaultValue, msg)
	case <-time.After(1 * time.Second):
		t.Fatal("Не дождались сообщения с дефолтным значением")
	}

	select {
	case msg := <-messageChan:
		assert.Equal(t, "50", msg)
	case <-time.After(1 * time.Second):
		t.Fatal("Не дождались сообщения с новым значением")
	}

	select {
	case <-messageChan:
		t.Fatal("Не должно быть лишних сообщений")
	default:
	}
}

func TestVirtualControlContextCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	controlName := testutils.RandString(10)

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     database,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: "0",
	}

	vc := NewVirtualControl(ctx, opt)

	// Проверяем, что контрол работает до отмены контекста
	vc.SetValue("25")
	assert.Equal(t, "25", vc.GetValue())

	// Отменяем контекст
	cancel()
	time.Sleep(100 * time.Millisecond) // Даем время на завершение

	// Проверяем, что после отмены контекста SetValue не работает
	assert.NotPanics(t, func() {
		vc.SetValue("50")
	})

	assert.Equal(t, "25", vc.GetValue()) // Значение не должно измениться

	assert.NotPanics(t, func() {
		vc.AddWatcher(func(payload control.WatcherPayloadString) {
		})
	})
}
