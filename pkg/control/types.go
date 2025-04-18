package control

import "time"

type ControlInterface interface {
	// GetValue returns the current value of the control as a string.
	GetValue() string

	// SetValue sets a new value for the control.
	SetValue(value string)

	// AddWatcher registers a callback function that will be invoked when
	// the control's value changes.
	AddWatcher(f func(payload WatcherPayloadString))

	// GetInfo returns metadata and topic information about the control.
	GetInfo() Info
}

type Converter[T any] interface {
	// Encode converts a value of type T to its string representation.
	Encode(value T) string

	// Decode parses a string into a value of type T.
	Decode(value string) (T, error)
}

// WatcherPayload represents a payload for watchers with a generic type T.
// It contains the new value, old value, and the topic where the change occurred.
type WatcherPayload[T any] struct {
	NewValue T      // The new value after the change.
	OldValue T      // The old value before the change.
	Topic    string // The MQTT topic where the change was published.
}

// WatcherPayloadInt represents a payload for watchers with the type int.
// It is an alias for WatcherPayload[int].
type WatcherPayloadInt = WatcherPayload[int]

// WatcherPayloadBool represents a payload for watchers with the type bool.
// It is an alias for WatcherPayload[bool].
type WatcherPayloadBool = WatcherPayload[bool]

// WatcherPayloadString represents a payload for watchers with the type string.
// It is an alias for WatcherPayload[string].
type WatcherPayloadString = WatcherPayload[string]

// WatcherPayloadFloat64 represents a payload for watchers with the type float64.
// It is an alias for WatcherPayload[float64].
type WatcherPayloadFloat64 = WatcherPayload[float64]

// WatcherPayloadTime represents a payload for watchers with the type time.Time.
// It is an alias for WatcherPayload[time.Time].
type WatcherPayloadTime = WatcherPayload[time.Time]

// WatcherPayloadRGB represents a payload for watchers with the type RgbValue.
// It is an alias for WatcherPayload[RgbValue].
type WatcherPayloadRGB = WatcherPayload[RgbValue]

// RgbValue represents an RGB color value with red, green, and blue components.
type RgbValue struct {
	Red   int // The red component of the RGB color.
	Green int // The green component of the RGB color.
	Blue  int // The blue component of the RGB color.
}

// Info base control information
type Info struct {
	Name         string
	ValueTopic   string
	CommandTopic string
	Meta         Meta
}

// Meta represents metadata for a control, including type, units, range, and display settings.
type Meta struct {
	Type      string                      `json:"type,omitempty"`      // Type of the control
	Units     string                      `json:"units,omitempty"`     // Units of measurement (only for type="value")
	Max       float64                     `json:"max,omitempty"`       // Maximum value
	Min       float64                     `json:"min,omitempty"`       // Minimum value
	Precision float64                     `json:"precision,omitempty"` // Precision
	Order     int                         `json:"order"`               // Display order
	ReadOnly  bool                        `json:"readonly"`            // Read-only flag
	Title     MultilingualText            `json:"title"`               // Title in multiple languages
	Enum      map[string]MultilingualEnum `json:"enum,omitempty"`      // Enum titles in multiple languages
}

// MultilingualEnum represents an enum title in multiple languages.
type MultilingualEnum struct {
	Title MultilingualText `json:"title"` // Enum title in multiple languages
}

// MultilingualText stores text values in multiple languages.
type MultilingualText map[string]string
