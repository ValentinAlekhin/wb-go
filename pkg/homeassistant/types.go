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
	Identifiers      string `json:"identifiers,omitempty"`
	Manufacturer     string `json:"manufacturer,omitempty"`
	Model            string `json:"model,omitempty"`
	Name             string `json:"name,omitempty"`
	ConfigurationURL string `json:"configuration_url,omitempty"`
	Connections      string `json:"connections,omitempty"`
	ModelID          string `json:"model_id,omitempty"`
	HardwareVersion  string `json:"hw_version,omitempty"`
	SoftwareVersion  string `json:"sw_version,omitempty"`
	SuggestedArea    string `json:"suggested_area,omitempty"`
	SerialNumber     string `json:"serial_number,omitempty"`
}

type MqttDiscoveryAvailability struct {
	Topic               string `json:"topic"`
	ValueTemplate       string `json:"value_template,omitempty"`
	PayloadNotAvailable bool   `json:"payload_not_available"`
	PayloadAvailable    bool   `json:"payload_available"`
}

type MqttDiscoveryConfig struct {
	Device                         MqttDiscoveryDevice         `json:"device,omitempty"`
	AvailabilityMode               string                      `json:"availability_mode,omitempty"`
	EnabledByDefault               bool                        `json:"enabled_by_default,omitempty"`
	Availability                   []MqttDiscoveryAvailability `json:"availability,omitempty"`
	RGBStateTopic                  string                      `json:"rgb_state_topic,omitempty"`
	RGBCommandTopic                string                      `json:"rgb_command_topic,omitempty"`
	RGBValueTemplate               string                      `json:"rgb_value_template,omitempty"`
	RGBCommandTemplate             string                      `json:"rgb_command_template,omitempty"`
	Name                           string                      `json:"name,omitempty"`
	UniqueID                       string                      `json:"unique_id,omitempty"`
	ObjectID                       string                      `json:"object_id,omitempty"`
	StateTopic                     string                      `json:"state_topic,omitempty"`
	PayloadOn                      string                      `json:"payload_on,omitempty"`
	PayloadOff                     string                      `json:"payload_off,omitempty"`
	CommandTopic                   string                      `json:"command_topic,omitempty"`
	BrightnessCommandTopic         string                      `json:"brightness_command_topic,omitempty"`
	BrightnessStateTopic           string                      `json:"brightness_state_topic,omitempty"`
	BrightnessScale                int                         `json:"brightness_scale,omitempty"`
	MaxMireds                      int                         `json:"max_mireds,omitempty"`
	MinMireds                      int                         `json:"min_mireds,omitempty"`
	ColorTempStateTopic            string                      `json:"color_temp_state_topic,omitempty"`
	ColorTempCommandTopic          string                      `json:"color_temp_command_topic,omitempty"`
	ColorTempValueTemplate         string                      `json:"color_temp_value_template,omitempty"`
	ColorTempCommandTemplate       string                      `json:"color_temp_command_template,omitempty"`
	UnitOfMeasurement              string                      `json:"unit_of_measurement,omitempty"`
	DeviceClass                    string                      `json:"device_class,omitempty"`
	Min                            float64                     `json:"min,omitempty"`
	Max                            float64                     `json:"max,omitempty"`
	ActionTopic                    string                      `json:"action_topic,omitempty"`
	ActionTemplate                 string                      `json:"action_template,omitempty"`
	AutomationType                 string                      `json:"automation_type,omitempty"`
	AuxCommandTopic                string                      `json:"aux_command_topic,omitempty"`
	AuxStateTopic                  string                      `json:"aux_state_topic,omitempty"`
	AuxStateTemplate               string                      `json:"aux_state_template,omitempty"`
	AvailableTones                 string                      `json:"available_tones,omitempty"`
	AvailabilityTopic              string                      `json:"availability_topic,omitempty"`
	AvailabilityTemplate           string                      `json:"availability_template,omitempty"`
	AwayModeCommandTopic           string                      `json:"away_mode_command_topic,omitempty"`
	AwayModeStateTopic             string                      `json:"away_mode_state_topic,omitempty"`
	AwayModeStateTemplate          string                      `json:"away_mode_state_template,omitempty"`
	BlueTemplate                   string                      `json:"blue_template,omitempty"`
	BrightnessCommandTemplate      string                      `json:"brightness_command_template,omitempty"`
	BrightnessTemplate             string                      `json:"brightness_template,omitempty"`
	BrightnessValueTemplate        string                      `json:"brightness_value_template,omitempty"`
	ColorTempTemplate              string                      `json:"color_temp_template,omitempty"`
	ColorMode                      string                      `json:"color_mode,omitempty"`
	ColorModeStateTopic            string                      `json:"color_mode_state_topic,omitempty"`
	ColorModeValueTemplate         string                      `json:"color_mode_value_template,omitempty"`
	CommandOffTemplate             string                      `json:"command_off_template,omitempty"`
	CommandOnTemplate              string                      `json:"command_on_template,omitempty"`
	CommandTemplate                string                      `json:"command_template,omitempty"`
	Components                     string                      `json:"components,omitempty"`
	CodeArmRequired                bool                        `json:"code_arm_required,omitempty"`
	CodeDisarmRequired             bool                        `json:"code_disarm_required,omitempty"`
	CodeTriggerRequired            bool                        `json:"code_trigger_required,omitempty"`
	ContentType                    string                      `json:"content_type,omitempty"`
	CurrentTemperatureTopic        string                      `json:"current_temperature_topic,omitempty"`
	CurrentTemperatureTemplate     string                      `json:"current_temperature_template,omitempty"`
	DirectionCommandTopic          string                      `json:"direction_command_topic,omitempty"`
	DirectionCommandTemplate       string                      `json:"direction_command_template,omitempty"`
	DirectionStateTopic            string                      `json:"direction_state_topic,omitempty"`
	DirectionValueTemplate         string                      `json:"direction_value_template,omitempty"`
	DisplayPrecision               string                      `json:"display_precision,omitempty"`
	Encoding                       string                      `json:"encoding,omitempty"`
	EntityCategory                 string                      `json:"entity_category,omitempty"`
	EntityPicture                  string                      `json:"entity_picture,omitempty"`
	EventTypes                     string                      `json:"event_types,omitempty"`
	ExpireAfter                    int                         `json:"expire_after,omitempty"`
	FanSpeedList                   string                      `json:"fan_speed_list,omitempty"`
	FlashTimeLong                  int                         `json:"flash_time_long,omitempty"`
	FlashTimeShort                 int                         `json:"flash_time_short,omitempty"`
	EffectCommandTopic             string                      `json:"effect_command_topic,omitempty"`
	EffectCommandTemplate          string                      `json:"effect_command_template,omitempty"`
	EffectList                     string                      `json:"effect_list,omitempty"`
	EffectStateTopic               string                      `json:"effect_state_topic,omitempty"`
	EffectTemplate                 string                      `json:"effect_template,omitempty"`
	EffectValueTemplate            string                      `json:"effect_value_template,omitempty"`
	FanModeCommandTopic            string                      `json:"fan_mode_command_topic,omitempty"`
	FanModeCommandTemplate         string                      `json:"fan_mode_command_template,omitempty"`
	FanModeStateTopic              string                      `json:"fan_mode_state_topic,omitempty"`
	FanModeStateTemplate           string                      `json:"fan_mode_state_template,omitempty"`
	ForceUpdate                    bool                        `json:"force_update,omitempty"`
	GreenTemplate                  string                      `json:"green_template,omitempty"`
	HSCommandTopic                 string                      `json:"hs_command_topic,omitempty"`
	HSCommandTemplate              string                      `json:"hs_command_template,omitempty"`
	HSStateTopic                   string                      `json:"hs_state_topic,omitempty"`
	HSValueTemplate                string                      `json:"hs_value_template,omitempty"`
	Icon                           string                      `json:"icon,omitempty"`
	ImageEncoding                  string                      `json:"image_encoding,omitempty"`
	ImageTopic                     string                      `json:"image_topic,omitempty"`
	Initial                        string                      `json:"initial,omitempty"`
	TargetHumidityCommandTopic     string                      `json:"target_humidity_command_topic,omitempty"`
	TargetHumidityCommandTemplate  string                      `json:"target_humidity_command_template,omitempty"`
	TargetHumidityStateTopic       string                      `json:"target_humidity_state_topic,omitempty"`
	TargetHumidityStateTemplate    string                      `json:"target_humidity_state_template,omitempty"`
	JSONAttributes                 string                      `json:"json_attributes,omitempty"`
	JSONAttributesTopic            string                      `json:"json_attributes_topic,omitempty"`
	JSONAttributesTemplate         string                      `json:"json_attributes_template,omitempty"`
	LatestVersionTopic             string                      `json:"latest_version_topic,omitempty"`
	LatestVersionTemplate          string                      `json:"latest_version_template,omitempty"`
	LastResetTopic                 string                      `json:"last_reset_topic,omitempty"`
	LastResetValueTemplate         string                      `json:"last_reset_value_template,omitempty"`
	MaxHumidity                    string                      `json:"max_humidity,omitempty"`
	MaxTemp                        string                      `json:"max_temp,omitempty"`
	MigrateDiscovery               bool                        `json:"migrate_discovery,omitempty"`
	MinHumidity                    string                      `json:"min_humidity,omitempty"`
	MinTemp                        string                      `json:"min_temp,omitempty"`
	Mode                           string                      `json:"mode,omitempty"`
	ModeCommandTopic               string                      `json:"mode_command_topic,omitempty"`
	ModeCommandTemplate            string                      `json:"mode_command_template,omitempty"`
	ModeStateTopic                 string                      `json:"mode_state_topic,omitempty"`
	ModeStateTemplate              string                      `json:"mode_state_template,omitempty"`
	Modes                          string                      `json:"modes,omitempty"`
	Origin                         string                      `json:"origin,omitempty"`
	OffDelay                       int                         `json:"off_delay,omitempty"`
	OnCommandType                  string                      `json:"on_command_type,omitempty"`
	Options                        string                      `json:"options,omitempty"`
	Optimistic                     bool                        `json:"optimistic,omitempty"`
	OscillationCommandTopic        string                      `json:"oscillation_command_topic,omitempty"`
	OscillationCommandTemplate     string                      `json:"oscillation_command_template,omitempty"`
	OscillationStateTopic          string                      `json:"oscillation_state_topic,omitempty"`
	OscillationValueTemplate       string                      `json:"oscillation_value_template,omitempty"`
	Platform                       string                      `json:"platform,omitempty"`
	PercentageCommandTopic         string                      `json:"percentage_command_topic,omitempty"`
	PercentageCommandTemplate      string                      `json:"percentage_command_template,omitempty"`
	PercentageStateTopic           string                      `json:"percentage_state_topic,omitempty"`
	PercentageValueTemplate        string                      `json:"percentage_value_template,omitempty"`
	Payload                        string                      `json:"payload,omitempty"`
	PayloadArmAway                 string                      `json:"payload_arm_away,omitempty"`
	PayloadArmCustomBypass         string                      `json:"payload_arm_custom_bypass,omitempty"`
	PayloadArmHome                 string                      `json:"payload_arm_home,omitempty"`
	PayloadArmNight                string                      `json:"payload_arm_night,omitempty"`
	PayloadArmVacation             string                      `json:"payload_arm_vacation,omitempty"`
	PayloadAvailable               string                      `json:"payload_available,omitempty"`
	PayloadCleanSpot               string                      `json:"payload_clean_spot,omitempty"`
	PayloadClose                   string                      `json:"payload_close,omitempty"`
	PayloadDirectionForward        string                      `json:"payload_direction_forward,omitempty"`
	PayloadDirectionReverse        string                      `json:"payload_direction_reverse,omitempty"`
	PayloadDisarm                  string                      `json:"payload_disarm,omitempty"`
	PayloadHome                    string                      `json:"payload_home,omitempty"`
	PayloadInstall                 string                      `json:"payload_install,omitempty"`
	PayloadLocate                  string                      `json:"payload_locate,omitempty"`
	PayloadLock                    string                      `json:"payload_lock,omitempty"`
	PayloadNotAvailable            string                      `json:"payload_not_available,omitempty"`
	PayloadNotHome                 string                      `json:"payload_not_home,omitempty"`
	PayloadOpen                    string                      `json:"payload_open,omitempty"`
	PayloadOscillationOff          string                      `json:"payload_oscillation_off,omitempty"`
	PayloadOscillationOn           string                      `json:"payload_oscillation_on,omitempty"`
	PayloadPause                   string                      `json:"payload_pause,omitempty"`
	PayloadStop                    string                      `json:"payload_stop,omitempty"`
	PayloadStart                   string                      `json:"payload_start,omitempty"`
	PayloadPress                   string                      `json:"payload_press,omitempty"`
	PayloadReturnToBase            string                      `json:"payload_return_to_base,omitempty"`
	PayloadReset                   string                      `json:"payload_reset,omitempty"`
	PayloadResetHumidity           string                      `json:"payload_reset_humidity,omitempty"`
	PayloadResetMode               string                      `json:"payload_reset_mode,omitempty"`
	PayloadResetPercentage         string                      `json:"payload_reset_percentage,omitempty"`
	PayloadResetPresetMode         string                      `json:"payload_reset_preset_mode,omitempty"`
	PayloadTurnOff                 string                      `json:"payload_turn_off,omitempty"`
	PayloadTurnOn                  string                      `json:"payload_turn_on,omitempty"`
	PayloadTrigger                 string                      `json:"payload_trigger,omitempty"`
	PayloadUnlock                  string                      `json:"payload_unlock,omitempty"`
	ReportsPosition                bool                        `json:"reports_position,omitempty"`
	PositionClosed                 string                      `json:"position_closed,omitempty"`
	PositionOpen                   string                      `json:"position_open,omitempty"`
	PresetModeCommandTopic         string                      `json:"preset_mode_command_topic,omitempty"`
	PresetModeCommandTemplate      string                      `json:"preset_mode_command_template,omitempty"`
	PresetModeStateTopic           string                      `json:"preset_mode_state_topic,omitempty"`
	PresetModeValueTemplate        string                      `json:"preset_mode_value_template,omitempty"`
	PresetModes                    string                      `json:"preset_modes,omitempty"`
	Pattern                        string                      `json:"pattern,omitempty"`
	RedTemplate                    string                      `json:"red_template,omitempty"`
	ReleaseSummary                 string                      `json:"release_summary,omitempty"`
	ReleaseURL                     string                      `json:"release_url,omitempty"`
	Retain                         bool                        `json:"retain,omitempty"`
	RGBWCommandTopic               string                      `json:"rgbw_command_topic,omitempty"`
	RGBWCommandTemplate            string                      `json:"rgbw_command_template,omitempty"`
	RGBWStateTopic                 string                      `json:"rgbw_state_topic,omitempty"`
	RGBWValueTemplate              string                      `json:"rgbw_value_template,omitempty"`
	RGBWWCommandTopic              string                      `json:"rgbww_command_topic,omitempty"`
	RGBWWCommandTemplate           string                      `json:"rgbww_command_template,omitempty"`
	RGBWWStateTopic                string                      `json:"rgbww_state_topic,omitempty"`
	RGBWWValueTemplate             string                      `json:"rgbww_value_template,omitempty"`
	SendCommandTopic               string                      `json:"send_command_topic,omitempty"`
	SendIfOff                      bool                        `json:"send_if_off,omitempty"`
	SetFanSpeedTopic               string                      `json:"set_fan_speed_topic,omitempty"`
	SetPositionTopic               string                      `json:"set_position_topic,omitempty"`
	SetPositionTemplate            string                      `json:"set_position_template,omitempty"`
	PositionTopic                  string                      `json:"position_topic,omitempty"`
	PositionTemplate               string                      `json:"position_template,omitempty"`
	SpeedRangeMin                  int                         `json:"speed_range_min,omitempty"`
	SpeedRangeMax                  int                         `json:"speed_range_max,omitempty"`
	SourceType                     string                      `json:"source_type,omitempty"`
	StateClass                     string                      `json:"state_class,omitempty"`
	StateClosing                   string                      `json:"state_closing,omitempty"`
	StateClosed                    string                      `json:"state_closed,omitempty"`
	StateJammed                    string                      `json:"state_jammed,omitempty"`
	StateLocked                    string                      `json:"state_locked,omitempty"`
	StateLocking                   string                      `json:"state_locking,omitempty"`
	StateOff                       string                      `json:"state_off,omitempty"`
	StateOn                        string                      `json:"state_on,omitempty"`
	StateOpen                      string                      `json:"state_open,omitempty"`
	StateOpening                   string                      `json:"state_opening,omitempty"`
	StateStopped                   string                      `json:"state_stopped,omitempty"`
	StateUnlocked                  string                      `json:"state_unlocked,omitempty"`
	StateUnlocking                 string                      `json:"state_unlocking,omitempty"`
	StateTemplate                  string                      `json:"state_template,omitempty"`
	StateValueTemplate             string                      `json:"state_value_template,omitempty"`
	Step                           string                      `json:"step,omitempty"`
	Subtype                        string                      `json:"subtype,omitempty"`
	SuggestedDisplayPrecision      string                      `json:"suggested_display_precision,omitempty"`
	SupportedColorModes            string                      `json:"supported_color_modes,omitempty"`
	SupportDuration                bool                        `json:"support_duration,omitempty"`
	SupportVolumeSet               bool                        `json:"support_volume_set,omitempty"`
	SupportedFeatures              string                      `json:"supported_features,omitempty"`
	SwingModeCommandTopic          string                      `json:"swing_mode_command_topic,omitempty"`
	SwingModeCommandTemplate       string                      `json:"swing_mode_command_template,omitempty"`
	SwingModeStateTopic            string                      `json:"swing_mode_state_topic,omitempty"`
	SwingModeStateTemplate         string                      `json:"swing_mode_state_template,omitempty"`
	Topic                          string                      `json:"topic,omitempty"`
	TemperatureCommandTopic        string                      `json:"temperature_command_topic,omitempty"`
	TemperatureCommandTemplate     string                      `json:"temperature_command_template,omitempty"`
	TemperatureHighCommandTopic    string                      `json:"temperature_high_command_topic,omitempty"`
	TemperatureHighCommandTemplate string                      `json:"temperature_high_command_template,omitempty"`
	TemperatureHighStateTopic      string                      `json:"temperature_high_state_topic,omitempty"`
	TemperatureHighStateTemplate   string                      `json:"temperature_high_state_template,omitempty"`
	TemperatureLowCommandTopic     string                      `json:"temperature_low_command_topic,omitempty"`
	TemperatureLowCommandTemplate  string                      `json:"temperature_low_command_template,omitempty"`
	TemperatureLowStateTopic       string                      `json:"temperature_low_state_topic,omitempty"`
	TemperatureLowStateTemplate    string                      `json:"temperature_low_state_template,omitempty"`
	TemperatureStateTopic          string                      `json:"temperature_state_topic,omitempty"`
	TemperatureStateTemplate       string                      `json:"temperature_state_template,omitempty"`
	TemperatureUnit                string                      `json:"temperature_unit,omitempty"`
	TiltClosedValue                string                      `json:"tilt_closed_value,omitempty"`
	TiltCommandTopic               string                      `json:"tilt_command_topic,omitempty"`
	TiltCommandTemplate            string                      `json:"tilt_command_template,omitempty"`
	TiltMax                        int                         `json:"tilt_max,omitempty"`
	TiltMin                        int                         `json:"tilt_min,omitempty"`
	TiltOpenedValue                string                      `json:"tilt_opened_value,omitempty"`
	TiltOptimistic                 bool                        `json:"tilt_optimistic,omitempty"`
	TiltStatusTopic                string                      `json:"tilt_status_topic,omitempty"`
	TiltStatusTemplate             string                      `json:"tilt_status_template,omitempty"`
	TopicTemplate                  string                      `json:"topic_template,omitempty"`
	UpdateInterval                 string                      `json:"update_interval,omitempty"`
	ValueTemplate                  string                      `json:"value_template,omitempty"`
	WhiteCommandTopic              string                      `json:"white_command_topic,omitempty"`
	WhiteCommandTemplate           string                      `json:"white_command_template,omitempty"`
	WhiteStateTopic                string                      `json:"white_state_topic,omitempty"`
	WhiteValueTemplate             string                      `json:"white_value_template,omitempty"`
	XyCommandTopic                 string                      `json:"xy_command_topic,omitempty"`
	XyCommandTemplate              string                      `json:"xy_command_template,omitempty"`
	XyStateTopic                   string                      `json:"xy_state_topic,omitempty"`
	XyValueTemplate                string                      `json:"xy_value_template,omitempty"`
}

type ConfigGetterFn func(deviceInfo basedevice.Info, controlInfo control.Info) MqttDiscoveryConfig
type ConfigMiddleware func(domain *string, config *MqttDiscoveryConfig, device basedevice.Info, control control.Info)
