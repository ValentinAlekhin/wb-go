package timeonly

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTimeOnly(t *testing.T) {
	// Test creating Time from hours, minutes, and seconds
	tm := NewTime(14, 30, 45)
	assert.Equal(t, 14, tm.Hour(), "Hour should be 14")
	assert.Equal(t, 30, tm.Minute(), "Minute should be 30")
	assert.Equal(t, 45, tm.Second(), "Second should be 45")
}

func TestFromSeconds(t *testing.T) {
	// Test creating Time from seconds since the start of the day
	tm := FromSeconds(52245) // 14:30:45
	assert.Equal(t, 14, tm.Hour(), "Hour should be 14")
	assert.Equal(t, 30, tm.Minute(), "Minute should be 30")
	assert.Equal(t, 45, tm.Second(), "Second should be 45")

	// Test wrap-around for more than 24 hours
	tm = FromSeconds(90000) // 25:00:00 -> 01:00:00
	assert.Equal(t, 1, tm.Hour(), "Hour should be 1")
	assert.Equal(t, 0, tm.Minute(), "Minute should be 0")
	assert.Equal(t, 0, tm.Second(), "Second should be 0")
}

func TestFromTime(t *testing.T) {
	// Test creating Time from time.Time
	timeObj := time.Date(2025, 1, 1, 14, 30, 45, 0, time.UTC)
	tm := FromTime(timeObj)
	assert.Equal(t, 14, tm.Hour(), "Hour should match")
	assert.Equal(t, 30, tm.Minute(), "Minute should match")
	assert.Equal(t, 45, tm.Second(), "Second should match")
}

func TestNow(t *testing.T) {
	// Test Now function to get the current time
	now := time.Now()
	tm := Now()
	assert.Equal(t, now.Hour(), tm.Hour(), "Hour should match current time")
	assert.Equal(t, now.Minute(), tm.Minute(), "Minute should match current time")
	assert.Equal(t, now.Second(), tm.Second(), "Second should match current time")
}

func TestHourMinuteSecond(t *testing.T) {
	// Test extracting hours, minutes, and seconds
	tm := NewTime(14, 30, 45)
	assert.Equal(t, 14, tm.Hour(), "Hour should be 14")
	assert.Equal(t, 30, tm.Minute(), "Minute should be 30")
	assert.Equal(t, 45, tm.Second(), "Second should be 45")
}

func TestEqual(t *testing.T) {
	// Test equality of two Time objects
	tm1 := NewTime(14, 30, 45)
	tm2 := NewTime(14, 30, 45)
	tm3 := NewTime(15, 0, 0)

	assert.True(t, tm1.Equal(tm2), "tm1 should equal tm2")
	assert.False(t, tm1.Equal(tm3), "tm1 should not equal tm3")
}

func TestBefore(t *testing.T) {
	// Test if one time is before another
	tm1 := NewTime(14, 30, 0)
	tm2 := NewTime(16, 0, 0)

	assert.True(t, tm1.Before(tm2), "tm1 should be before tm2")
	assert.False(t, tm2.Before(tm1), "tm2 should not be before tm1")
}

func TestAfter(t *testing.T) {
	// Test if one time is after another
	tm1 := NewTime(14, 30, 0)
	tm2 := NewTime(16, 0, 0)

	assert.True(t, tm2.After(tm1), "tm2 should be after tm1")
	assert.False(t, tm1.After(tm2), "tm1 should not be after tm2")
}

func TestIsZero(t *testing.T) {
	// Test if a Time object represents midnight
	tm := NewTime(0, 0, 0)
	assert.True(t, tm.IsZero(), "tm should represent midnight")

	nonZeroTm := NewTime(1, 0, 0)
	assert.False(t, nonZeroTm.IsZero(), "nonZeroTm should not represent midnight")
}

func TestToTime(t *testing.T) {
	// Test converting Time to time.Time
	tm := NewTime(14, 30, 45)
	convertedTime := tm.ToTime()

	assert.Equal(t, 14, convertedTime.Hour(), "Hour should match")
	assert.Equal(t, 30, convertedTime.Minute(), "Minute should match")
	assert.Equal(t, 45, convertedTime.Second(), "Second should match")
}

func TestParseString(t *testing.T) {
	// Test parsing a valid time string
	input := "14:30:45"
	tm, err := ParseString(input)
	assert.NoError(t, err, "Parsing should not produce an error")
	assert.Equal(t, 14, tm.Hour(), "Hour should match")
	assert.Equal(t, 30, tm.Minute(), "Minute should match")
	assert.Equal(t, 45, tm.Second(), "Second should match")

	// Test parsing an invalid time string
	invalidInput := "invalid"
	_, err = ParseString(invalidInput)
	assert.Error(t, err, "Parsing invalid input should produce an error")
}

func TestString(t *testing.T) {
	// Test converting Time to a string
	tm := NewTime(14, 30, 45)
	assert.Equal(t, "14:30:45", tm.String(), "String representation should match")
}
