package gen

import (
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
)

const mqttTopic = "/devices/+/controls/+/meta"
const deviceTemplateFile = "templates/device.txt"

var controlValueTypeMap = map[string][]string{
	"value": {
		conventions.CONV_TYPE_TEMPERATURE,
		conventions.CONV_TYPE_REL_HUMIDITY,
		conventions.CONV_TYPE_ATMOSPHERIC_PRESSURE,
		conventions.CONV_TYPE_RAINFALL,
		conventions.CONV_TYPE_WIND_SPEED,
		conventions.CONV_TYPE_POWER,
		conventions.CONV_TYPE_POWER_CONSUMPTION,
		conventions.CONV_TYPE_VOLTAGE,
		conventions.CONV_TYPE_WATER_FLOW,
		conventions.CONV_TYPE_WATER_CONSUMPTION,
		conventions.CONV_TYPE_RESISTANCE,
		conventions.CONV_TYPE_CONCENTRATION,
		conventions.CONV_TYPE_PRESSURE,
		conventions.CONV_TYPE_ILLUMINANCE,
		conventions.CONV_TYPE_SOUND_LEVEL,
		conventions.CONV_TYPE_HEAT_POWER,
		conventions.CONV_TYPE_HEAT_ENERGY,
		conventions.CONV_TYPE_CURRENT,
	},
}
