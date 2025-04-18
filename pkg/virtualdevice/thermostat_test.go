package virtualdevice

import (
	"context"
	"testing"
	"time"

	"github.com/ValentinAlekhin/wb-go/internal/dbmock"
	"github.com/ValentinAlekhin/wb-go/internal/mqttmock"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewThermostat_InvalidConfig проверяет корректную обработку некорректных конфигураций
func TestNewThermostat_InvalidConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		ctx    context.Context
		config ThermostatConfig
		errMsg string
	}{
		{
			name:   "Nil Client",
			ctx:    context.Background(),
			config: ThermostatConfig{DB: dbmock.NewDBMock(), Device: "TestThermostat"},
			errMsg: "client is nil",
		},
		{
			name:   "Nil Queries",
			ctx:    context.Background(),
			config: ThermostatConfig{Client: mqttmock.NewMockClient(), Device: "TestThermostat"},
			errMsg: "db is nil",
		},
		{
			name:   "Empty Device",
			ctx:    context.Background(),
			config: ThermostatConfig{DB: dbmock.NewDBMock(), Client: mqttmock.NewMockClient()},
			errMsg: "device is empty",
		},
		{
			name:   "Negative Hysteresis",
			ctx:    context.Background(),
			config: ThermostatConfig{DB: dbmock.NewDBMock(), Client: mqttmock.NewMockClient(), Device: "TestThermostat", Hysteresis: -0.5},
			errMsg: "hysteresis must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thermostat, err := NewThermostat(tt.ctx, tt.config)
			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.errMsg)
				assert.Nil(t, thermostat)
			}
		})
	}
}

// TestNewThermostat_Initialization проверяет корректную инициализацию термостата
func TestNewThermostat_Initialization(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	config := ThermostatConfig{
		DB:                database,
		Client:            client,
		Device:            "TestThermostat",
		TargetTemperature: 25,
		Hysteresis:        0.5,
	}

	thermostat, err := NewThermostat(ctx, config)
	require.NoError(t, err)
	require.NotNil(t, thermostat)

	assert.Equal(t, "TestThermostat", thermostat.meta.Name)
	assert.Equal(t, 0.5, thermostat.hysteresis)
	assert.Equal(t, DefaultUpdateInterval, thermostat.updateInterval)
	assert.Equal(t, 25, thermostat.Controls.TargetTemperature.GetValue())
	assert.True(t, thermostat.Controls.Enabled.GetValue())
	assert.False(t, thermostat.Controls.Relay.GetValue())
}

// TestThermostat_UpdateRelay проверяет корректную работу реле в зависимости от температуры
func TestThermostat_UpdateRelay(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	config := ThermostatConfig{
		DB:                database,
		Client:            client,
		Device:            "TestThermostat",
		TargetTemperature: 25,
		Hysteresis:        0.5,
	}

	thermostat, err := NewThermostat(ctx, config)
	require.NoError(t, err)

	tests := []struct {
		name          string
		enabled       bool
		currentTemp   float64
		expectedRelay bool
	}{
		{
			name:          "Thermostat disabled",
			enabled:       false,
			currentTemp:   24.0,
			expectedRelay: false,
		},
		{
			name:          "Temperature below lower threshold",
			enabled:       true,
			currentTemp:   24.0, // 25 - 0.5 = 24.5
			expectedRelay: true,
		},
		{
			name:          "Temperature above upper threshold",
			enabled:       true,
			currentTemp:   26.0, // 25 + 0.5 = 25.5
			expectedRelay: false,
		},
		{
			name:          "Temperature between thresholds",
			enabled:       true,
			currentTemp:   25.0,
			expectedRelay: false, // Сохраняет предыдущее состояние
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thermostat.Controls.Enabled.SetValue(tt.enabled)
			thermostat.Controls.CurrentTemperature.SetValue(tt.currentTemp)
			thermostat.updateRelay()
			assert.Equal(t, tt.expectedRelay, thermostat.Controls.Relay.GetValue())
		})
	}
}

