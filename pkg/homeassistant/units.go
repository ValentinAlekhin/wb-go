package homeassistant

import (
	"errors"
)

// ConvertWBUnitToHA converts a unit from Wiren Board format to Home Assistant format
func ConvertWBUnitToHA(wbUnit string) (string, error) {
	if haUnit, exists := UnitConversionMap[wbUnit]; exists {
		return haUnit, nil
	}
	return "", errors.New("unknown unit: " + wbUnit)
}

// ConvertMetaTypeToUnit converts a meta/type to the corresponding unit
func ConvertMetaTypeToUnit(metaType string) (string, error) {
	if unit, exists := TypeToUnitMap[metaType]; exists {
		return unit, nil
	}
	return "", errors.New("unknown meta/type: " + metaType)
}
