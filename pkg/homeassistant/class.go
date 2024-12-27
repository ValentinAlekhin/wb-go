package homeassistant

import (
	"errors"
)

// UnitToDeviceClassMap Map for converting units to device_class for Home Assistant
var UnitToDeviceClassMap = map[string]string{
	"deg C":  "temperature",                      // Температура
	"%, RH":  "humidity",                         // Относительная влажность
	"mbar":   "pressure",                         // Давление
	"mm/h":   "precipitation_intensity",          // Интенсивность осадков
	"m/s":    "wind_speed",                       // Скорость ветра
	"W":      "power",                            // Мощность
	"kWh":    "energy",                           // Энергопотребление
	"V":      "voltage",                          // Напряжение
	"mV":     "voltage",                          // Напряжение (милливольты)
	"m^3/h":  "volume_flow_rate",                 // Расход воды
	"m^3":    "volume",                           // Объем воды
	"Ohm":    "none",                             // Сопротивление (нет точного класса)
	"ppm":    "volatile_organic_compounds_parts", // Концентрация газа
	"Gcal/h": "power",                            // Тепловая мощность
	"Gcal":   "energy",                           // Тепловая энергия
	"A":      "current",                          // Ток
	"mA":     "current",                          // Ток (миллиамперы)
	"bar":    "pressure",                         // Давление
	"lx":     "illuminance",                      // Освещённость
	"dB":     "sound_pressure",                   // Уровень звука
}

// ConvertUnitToDeviceClass converts a unit from Wiren Board to device_class for Home Assistant
func ConvertUnitToDeviceClass(unit string) (string, error) {
	if deviceClass, exists := UnitToDeviceClassMap[unit]; exists {
		return deviceClass, nil
	}
	return "", errors.New("unknown unit: " + unit)
}
