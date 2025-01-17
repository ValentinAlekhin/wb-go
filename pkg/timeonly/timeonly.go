package timeonly

import (
	"fmt"
	"time"
)

// Time represents time in seconds since the start of the day
type Time struct {
	seconds int
}

// NewTime creates a new Time object
func NewTime(hour, minute, second int) Time {
	return Time{seconds: hour*3600 + minute*60 + second}
}

// FromSeconds creates a Time object from the number of seconds since the start of the day
func FromSeconds(seconds int) Time {
	return Time{seconds: seconds % 86400} // Ensure the value does not exceed 24 hours
}

// FromTime creates a Time object from a time.Time object
func FromTime(t time.Time) Time {
	return NewTime(t.Hour(), t.Minute(), t.Second())
}

// Now returns the current time as a Time object
func Now() Time {
	return FromTime(time.Now())
}

// Hour returns the number of hours
func (t Time) Hour() int {
	return t.seconds / 3600
}

// Minute returns the number of minutes
func (t Time) Minute() int {
	return (t.seconds % 3600) / 60
}

// Second returns the number of seconds (0 to 59)
func (t Time) Second() int {
	return t.seconds % 60
}

// Equal checks if two Time objects are equal
func (t Time) Equal(other Time) bool {
	return t.seconds == other.seconds
}

// Before checks if the current time is earlier than another
func (t Time) Before(other Time) bool {
	return t.seconds < other.seconds
}

// After checks if the current time is later than another
func (t Time) After(other Time) bool {
	return t.seconds > other.seconds
}

// IsZero checks if the object represents midnight
func (t Time) IsZero() bool {
	return t.seconds == 0
}

// ToTime converts Time to a time.Time object with a fixed date
func (t Time) ToTime() time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
}

// ParseString parses a string in the HH:MM:SS format and returns a Time object
func ParseString(s string) (Time, error) {
	t, err := time.Parse(time.TimeOnly, s)
	if err != nil {
		return Time{}, err
	}

	return FromTime(t), nil
}

// String returns the string representation of the time in HH:MM:SS format
func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
}
