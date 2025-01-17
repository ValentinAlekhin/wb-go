package virtuladevice

import (
	"fmt"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/ValentinAlekhin/wb-go/pkg/virualcontrol"
	"github.com/ValentinAlekhin/wb-go/testutils"
	"github.com/ValentinAlekhin/wb-go/testutils/test_mqtt_server"
	"github.com/glebarez/sqlite"
	"testing"

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
	assert.NotNil(t, al.Controls.Sunrise)
	assert.NotNil(t, al.Controls.Sunset)
	assert.NotNil(t, al.Controls.SleepStart)
	assert.NotNil(t, al.Controls.SleepEnd)
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

	sleepStart := timeonly.NewTime(23, 0, 0)
	sleepEnd := timeonly.NewTime(6, 0, 0)

	al.now = timeonly.NewTime(23, 30, 0)
	assert.True(t, al.getSleepMode(sleepStart, sleepEnd, al.now))

	al.now = timeonly.NewTime(12, 0, 0)
	assert.False(t, al.getSleepMode(sleepStart, sleepEnd, al.now))
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

//func TestAdaptiveLight_MetaPublishing(t *testing.T) {
//	config := AdaptiveLightConfig{
//		DB:     testDB,
//		Client: testClient,
//		Device: "TestDevice",
//	}
//
//	al, err := NewAdaptiveLight(config)
//	require.NoError(t, err)
//
//	// Check if meta is published correctly
//	meta := Meta{
//		Name:   "TestDevice",
//		Driver: "wb-go",
//	}
//	metaBytes, _ := json.Marshal(meta)
//
//	require.Len(t, client.publishedPayloads, 1)
//	assert.Equal(t, client.publishedPayloads[0].Value, string(metaBytes))
//	assert.Equal(t, client.publishedPayloads[0].Topic, al.metaTopic)
//}
