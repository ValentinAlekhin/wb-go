package conventions

const (
	//
	// Format strings for fmt.Printf to form topic names
	// Use '+' as values to form subscription topics

	// Control value topic format string
	// Parameters:
	// 1. Device name
	// 2. Control name
	CONV_CONTROL_VALUE_FMT = "/devices/%s/controls/%s"

	// Control 'on' value topic format string
	// Used to pass new value to device from external software
	// Parameters:
	// 1. Device name
	// 2. Control name
	CONV_CONTROL_ON_VALUE_FMT = CONV_CONTROL_VALUE_FMT + "/on"

	// Device meta info topic format string
	// Parameters:
	// 1. Device name
	// 2. Meta subtopic name
	CONV_DEVICE_META_FMT = "/devices/%s/meta/%s"

	// Device meta info v2 topic format string
	// Parameters:
	// 1. Device name
	CONV_DEVICE_META_V2_FMT = "/devices/%s/meta"

	// Device driver topic format string
	// Parameters:
	// 1. Device name
	CONV_DEVICE_META_DRIVER_FMT = "/devices/%s/meta/driver"

	// Device control meta info topic format string
	// Parameters:
	// 1. Device name
	// 2. Control name
	// 3. Meta subtopic name
	CONV_CONTROL_META_FMT = "/devices/%s/controls/%s/meta/%s"

	// Device control meta info v2 topic format string
	// Parameters:
	// 1. Device name
	// 2. Control name
	CONV_CONTROL_META_V2_FMT = "/devices/%s/controls/%s/meta"

	// Device control all meta info topic format string
	// Parameters:
	// 1. Device name
	// 2. Control name
	CONV_CONTROL_ALL_META_FMT = "/devices/%s/controls/%s/meta/+"

	//
	// Meta information subtopics

	CONV_META_SUBTOPIC_DRIVER        = "driver"      // for /devices/+/meta/driver
	CONV_META_SUBTOPIC_TITLE         = "name"        // for /devices/+/meta/name ('name' is legacy)
	CONV_META_SUBTOPIC_TITLE_V2      = "title"       // for /devices/+/meta/title
	CONV_META_SUBTOPIC_ERROR         = "error"       // for /devices/+/controls/+/meta/error and /devices/+/meta/error
	CONV_META_SUBTOPIC_ORDER         = "order"       // for /devices/+/controls/+/meta/order
	CONV_META_SUBTOPIC_TYPE          = "type"        // for /devices/+/controls/+/meta/type
	CONV_META_SUBTOPIC_UNITS         = "units"       // for /devices/+/controls/+/meta/units
	CONV_META_SUBTOPIC_MAX           = "max"         // for /devices/+/controls/+/meta/max
	CONV_META_SUBTOPIC_MIN           = "min"         // for /devices/+/controls/+/meta/min
	CONV_META_SUBTOPIC_PRECISION     = "precision"   // for /devices/+/controls/+/meta/precision
	CONV_META_SUBTOPIC_DESCRIPTION   = "description" // for /devices/+/controls/+/meta/description
	CONV_META_SUBTOPIC_CONTROL_TITLE = "title"       // for /devices/+/controls/+/meta/title
	CONV_META_SUBTOPIC_READONLY      = "readonly"    // for /devices/+/controls/+/meta/readonly
	CONV_META_SUBTOPIC_CONTROL_ENUM  = "enum"        // for /devices/+/controls/+/meta/enum

	// Type names
	CONV_TYPE_SWITCH     = "switch"
	CONV_TYPE_ALARM      = "alarm"
	CONV_TYPE_PUSHBUTTON = "pushbutton"
	CONV_TYPE_RANGE      = "range"
	CONV_TYPE_RGB        = "rgb"
	CONV_TYPE_TEXT       = "text"
	CONV_TYPE_VALUE      = "value"

	// Meta types (types derived from 'value')
	CONV_TYPE_TEMPERATURE          = "temperature"
	CONV_TYPE_REL_HUMIDITY         = "rel_humidity"
	CONV_TYPE_ATMOSPHERIC_PRESSURE = "atmospheric_pressure"
	CONV_TYPE_RAINFALL             = "rainfall"
	CONV_TYPE_WIND_SPEED           = "wind_speed"
	CONV_TYPE_POWER                = "power"
	CONV_TYPE_POWER_CONSUMPTION    = "power_consumption"
	CONV_TYPE_VOLTAGE              = "voltage"
	CONV_TYPE_WATER_FLOW           = "water_flow"
	CONV_TYPE_WATER_CONSUMPTION    = "water_consumption"
	CONV_TYPE_RESISTANCE           = "resistance"
	CONV_TYPE_CONCENTRATION        = "concentration"
	CONV_TYPE_PRESSURE             = "pressure"
	CONV_TYPE_ILLUMINANCE          = "lux"
	CONV_TYPE_SOUND_LEVEL          = "sound_level"
	CONV_TYPE_HEAT_POWER           = "heat_power"
	CONV_TYPE_HEAT_ENERGY          = "heat_energy"
	CONV_TYPE_CURRENT              = "current"

	// Default data type for unknown meta type
	CONV_DEFAULT_DATATYPE = CONV_DATATYPE_STRING

	CONV_META_BOOL_TRUE  = "1"
	CONV_META_BOOL_FALSE = "0"

	// Special values for types
	CONV_SWITCH_VALUE_TRUE  = CONV_META_BOOL_TRUE
	CONV_SWITCH_VALUE_FALSE = CONV_META_BOOL_FALSE
	CONV_ALARM_VALUE_TRUE   = CONV_META_BOOL_TRUE
	CONV_ALARM_VALUE_FALSE  = CONV_META_BOOL_FALSE

	CONV_SWITCH_DEFAULT_VALUE = CONV_SWITCH_VALUE_FALSE
	CONV_ALARM_DEFAULT_VALUE  = CONV_ALARM_VALUE_FALSE

	CONV_SUBTOPIC_ALL = "+"
)

// ControlDataType is a real data types used in representations
type ControlDataType int

const (
	CONV_DATATYPE_STRING ControlDataType = iota
	CONV_DATATYPE_BOOLEAN
	CONV_DATATYPE_FLOAT
	CONV_DATATYPE_BUTTON
)
