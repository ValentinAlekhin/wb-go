package virualcontrol

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"testing"
	"time"

	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/testutils"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const device = "test-device" // Константа для устройства

var testDB *gorm.DB

func TestMain(m *testing.M) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = database.AutoMigrate(&ControlModel{})
	fmt.Println("MOGRATE", err)
	if err != nil {
		panic(err)
	}

	testDB = database

	m.Run()

	database.Where("1 = 1").Delete(&ControlModel{})
}

func TestVirtualControlInitialization(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     testDB,
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

	vc := NewVirtualControl(opt)

	assert.Equal(t, controlName, vc.GetInfo().Name)
	assert.Equal(t, "0", vc.GetValue())
}

func TestVirtualControlSetValue(t *testing.T) {
	//t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10) // Генерация случайного имени для контрола

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,      // Используем константу для устройства
			Name:   controlName, // Используем сгенерированное имя
			Meta:   control.Meta{},
		},
		DefaultValue: "0",
	}

	vc := NewVirtualControl(opt)

	vc.SetValue("25")
	assert.Equal(t, "25", vc.GetValue())

	var model ControlModel
	err := testDB.First(&model, "topic = ?", vc.GetInfo().ValueTopic).Error
	require.NoError(t, err)
	assert.Equal(t, "25", model.Value)
}

func TestVirtualControlWatchers(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10) // Генерация случайного имени для контрола

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,      // Используем константу для устройства
			Name:   controlName, // Используем сгенерированное имя
			Meta:   control.Meta{},
		},
		DefaultValue: "0",
	}

	vc := NewVirtualControl(opt)

	var payloads []control.WatcherPayload
	vc.AddWatcher(func(payload control.WatcherPayload) {
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

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10) // Генерация случайного имени для контрола

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,      // Используем константу для устройства
			Name:   controlName, // Используем сгенерированное имя
			Meta:   control.Meta{},
		},
		DefaultValue: "0", // Устанавливаем дефолтное значение
	}

	vc := NewVirtualControl(opt)

	// Перед подпиской отправляем пустое сообщение с флагом Retained, чтобы очистить старое значение
	_ = client.Publish(wb.PublishPayload{
		Topic:    vc.GetInfo().ValueTopic,
		Value:    "", // Отправляем пустое значение
		QOS:      1,
		Retained: true, // Устанавливаем Retained, чтобы очистить топик
	})

	// Подписываемся на MQTT-топик и проверяем сообщения
	messageChan := make(chan string, 1)
	err := client.Subscribe(vc.GetInfo().ValueTopic, func(client mqtt.Client, msg mqtt.Message) {
		messageChan <- string(msg.Payload())
	})
	require.NoError(t, err)

	// Устанавливаем новое значение, которое должно быть отправлено в MQTT
	vc.SetValue("50")

	// Проверяем, что сообщение с новым значением пришло в канал
	select {
	case msg := <-messageChan:
		assert.Equal(t, "50", msg) // Ожидаем значение "50"
	case <-time.After(1 * time.Second):
		t.Fatal("Не дождались сообщения в MQTT-топике")
	}
}

func TestVirtualControlDefaultValue(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: "42", // Устанавливаем дефолтное значение
	}

	vc := NewVirtualControl(opt)

	// Проверяем, что дефолтное значение установлено правильно
	assert.Equal(t, "42", vc.GetValue())
}

func TestVirtualControlMetaData(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

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
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   meta, // Устанавливаем метаданные
		},
		DefaultValue: "0",
	}

	vc := NewVirtualControl(opt)

	// Проверяем, что метаданные установлены правильно
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

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)

	// Создаем флаг для отслеживания изменений
	var handlerCalled bool
	var lastSetValue string

	// Создаем кастомный OnHandler, который будет вызываться при изменении значения
	onHandler := func(payload OnHandlerPayload) {
		handlerCalled = true
		lastSetValue = payload.Value
	}

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: "0",
		OnHandler:    onHandler, // Устанавливаем кастомный OnHandler
	}

	vc := NewVirtualControl(opt)

	// Проверяем, что изначально OnHandler не был вызван
	assert.False(t, handlerCalled)

	// Устанавливаем новое значение через SetValue
	err := client.Publish(wb.PublishPayload{
		Value: "99",
		QOS:   0,
		Topic: vc.GetInfo().CommandTopic,
	})
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Проверяем, что OnHandler был вызван и значение установлено правильно
	assert.True(t, handlerCalled)
	assert.Equal(t, "99", lastSetValue)
}

