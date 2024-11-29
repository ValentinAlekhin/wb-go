package homeassistant

import "fmt"

func GetConfigForCct(name string, modbusAddress string, control string) (Config, string) {
	identifiers := fmt.Sprintf("%s_%s", name, modbusAddress)
	controlTopic := fmt.Sprintf("/devices/%s/controls/%s", identifiers, control)

	device := newDevice(identifiers, name)
	availabilityList := newAvailabilityList(controlTopic)

	config := Config{
		Device:                   device,
		AvailabilityMode:         "latest",
		EnabledByDefault:         true,
		Availability:             availabilityList,
		BrightnessStateTopic:     controlTopic + " Brightness",
		BrightnessCommandTopic:   controlTopic + " Brightness/on",
		BrightnessScale:          100,
		MaxMireds:                454,
		MinMireds:                154,
		ColorTempStateTopic:      controlTopic + " Temperature",
		ColorTempCommandTopic:    controlTopic + " Temperature/on",
		ColorTempValueTemplate:   "{{ ((((100 - value | float) / 100) * (this.attributes.max_mireds - this.attributes.min_mireds)) + this.attributes.min_mireds) | round(0) }}",
		ColorTempCommandTemplate: "{{ (100 - (((value - this.attributes.min_mireds) / (this.attributes.max_mireds - this.attributes.min_mireds)) * 100)) | round(0) }}",
		Name:                     "Лента " + control,
		UniqueID:                 identifiers,
		ObjectID:                 identifiers,
		StateTopic:               controlTopic,
		CommandTopic:             controlTopic + "/on",
		PayloadOn:                1,
		PayloadOff:               0,
	}

	return config, fmt.Sprintf("homeassistant/light/%s/%s/config", identifiers, control)
}

func newDevice(identifiers string, name string) Device {
	return Device{
		Identifiers:  identifiers,
		Manufacturer: "Wiren Board",
		Model:        name,
		Name:         name,
	}

}

func newAvailabilityList(topic string) []Availability {
	suffixes := []string{"", "/meta", "/meta/error"}
	var availabilityList []Availability

	for _, suffix := range suffixes {
		availability := Availability{
			Topic:               topic + suffix,
			ValueTemplate:       "{{ False if value == '' else True }}",
			PayloadNotAvailable: false,
			PayloadAvailable:    true,
		}
		availabilityList = append(availabilityList, availability)
	}

	return availabilityList
}
