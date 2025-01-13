package homeassistant

const DefaultDiscoveryName = "wb-go"
const DiscoveryMetaTopic = "wb-go-meta"

var TypeToUnitMap = map[string]string{
	"temperature":          "deg C",  // Температура
	"rel_humidity":         "%, RH",  // Относительная влажность
	"atmospheric_pressure": "mbar",   // Атмосферное давление
	"rainfall":             "mm/h",   // Интенсивность осадков
	"wind_speed":           "m/s",    // Скорость ветра
	"power":                "W",      // Мощность
	"power_consumption":    "kWh",    // Потребление энергии
	"voltage":              "V",      // Напряжение
	"water_flow":           "m^3/h",  // Расход воды
	"water_consumption":    "m^3",    // Объем воды
	"resistance":           "Ohm",    // Сопротивление
	"concentration":        "ppm",    // Концентрация газа
	"heat_power":           "Gcal/h", // Тепловая мощность
	"heat_energy":          "Gcal",   // Тепловая энергия
	"current":              "A",      // Сила тока
	"pressure":             "bar",    // Давление
	"lux":                  "lx",     // Освещенность
	"sound_level":          "dB",     // Уровень звука
}

// UnitConversionMap Map for unit conversion from Wiren Board to Home Assistant
var UnitConversionMap = map[string]string{
	"mm/h":   "mm/h",   // Осадки
	"m/s":    "m/s",    // Скорость
	"W":      "W",      // Мощность
	"kWh":    "kWh",    // Потребление энергии
	"V":      "V",      // Напряжение
	"mV":     "mV",     // Милливольты
	"m^3/h":  "m³/h",   // Расход воды
	"m^3":    "m³",     // Объем воды
	"Gcal/h": "Gcal/h", // Тепловая мощность
	"cal":    "cal",    // Калории
	"Gcal":   "Gcal",   // Гига калории
	"Ohm":    "Ω",      // Сопротивление
	"mOhm":   "mΩ",     // Миллиомы
	"bar":    "bar",    // Давление
	"mbar":   "mbar",   // Миллибары
	"s":      "s",      // Секунды
	"min":    "min",    // Минуты
	"h":      "h",      // Часы
	"m":      "m",      // Метры
	"g":      "g",      // Граммы
	"kg":     "kg",     // Килограммы
	"mol":    "mol",    // Моль
	"cd":     "cd",     // Канделы
	"%, RH":  "%",      // Влажность
	"deg C":  "°C",     // Градусы Цельсия
	"%":      "%",      // Проценты
	"ppm":    "ppm",    // Части на миллион
	"ppb":    "ppb",    // Части на миллиард
	"A":      "A",      // Амперы
	"mA":     "mA",     // Миллиамперы
	"deg":    "°",      // Градусы
	"rad":    "rad",    // Радианы
	"lx":     "lx",     // Люксы
	"dB":     "dB",     // Децибелы
	"Hz":     "Hz",     // Герцы
	"rpm":    "rpm",    // Обороты в минуту
	"Pa":     "Pa",     // Паскали
}