func TestVirtualControlDefaultValueInTopic(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)

	// Устанавливаем дефолтное значение
	defaultValue := "42"

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualControl(opt)

	// Подписываемся на топик значения, чтобы проверить, что дефолтное значение приходит
	messageChan := make(chan string, 1)
	err := client.Subscribe(vc.GetInfo().ValueTopic, func(client mqtt.Client, msg mqtt.Message) {
		messageChan <- string(msg.Payload())
	})
	require.NoError(t, err)

	// Ждем сообщение с дефолтным значением
	select {
	case msg := <-messageChan:
		assert.Equal(t, defaultValue, msg)
	case <-time.After(1 * time.Second):
		t.Fatal("Не дождались сообщения с дефолтным значением в MQTT-топике")
	}
}

func TestVirtualControlMetaInTopic(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)

	// Устанавливаем метаданные
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
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   meta, // Устанавливаем метаданные
		},
		DefaultValue: "0",
	}

	vc := NewVirtualControl(opt)

	// Подписываемся на топик метаданных, чтобы проверить, что метаданные приходят
	metaTopic := fmt.Sprintf("%s/meta", vc.GetInfo().ValueTopic)
	messageChan := make(chan string, 1)
	err := client.Subscribe(metaTopic, func(client mqtt.Client, msg mqtt.Message) {
		messageChan <- string(msg.Payload())
	})
	require.NoError(t, err)

	// Ждем сообщение с метаданными
	select {
	case msg := <-messageChan:
		// Проверяем, что метаданные пришли в правильном формате
		expectedMeta := `{"type":"value","units":"°C","max":100,"min":11,"precision":0.1,"order":1,"readonly":false,"title":{"en":"Temperature","ru":"Температура"}}`
		assert.JSONEq(t, expectedMeta, msg)
	case <-time.After(1 * time.Second):
		t.Fatal("Не дождались сообщения с метаданными в MQTT-топике")
	}
}

func TestVirtualControlNoDuplicatePushesWithMqtt(t *testing.T) {
	t.Parallel()

	client, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	controlName := testutils.RandString(10)
	defaultValue := "42"
	topic := fmt.Sprintf(conventions.CONV_CONTROL_VALUE_FMT, device, controlName)

	opt := Options{
		BaseOptions: BaseOptions{
			DB:     testDB,
			Client: client,
			Device: device,
			Name:   controlName,
			Meta:   control.Meta{},
		},
		DefaultValue: defaultValue,
	}

	vc := NewVirtualControl(opt)

	// Канал для получения сообщений MQTT
	messageChan := make(chan string, 10)

	// Подписываемся на MQTT-топик, на который будет отправляться значение
	err := client.Subscribe(topic, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println(string(msg.Payload()))
		messageChan <- string(msg.Payload())
	})
	require.NoError(t, err)

	// Устанавливаем значение первый раз
	vc.SetValue(defaultValue)

	// Устанавливаем то же самое значение второй раз
	vc.SetValue(defaultValue)

	// Устанавливаем то же самое значение третий раз
	vc.SetValue(defaultValue)

	// Устанавливаем другое значение
	vc.SetValue("50")

	// Ждем, пока публикации обработаются
	select {
	case msg := <-messageChan:
		fmt.Println("read: ", msg)
		assert.Equal(t, defaultValue, msg) // Получаем сообщение с дефолтным значением
	case <-time.After(1 * time.Second):
		t.Fatal("Не дождались сообщения с дефолтным значением")
	}

	// Ждем, пока будет отправлено новое значение
	select {
	case msg := <-messageChan:
		assert.Equal(t, "50", msg) // Получаем сообщение с новым значением
	case <-time.After(1 * time.Second):
		t.Fatal("Не дождались сообщения с новым значением")
	}

	// Проверяем, что не было лишних сообщений
	select {
	case <-messageChan:
		t.Fatal("Не должно быть лишних сообщений")
	default:
		// Это нормально, если сообщений больше не пришло
	}
}
