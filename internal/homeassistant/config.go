package homeassistant

type Device struct {
	Identifiers  string `json:"identifiers"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Name         string `json:"name"`
}

type Availability struct {
	Topic               string `json:"topic"`
	ValueTemplate       string `json:"value_template"`
	PayloadNotAvailable bool   `json:"payload_not_available"`
	PayloadAvailable    bool   `json:"payload_available"`
}

type Config struct {
	Device                   Device         `json:"device"`
	AvailabilityMode         string         `json:"availability_mode"`
	EnabledByDefault         bool           `json:"enabled_by_default"`
	Availability             []Availability `json:"availability"`
	BrightnessStateTopic     string         `json:"brightness_state_topic"`
	BrightnessCommandTopic   string         `json:"brightness_command_topic"`
	BrightnessScale          int            `json:"brightness_scale"`
	MaxMireds                int            `json:"max_mireds"`
	MinMireds                int            `json:"min_mireds"`
	ColorTempStateTopic      string         `json:"color_temp_state_topic"`
	ColorTempCommandTopic    string         `json:"color_temp_command_topic"`
	ColorTempValueTemplate   string         `json:"color_temp_value_template"`
	ColorTempCommandTemplate string         `json:"color_temp_command_template"`
	Name                     string         `json:"name"`
	UniqueID                 string         `json:"unique_id"`
	ObjectID                 string         `json:"object_id"`
	StateTopic               string         `json:"state_topic"`
	PayloadOn                int            `json:"payload_on"`
	PayloadOff               int            `json:"payload_off"`
	CommandTopic             string         `json:"command_topic"`
}
