package gen

const mqttTopic = "/devices/+/controls/+/meta"
const deviceTemplateFile = "templates/device.txt"

var controlValueTypeMap = map[string][]string{
	"value": {
		"value",
		"voltage",
		"lux",
		"rel_humidity",
		"temperature",
		"concentration",
		"sound_level",
	},
}

var controlFileNames = []string{
	"control",
	"pushbutton_control",
	"range_control",
	"switch_control",
	"text_control",
	"value_control",
}
