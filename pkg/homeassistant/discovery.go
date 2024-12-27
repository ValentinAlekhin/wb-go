package homeassistant

import (
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
)

type MqttDiscoveryDevice struct {
	Identifiers  string `json:"identifiers,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Model        string `json:"model,omitempty"`
	Name         string `json:"name,omitempty"`
}

type MqttDiscoveryAvailability struct {
	Topic               string `json:"topic"`
	ValueTemplate       string `json:"value_template,omitempty"`
	PayloadNotAvailable bool   `json:"payload_not_available"`
	PayloadAvailable    bool   `json:"payload_available"`
}

type MqttDiscoveryConfig struct {
	Device                   MqttDiscoveryDevice         `json:"device,omitempty"`
	AvailabilityMode         string                      `json:"availability_mode,omitempty"`
	EnabledByDefault         bool                        `json:"enabled_by_default,omitempty"`
	Availability             []MqttDiscoveryAvailability `json:"availability,omitempty"`
	RGBStateTopic            string                      `json:"rgb_state_topic,omitempty"`
	RGBCommandTopic          string                      `json:"rgb_command_topic,omitempty"`
	RGBValueTemplate         string                      `json:"rgb_value_template,omitempty"`
	RGBCommandTemplate       string                      `json:"rgb_command_template,omitempty"`
	Name                     string                      `json:"name,omitempty"`
	UniqueID                 string                      `json:"unique_id,omitempty"`
	ObjectID                 string                      `json:"object_id,omitempty"`
	StateTopic               string                      `json:"state_topic,omitempty"`
	PayloadOn                string                      `json:"payload_on,omitempty"`
	PayloadOff               string                      `json:"payload_off,omitempty"`
	CommandTopic             string                      `json:"command_topic,omitempty"`
	BrightnessCommandTopic   string                      `json:"brightness_command_topic,omitempty"`
	BrightnessStateTopic     string                      `json:"brightness_state_topic,omitempty"`
	BrightnessScale          int                         `json:"brightness_scale,omitempty"`
	MaxMireds                int                         `json:"max_mireds,omitempty"`
	MinMireds                int                         `json:"min_mireds,omitempty"`
	ColorTempStateTopic      string                      `json:"color_temp_state_topic,omitempty"`
	ColorTempCommandTopic    string                      `json:"color_temp_command_topic,omitempty"`
	ColorTempValueTemplate   string                      `json:"color_temp_value_template,omitempty"`
	ColorTempCommandTemplate string                      `json:"color_temp_command_template,omitempty"`
	UnitOfMeasurement        string                      `json:"unit_of_measurement,omitempty"`
	DeviceClass              string                      `json:"device_class,omitempty"`
	Min                      float64                     `json:"min,omitempty"`
	Max                      float64                     `json:"max,omitempty"`
}

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