// TestThermostat_SetUpdateInterval проверяет корректную работу метода SetUpdateInterval
func TestThermostat_SetUpdateInterval(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	config := ThermostatConfig{
		DB:                database,
		Client:            client,
		Device:            "TestThermostat",
		TargetTemperature: 25,
		Hysteresis:        0.5,
	}

	thermostat, err := NewThermostat(ctx, config)
	require.NoError(t, err)

	// Проверяем начальное значение
	assert.Equal(t, DefaultUpdateInterval, thermostat.updateInterval)

	// Устанавливаем новый интервал
	newInterval := 2 * time.Second
	thermostat.SetUpdateInterval(newInterval)
	assert.Equal(t, newInterval, thermostat.updateInterval)

	// Проверяем, что отрицательный интервал игнорируется
	thermostat.SetUpdateInterval(-1 * time.Second)
	assert.Equal(t, newInterval, thermostat.updateInterval)

	// Проверяем, что нулевой интервал игнорируется
	thermostat.SetUpdateInterval(0)
	assert.Equal(t, newInterval, thermostat.updateInterval)
}

// TestThermostat_ContextCancellation проверяет корректную остановку тикера при отмене контекста
func TestThermostat_ContextCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	config := ThermostatConfig{
		DB:                database,
		Client:            client,
		Device:            "TestThermostat",
		TargetTemperature: 25,
		Hysteresis:        0.5,
	}

	thermostat, err := NewThermostat(ctx, config)
	require.NoError(t, err)

	thermostat.Controls.CurrentTemperature.SetValue(20)

	// Отменяем контекст
	cancel()

	// Даем время горутине завершиться
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что тикер больше не вызывает update
	thermostat.Controls.CurrentTemperature.SetValue(30.0)
	time.Sleep(2 * DefaultUpdateInterval) // Ждем достаточное время для срабатывания тикера
	assert.Equal(t, 20.0, thermostat.Controls.CurrentTemperature.GetValue(), "После отмены контекста значение не должно изменяться автоматически")
}

// TestThermostat_MetaPublishing проверяет корректную публикацию метаданных в MQTT
func TestThermostat_MetaPublishing(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	config := ThermostatConfig{
		DB:                database,
		Client:            client,
		Device:            "TestThermostat",
		TargetTemperature: 25,
		Hysteresis:        0.5,
	}

	thermostat, err := NewThermostat(ctx, config)
	require.NoError(t, err)

	messageChan := make(chan string, 1)
	err = client.Subscribe(thermostat.metaTopic, func(client mqtt.Client, msg mqtt.Message) {
		messageChan <- string(msg.Payload())
	})
	require.NoError(t, err)

	select {
	case msg := <-messageChan:
		expectedMeta := `{"name":"TestThermostat","driver":"wb-go"}`
		assert.JSONEq(t, expectedMeta, msg)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Не дождались сообщения с метаданными в MQTT-топике")
	}
}

// TestThermostat_Stop проверяет корректную остановку термостата
func TestThermostat_Stop(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := mqttmock.NewMockClient()
	database := dbmock.NewDBMock()

	config := ThermostatConfig{
		DB:                database,
		Client:            client,
		Device:            "TestThermostat",
		TargetTemperature: 25,
		Hysteresis:        0.5,
	}

	thermostat, err := NewThermostat(ctx, config)
	require.NoError(t, err)

	// Запоминаем состояние
	thermostat.Controls.CurrentTemperature.SetValue(20.0)
	thermostat.Controls.Enabled.SetValue(true)
	thermostat.Controls.Relay.SetValue(true)
	relayState := thermostat.Controls.Relay.GetValue()

	// Останавливаем термостат
	thermostat.Stop()

	// Изменяем температуру до значения, которое должно переключить реле
	thermostat.Controls.CurrentTemperature.SetValue(30.0)

	// Вызываем update напрямую (это не должно происходить автоматически после Stop)
	thermostat.update()

	// Проверяем, что реле изменило состояние при вызове update напрямую
	assert.NotEqual(t, relayState, thermostat.Controls.Relay.GetValue(),
		"Состояние реле должно измениться при ручном вызове update")
}
