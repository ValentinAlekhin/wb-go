package virtualdevice

import (
	"fmt"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/pkg/virualcontrol"
	"github.com/ValentinAlekhin/wb-go/testutils"
	"github.com/ValentinAlekhin/wb-go/testutils/test_mqtt_server"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/glebarez/sqlite"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ValentinAlekhin/wb-go/pkg/timeonly"
	"gorm.io/gorm"
)

var testClient wb.ClientInterface
var testDB *gorm.DB

func TestMain(m *testing.M) {
	test_mqtt_server.StartMQTTBroker(false)
	testClient = testutils.GetMqttClient()

	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = database.AutoMigrate(&virualcontrol.ControlModel{})
	if err != nil {
		panic(err)
	}

	testDB = database

	fmt.Println("DONE")

	m.Run()

	database.Where("1 = 1").Delete(&virualcontrol.ControlModel{})
}

func TestNewAdaptiveLight_NilDB(t *testing.T) {
	config := AdaptiveLightConfig{
		DB:     nil,
		Client: testClient,
		Device: "TestDevice",
	}

	al, err := NewAdaptiveLight(config)
	assert.Nil(t, al)
	assert.EqualError(t, err, "db is nil")
}

func TestNewAdaptiveLight_EmptyDevice(t *testing.T) {
	config := AdaptiveLightConfig{
		DB:     testDB,
		Client: testClient,
		Device: "",
	}

	al, err := NewAdaptiveLight(config)
	assert.Nil(t, al)
	assert.EqualError(t, err, "device is empty")
}

func TestNewAdaptiveLight_Initialization(t *testing.T) {
	config := AdaptiveLightConfig{
		DB:     testDB,
		Client: testClient,
		Device: "TestDevice",
	}

	al, err := NewAdaptiveLight(config)
	require.NoError(t, err)
	require.NotNil(t, al)

	assert.Equal(t, fmt.Sprintf("%s_TestDevice", DevicePrefix), al.GetInfo().Name)
	assert.NotNil(t, al.Controls.Enabled)
	assert.NotNil(t, al.Controls.MinTemp)
	assert.NotNil(t, al.Controls.MaxTemp)
	assert.NotNil(t, al.Controls.CurrentTemp)
}

func TestNewAdaptiveLight_InvalidConfig(t *testing.T) {
	_, err := NewAdaptiveLight(AdaptiveLightConfig{
		DB:     nil,
		Client: nil,
		Device: "",
	})
	assert.Error(t, err)
}

func TestAdaptiveLight_GetInfo(t *testing.T) {
	config := AdaptiveLightConfig{
		DB:     testDB,
		Client: testClient,
		Device: "TestDevice",
	}

	al, err := NewAdaptiveLight(config)
	require.NoError(t, err)

	info := al.GetInfo()
	assert.Equal(t, fmt.Sprintf("%s_TestDevice", DevicePrefix), info.Name)
}

func TestAdaptiveLight_Update_Disabled(t *testing.T) {
	config := AdaptiveLightConfig{
		DB:     testDB,
		Client: testClient,
		Device: "TestDevice",
	}

	al, err := NewAdaptiveLight(config)
	require.NoError(t, err)

	// Disable AdaptiveLight
	al.Controls.Enabled.SetValue(false)

	// Ensure update does not proceed
	al.update()
	assert.False(t, al.Controls.SleepMode.GetValue())
}

