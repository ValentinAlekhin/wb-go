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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thermostat, err := NewThermostat(tt.ctx, tt.config)
			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.errMsg)
			}
			assert.NotNil(t, thermostat)
		})
	}
}

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
	assert.Equal(t, 25, thermostat.Controls.TargetTemperature.GetValue())
	assert.True(t, thermostat.Controls.Enabled.GetValue())
	assert.False(t, thermostat.Controls.Relay.GetValue())
}

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

	// Термостат выключен
	thermostat.Controls.Enabled.SetValue(false)
	thermostat.Controls.CurrentTemperature.SetValue(24.0)
	thermostat.updateRelay()
	assert.False(t, thermostat.Controls.Relay.GetValue())

	// Термостат включен, температура ниже нижнего порога
	thermostat.Controls.Enabled.SetValue(true)
	thermostat.Controls.CurrentTemperature.SetValue(24.0) // 25 - 0.5 = 24.5
	thermostat.updateRelay()
	assert.True(t, thermostat.Controls.Relay.GetValue())

	// Температура выше верхнего порога
	thermostat.Controls.CurrentTemperature.SetValue(26.0) // 25 + 0.5 = 25.5
	thermostat.updateRelay()
	assert.False(t, thermostat.Controls.Relay.GetValue())
}

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
	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, 20.0, thermostat.Controls.CurrentTemperature.GetValue())
}

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
