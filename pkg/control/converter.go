package control

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/timeonly"
	"strconv"
	"strings"
)

// RGBConverter converts RgbValue to/from string.
type RGBConverter struct{}

// Encode converts RgbValue to a string in the format "R;G;B".
func (c *RGBConverter) Encode(value RgbValue) string {
	return fmt.Sprintf("%d;%d;%d", value.Red, value.Green, value.Blue)
}

// Decode converts a string in the format "R;G;B" to RgbValue.
func (c *RGBConverter) Decode(value string) (RgbValue, error) {
	parts := strings.Split(value, ";")
	if len(parts) != 3 {
		return RgbValue{}, fmt.Errorf("invalid RGB format: %s", value)
	}

	red, err := strconv.Atoi(parts[0])
	if err != nil {
		return RgbValue{}, fmt.Errorf("invalid red value: %s", parts[0])
	}

	green, err := strconv.Atoi(parts[1])
	if err != nil {
		return RgbValue{}, fmt.Errorf("invalid green value: %s", parts[1])
	}

	blue, err := strconv.Atoi(parts[2])
	if err != nil {
		return RgbValue{}, fmt.Errorf("invalid blue value: %s", parts[2])
	}

	return RgbValue{Red: red, Green: green, Blue: blue}, nil
}

// RangeConverter converts int to/from string.
type RangeConverter struct{}

// Encode converts int to string.
func (c *RangeConverter) Encode(value int) string {
	return strconv.Itoa(value)
}

// Decode converts string to int.
func (c *RangeConverter) Decode(value string) (int, error) {
	return strconv.Atoi(value)
}

// SwitchConverter converts bool to/from string.
type SwitchConverter struct{}

// Encode converts bool to string.
func (c *SwitchConverter) Encode(value bool) string {
	if value {
		return "1"
	} else {
		return "0"
	}
}

// Decode converts string to bool.
func (c *SwitchConverter) Decode(value string) (bool, error) {
	return strconv.ParseBool(value)
}

// ValueConverter converts float64 to/from string.
type ValueConverter struct{}

// Encode converts float64 to string.
func (c *ValueConverter) Encode(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

// Decode converts string to float64.
func (c *ValueConverter) Decode(value string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(value), 64)
}

// TimeOnlyConverter converts timeonly.Time to/from string.
type TimeOnlyConverter struct{}

// Encode converts timeonly.Time to a string.
func (c *TimeOnlyConverter) Encode(value timeonly.Time) string {
	return value.String()
}

// Decode converts a string to timeonly.Time.
func (c *TimeOnlyConverter) Decode(value string) (timeonly.Time, error) {
	t, err := timeonly.ParseString(value)
	if err != nil {
		return timeonly.Time{}, fmt.Errorf("invalid time format: %s", value)
	}
	return t, nil
}
