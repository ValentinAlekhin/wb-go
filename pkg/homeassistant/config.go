package homeassistant

import (
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
)

func GetConfig(config MqttDiscoveryConfig) MqttDiscoveryConfig {
	cfg := &config
	applyDefault(cfg)
	applyAvailability(cfg)

	return *cfg
}

func applyDefault(config *MqttDiscoveryConfig) {
	config.EnabledByDefault = true

	if config.Device.Manufacturer == "" {
		config.Device.Manufacturer = "Wiren Board"
	}

	if config.PayloadOn == "" {
		config.PayloadOn = conventions.CONV_META_BOOL_TRUE
	}

	if config.PayloadOff == "" {
		config.PayloadOff = conventions.CONV_META_BOOL_FALSE
	}

	if config.AvailabilityMode == "" {
		config.AvailabilityMode = "latest"
	}

	if config.BrightnessScale == 0 {
		config.BrightnessScale = 100
	}
}

func applyAvailability(config *MqttDiscoveryConfig) {
	if config.Availability != nil {
		return
	}

	list := []MqttDiscoveryAvailability{
		{
			Topic:               config.StateTopic,
			ValueTemplate:       "{{ False if value == '' else True }}",
			PayloadNotAvailable: false,
			PayloadAvailable:    true,
		},
		{
			Topic:               config.StateTopic + "/meta",
			ValueTemplate:       "{{ False if value == '' else True }}",
			PayloadNotAvailable: false,
			PayloadAvailable:    true,
		},
		{
			Topic:               config.StateTopic + "/meta/error",
			ValueTemplate:       "{{ True if value == '' else False }}",
			PayloadNotAvailable: false,
			PayloadAvailable:    true,
		},
	}

	config.Availability = list
}
