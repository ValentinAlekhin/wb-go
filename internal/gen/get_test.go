package gen

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/testutils"
	"github.com/ValentinAlekhin/wb-go/testutils/test_mqtt_server"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/require"
	"os"
	"sync"
	"testing"
	"time"
)

var testClient wb.ClientInterface

func TestMain(m *testing.M) {
	test_mqtt_server.StartMQTTBroker(false)
	testClient = testutils.GetMqttClient()
	m.Run()
}

func TestGenerateService_CollectDevicesData(t *testing.T) {
	// Публикация тестовых данных о девайсах
	testDeviceMeta := `{"driver": "test-driver"}`
	devName := testutils.RandString(10)
	err := testClient.Publish(wb.PublishPayload{
		Topic:    fmt.Sprintf(conventions.CONV_DEVICE_META_V2_FMT, devName),
		Value:    testDeviceMeta,
		QOS:      1,
		Retained: true,
	})
	require.NoError(t, err)

	// Инициализация сервиса
	service := NewGenerateService(testClient, "./output", "testpkg")
	devices := service.collectDevicesData()

	// Проверка результатов
	require.Len(t, devices, 1)
	require.Equal(t, devName, devices[0].Name)
	require.Equal(t, "test-driver", devices[0].Meta.Driver)

	err = ClearTopics(testClient)
	require.NoError(t, err)
}

func TestGenerateService_CollectControlsData(t *testing.T) {
	testControlMeta := `{"type": "switch", "readOnly": false}`
	devName := testutils.RandString(10)
	controlName := testutils.RandString(10)
	err := testClient.Publish(wb.PublishPayload{
		Topic:    fmt.Sprintf(conventions.CONV_CONTROL_META_V2_FMT, devName, controlName),
		Value:    testControlMeta,
		QOS:      1,
		Retained: true,
	})
	require.NoError(t, err)

	// Инициализация сервиса
	service := NewGenerateService(testClient, "./output", "testpkg")
	controls := service.collectControlsData()

	// Проверка результатов
	require.Len(t, controls, 1)
	require.Equal(t, devName, controls[0].DeviceName)
	require.Equal(t, controlName, controls[0].Control)
	require.Equal(t, "switch", controls[0].Meta.Type)
	require.False(t, controls[0].Meta.ReadOnly)

	err = ClearTopics(testClient)
	require.NoError(t, err)
}

func TestGenerateService_GenerateFiles(t *testing.T) {
	deviceMeta := `{"driver": "test-driver"}`
	controlMeta := `{"type": "switch", "readOnly": false}`
	devName := testutils.RandString(10)
	controlName := testutils.RandString(10)
	err := testClient.Publish(wb.PublishPayload{
		Topic:    fmt.Sprintf(conventions.CONV_DEVICE_META_V2_FMT, devName),
		Value:    deviceMeta,
		QOS:      1,
		Retained: true,
	})
	require.NoError(t, err)

	err = testClient.Publish(wb.PublishPayload{
		Topic:    fmt.Sprintf(conventions.CONV_CONTROL_META_V2_FMT, devName, controlName),
		Value:    controlMeta,
		QOS:      1,
		Retained: true,
	})
	require.NoError(t, err)

	// Создание директории для вывода
	outputDir := "./test_output"
	defer func(path string) {
		err := os.RemoveAll(path)
		require.NoError(t, err)
	}(outputDir)

	// Инициализация сервиса
	service := NewGenerateService(testClient, outputDir, "testpkg")
	service.Run()

	// Проверка, что файлы созданы
	files, err := os.ReadDir(outputDir)
	require.NoError(t, err)
	require.NotEmpty(t, files)

	err = ClearTopics(testClient)
	require.NoError(t, err)
}

func ClearTopics(client wb.ClientInterface) error {

	topics := make([]string, 0)

	tCh := make(chan string)

	watcher := func(client mqtt.Client, msg mqtt.Message) {
		tCh <- msg.Topic()
	}

	controlTopic := fmt.Sprintf(conventions.CONV_CONTROL_META_V2_FMT, "+", "+")
	devTopic := fmt.Sprintf(conventions.CONV_DEVICE_META_V2_FMT, "+")
	err := client.Subscribe(controlTopic, watcher)
	if err != nil {
		return err
	}
	err = client.Subscribe(devTopic, watcher)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		duration := 100 * time.Millisecond
		timer := time.NewTimer(duration)
		defer timer.Stop()

		for {
			select {
			case topic := <-tCh:
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(duration)
				topics = append(topics, topic)
			case <-timer.C:
				return
			}
		}
	}()

	wg.Wait()

	err = client.Unsubscribe(devTopic, controlTopic)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		fmt.Printf("Topic: %s\n", topic)

		err := client.Publish(wb.PublishPayload{
			Topic:    topic,
			Value:    "",
			QOS:      1,
			Retained: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
