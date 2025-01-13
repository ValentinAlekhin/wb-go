package homeassistant

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"regexp"
	"time"
)

type Discovery struct {
	client wb.ClientInterface
	prefix string
	name   string
}

type DiscoveryOptions struct {
	Client wb.ClientInterface
	Prefix string
	Name   string
}

type DiscoveryMeta struct {
	ClientName string    `json:"client_name"`
	CreatedAt  time.Time `json:"created_at"`
}

type DeviceConfig struct {
	getters         []*ConfigGetter
	ignoreRegexpStr []string
	ignoreRegexp    []*regexp.Regexp
}

type ConfigGetter struct {
	regexpStr string
	regexp    *regexp.Regexp
	getter    ConfigGetterFn
	domain    string
}

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

type ConfigGetterFn func(deviceInfo basedevice.Info, controlInfo control.Info) MqttDiscoveryConfig
type ConfigMiddleware func(config *MqttDiscoveryConfig, device basedevice.Info, control control.Info)
