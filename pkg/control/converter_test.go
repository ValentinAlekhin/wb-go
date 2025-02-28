package control

import (
	"github.com/ValentinAlekhin/wb-go/pkg/timeonly"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRGBConverter_Encode(t *testing.T) {
	t.Parallel()

	converter := &RGBConverter{}

	// Тест кодирования
	value := RgbValue{Red: 255, Green: 0, Blue: 128}
	encoded := converter.Encode(value)
	assert.Equal(t, "255;0;128", encoded, "RGB encoding failed")
}

func TestRGBConverter_Decode(t *testing.T) {
	converter := &RGBConverter{}

	// Тест декодирования корректного значения
	encoded := "255;0;128"
	decoded, err := converter.Decode(encoded)
	assert.NoError(t, err, "RGB decoding failed")
	assert.Equal(t, RgbValue{Red: 255, Green: 0, Blue: 128}, decoded, "RGB decoding result is incorrect")

	// Тест с некорректным форматом (недостаточно частей)
	_, err = converter.Decode("255;0")
	assert.Error(t, err, "Expected error for invalid RGB format")

	// Тест с некорректным значением для красного компонента
	_, err = converter.Decode("abc;0;128")
	assert.Error(t, err, "Expected error for invalid red value")

	// Тест с некорректным значением для зеленого компонента
	_, err = converter.Decode("255;abc;128")
	assert.Error(t, err, "Expected error for invalid green value")
	assert.Contains(t, err.Error(), "invalid green value", "Error message should mention green value")

	// Тест с некорректным значением для синего компонента
	_, err = converter.Decode("255;0;abc")
	assert.Error(t, err, "Expected error for invalid blue value")
	assert.Contains(t, err.Error(), "invalid blue value", "Error message should mention blue value")
}

func TestRangeConverter_Encode(t *testing.T) {
	t.Parallel()

	converter := &RangeConverter{}

	// Тест кодирования
	value := 42
	encoded := converter.Encode(value)
	assert.Equal(t, "42", encoded, "Range encoding failed")
}

func TestRangeConverter_Decode(t *testing.T) {
	t.Parallel()

	converter := &RangeConverter{}

	// Тест декодирования
	encoded := "42"
	decoded, err := converter.Decode(encoded)
	assert.NoError(t, err, "Range decoding failed")
	assert.Equal(t, 42, decoded, "Range decoding result is incorrect")

	// Тест с некорректным значением
	_, err = converter.Decode("abc")
	assert.Error(t, err, "Expected error for invalid range value")
}

func TestSwitchConverter_Encode(t *testing.T) {
	t.Parallel()

	converter := &SwitchConverter{}

	// Тест кодирования
	value := true
	encoded := converter.Encode(value)
	assert.Equal(t, "true", encoded, "Switch encoding failed")

	value = false
	encoded = converter.Encode(value)
	assert.Equal(t, "false", encoded, "Switch encoding failed")
}

func TestSwitchConverter_Decode(t *testing.T) {
	t.Parallel()

	converter := &SwitchConverter{}

	// Тест декодирования
	encoded := "true"
	decoded, err := converter.Decode(encoded)
	assert.NoError(t, err, "Switch decoding failed")
	assert.Equal(t, true, decoded, "Switch decoding result is incorrect")

	encoded = "false"
	decoded, err = converter.Decode(encoded)
	assert.NoError(t, err, "Switch decoding failed")
	assert.Equal(t, false, decoded, "Switch decoding result is incorrect")

	// Тест с некорректным значением
	_, err = converter.Decode("abc")
	assert.Error(t, err, "Expected error for invalid switch value")
}

func TestValueConverter_Encode(t *testing.T) {
	t.Parallel()

	converter := &ValueConverter{}

	// Тест кодирования
	value := 42.42
	encoded := converter.Encode(value)
	assert.Equal(t, "42.42", encoded, "Value encoding failed")
}

func TestValueConverter_Decode(t *testing.T) {
	t.Parallel()

	converter := &ValueConverter{}

	// Тест декодирования
	encoded := "42.42"
	decoded, err := converter.Decode(encoded)
	assert.NoError(t, err, "Value decoding failed")
	assert.Equal(t, 42.42, decoded, "Value decoding result is incorrect")

	// Тест с некорректным значением
	_, err = converter.Decode("abc")
	assert.Error(t, err, "Expected error for invalid value")
}

func TestTimeOnlyConverter_Encode(t *testing.T) {
	t.Parallel()

	converter := &TimeOnlyConverter{}

	// Тест кодирования
	value := timeonly.NewTime(14, 30, 0)
	encoded := converter.Encode(value)
	assert.Equal(t, "14:30:00", encoded, "Time encoding failed")
}

func TestTimeOnlyConverter_Decode(t *testing.T) {
	t.Parallel()

	converter := &TimeOnlyConverter{}

	// Тест декодирования
	encoded := "14:30:00"
	decoded, err := converter.Decode(encoded)
	assert.NoError(t, err, "Time decoding failed")
	assert.Equal(t, timeonly.NewTime(14, 30, 0), decoded, "Time decoding result is incorrect")

	// Тест с некорректным форматом
	_, err = converter.Decode("14:30")
	assert.Error(t, err, "Expected error for invalid time format")

	// Тест с некорректным значением
	_, err = converter.Decode("abc")
	assert.Error(t, err, "Expected error for invalid time value")
}
