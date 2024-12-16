package devices

import (
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"strconv"
)

type RangeControl struct {
	control *Control
}

type RangeControlWatcherPayload struct {
	NewValue    int
	OldValue    int
	Topic       string
	ControlName string
}

func (rc *RangeControl) GetValue() int {
	value, err := strconv.Atoi(rc.control.GetValue())
	if err != nil {
		return 0
	}

	return value
}

func (rc *RangeControl) AddWatcher(f func(payload RangeControlWatcherPayload)) {
	rc.control.AddWatcher(func(p ControlWatcherPayload) {
		f(RangeControlWatcherPayload{
			NewValue:    rc.decode(p.NewValue),
			OldValue:    rc.decode(p.OldValue),
			Topic:       p.Topic,
			ControlName: p.ControlName,
		})
	})
}

func (rc *RangeControl) SetValue(value int) {
	rc.control.SetValue(rc.encode(value))
}

func (rc *RangeControl) encode(value int) string {
	return strconv.Itoa(value)
}

func (rc *RangeControl) decode(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return v
}

func NewRangeControl(client *wb.Client, topic string) *RangeControl {
	control := NewControl(client, topic)
	return &RangeControl{control: control}
}
