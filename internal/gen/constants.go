package gen

var DevicesToGenerate = []string{"wb-led", "wb-msw-v4", "wb-ms", "wb-adc", "wb-mwac-v2", "wb-mr6cu", "wb-mdm3", "buzzer", "network", "hwmon", "power_status"}
var ControlValueTypeMap = map[string][]string{"value": {"value", "voltage", "lux", "rel_humidity", "temperature", "concentration", "sound_level"}}
