package gen

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/testutils"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGenerateService_CollectDevicesData(t *testing.T) {
	t.Parallel()

	testClient, _, destroy := testutils.GetClientWithBroker()
	defer destroy()
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
}

//func TestGenerateService_CollectControlsData(t *testing.T) {
//	t.Parallel()
//
//	testClient, _, destroy := testutils.GetClientWithBroker()
//	defer destroy()
//
//	testControlMeta := `{"type": "switch", "readOnly": false}`
//	devName := testutils.RandString(10)
//
//	controlNames := make([]string, 0)
//	for i := range 10 {
//		name := fmt.Sprintf("control_%d", i)
//		controlNames = append(controlNames, name)
//
//		err := testClient.Publish(wb.PublishPayload{
//			Topic:    fmt.Sprintf(conventions.CONV_CONTROL_META_V2_FMT, devName, name),
//			Value:    testControlMeta,
//			QOS:      1,
//			Retained: true,
//		})
//		require.NoError(t, err)
//	}
//
//	// Инициализация сервиса
//	service := NewGenerateService(testClient, "./output", "testpkg")
//	controls := service.collectControlsData()
//
//	// Проверка результатов
//
//	require.Len(t, controls, len(controlNames))
//
//	for _, name := range controlNames {
//		var c watchControlResultItem
//		for _, control := range controls {
//			if name == control.Control {
//				c = control
//			}
//		}
//
//		require.NotNil(t, c)
//		require.Equal(t, devName, c.DeviceName)
//		require.Equal(t, name, c.Control)
//		require.Equal(t, "switch", c.Meta.Type)
//		require.False(t, c.Meta.ReadOnly)
//	}
//}

func TestGenerateService_CollectControlsData(t *testing.T) {
	t.Parallel()

	testClient, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

	testControlMeta := `{"type": "switch", "readOnly": false}`
	devName := testutils.RandString(10)

	controlNames := make([]string, 0)
	for i := range 10 {
		name := fmt.Sprintf("control_%d", i)
		controlNames = append(controlNames, name)

		err := testClient.Publish(wb.PublishPayload{
			Topic:    fmt.Sprintf(conventions.CONV_CONTROL_META_V2_FMT, devName, name),
			Value:    testControlMeta,
			QOS:      1,
			Retained: true,
		})
		require.NoError(t, err)
	}

	// Инициализация сервиса
	service := NewGenerateService(testClient, "./output", "testpkg")
	controls := service.collectControlsData()

	// Проверка результатов

	require.Len(t, controls, len(controlNames))

	for _, name := range controlNames {
		var c watchControlResultItem
		for _, control := range controls {
			if name == control.Control {
				c = control
			}
		}

		require.NotNil(t, c)
		require.Equal(t, devName, c.DeviceName)
		require.Equal(t, name, c.Control)
		require.Equal(t, "switch", c.Meta.Type)
		require.False(t, c.Meta.ReadOnly)
	}
}

func TestGenerateService_GenerateFiles(t *testing.T) {
	t.Parallel()

	testClient, _, destroy := testutils.GetClientWithBroker()
	defer destroy()

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
}