func TestAdaptiveLight_SleepMode(t *testing.T) {
	config := AdaptiveLightConfig{
		DB:     testDB,
		Client: testClient,
		Device: "TestDevice",
	}

	al, err := NewAdaptiveLight(config)
	require.NoError(t, err)

	tests := []struct {
		name     string
		start    timeonly.Time
		end      timeonly.Time
		now      timeonly.Time
		expected bool
	}{
		{
			"SleepMode: now is within active time range",
			timeonly.NewTime(23, 0, 0),
			timeonly.NewTime(6, 0, 0),
			timeonly.NewTime(8, 0, 0),
			false,
		},
		{
			"SleepMode: now equals start time",
			timeonly.NewTime(23, 0, 0),
			timeonly.NewTime(6, 0, 0),
			timeonly.NewTime(23, 0, 0),
			false,
		},
		{
			"SleepMode: now is after start time but before end time",
			timeonly.NewTime(23, 0, 0),
			timeonly.NewTime(6, 0, 0),
			timeonly.NewTime(23, 0, 1),
			true,
		},
		{
			"SleepMode: now is just after midnight and falls into sleep time",
			timeonly.NewTime(23, 0, 0),
			timeonly.NewTime(6, 0, 0),
			timeonly.NewTime(0, 0, 1),
			true,
		},
		{
			"SleepMode: now equals end time",
			timeonly.NewTime(23, 0, 0),
			timeonly.NewTime(6, 0, 0),
			timeonly.NewTime(6, 0, 0),
			false,
		},
		{
			"SleepMode: now is outside the sleep time range, later in the day",
			timeonly.NewTime(0, 30, 0),
			timeonly.NewTime(6, 0, 0),
			timeonly.NewTime(23, 0, 0),
			false,
		},
		{
			"SleepMode: now is before sleep time range starts",
			timeonly.NewTime(0, 30, 0),
			timeonly.NewTime(6, 0, 0),
			timeonly.NewTime(0, 25, 0),
			false,
		},
		{
			"SleepMode: now is just after the start of sleep time range",
			timeonly.NewTime(0, 30, 0),
			timeonly.NewTime(6, 0, 0),
			timeonly.NewTime(1, 0, 0),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := al.getSleepMode(tt.start, tt.end, tt.now)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAdaptiveLight_GetBrightness(t *testing.T) {
	config := AdaptiveLightConfig{
		DB:     testDB,
		Client: testClient,
		Device: "TestDevice",
	}

	al, err := NewAdaptiveLight(config)
	require.NoError(t, err)

	maxBrightness := 100
	minBrightness := 20

	assert.Equal(t, minBrightness, al.getBrightness(maxBrightness, minBrightness, true))
	assert.Equal(t, maxBrightness, al.getBrightness(maxBrightness, minBrightness, false))
}

func TestAdaptiveLight_GetColorTemp(t *testing.T) {
	config := AdaptiveLightConfig{
		DB:     testDB,
		Client: testClient,
		Device: "TestDevice",
	}

	al, err := NewAdaptiveLight(config)
	require.NoError(t, err)

	maxTemp := 6500
	minTemp := 2700
	sunrise := timeonly.NewTime(6, 0, 0)
	sunset := timeonly.NewTime(18, 0, 0)

	al.now = timeonly.NewTime(12, 0, 0) // Noon
	assert.Greater(t, al.getColorTemp(false, maxTemp, minTemp, sunrise, sunset, al.now), minTemp)

	al.now = timeonly.NewTime(23, 0, 0) // Night
	assert.Equal(t, minTemp, al.getColorTemp(true, maxTemp, minTemp, sunrise, sunset, al.now))
}

func TestAdaptiveLight_MetaPublishing(t *testing.T) {
	config := AdaptiveLightConfig{
		DB:     testDB,
		Client: testClient,
		Device: "TestDevice",
	}

	al, err := NewAdaptiveLight(config)
	require.NoError(t, err)

	messageChan := make(chan string, 1)
	err = testClient.Subscribe(al.GetInfo().MetaTopic, func(client mqtt.Client, msg mqtt.Message) {
		messageChan <- string(msg.Payload())
	})
	require.NoError(t, err)

	select {
	case msg := <-messageChan:
		// Проверяем, что метаданные пришли в правильном формате
		expectedMeta := `{"name":"TestDevice","driver":"wb-go"}`
		assert.JSONEq(t, expectedMeta, msg)
	case <-time.After(1 * time.Second):
		t.Fatal("Не дождались сообщения с метаданными в MQTT-топике")
	}

}
