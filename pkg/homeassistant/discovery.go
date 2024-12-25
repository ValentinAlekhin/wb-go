package homeassistant

import (
	"encoding/json"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
)

type MqttDiscoveryDevice struct {
	Identifiers  string `json:"identifiers"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Name         string `json:"name"`
}

type MqttDiscoveryAvailability struct {
	Topic               string `json:"topic"`
	ValueTemplate       string `json:"value_template"`
	PayloadNotAvailable bool   `json:"payload_not_available"`
	PayloadAvailable    bool   `json:"payload_available"`
}

type MqttDiscoveryConfig struct {
	Device             MqttDiscoveryDevice         `json:"device"`
	AvailabilityMode   string                      `json:"availability_mode"`
	EnabledByDefault   bool                        `json:"enabled_by_default"`
	Availability       []MqttDiscoveryAvailability `json:"availability"`
	RGBStateTopic      string                      `json:"rgb_state_topic"`
	RGBCommandTopic    string                      `json:"rgb_command_topic"`
	RGBValueTemplate   string                      `json:"rgb_value_template"`
	RGBCommandTemplate string                      `json:"rgb_command_template"`
	Name               string                      `json:"name"`
	UniqueID           string                      `json:"unique_id"`
	ObjectID           string                      `json:"object_id"`
	StateTopic         string                      `json:"state_topic"`
	PayloadOn          string                      `json:"payload_on"`
	PayloadOff         string                      `json:"payload_off"`
	CommandTopic       string                      `json:"command_topic"`
}

func GetJsonConfig(config MqttDiscoveryConfig) (string, error) {
	cfg := &config
	applyDefault(cfg)

	val, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}

	return string(val), nil
}

func applyDefault(config *MqttDiscoveryConfig) {
	if config.PayloadOn == "" {
		config.PayloadOn = conventions.CONV_META_BOOL_TRUE
	}

	if config.PayloadOff == "" {
		config.PayloadOff = conventions.CONV_META_BOOL_FALSE
	}

	if config.AvailabilityMode == "" {
		config.AvailabilityMode = "latest"
	}
}
