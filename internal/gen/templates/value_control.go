package devices

import (
	"strconv"
	"strings"
	wb "wb-go/pkg/mqtt"
)

type ValueControl struct {
	control *Control
}

type ValueControlWatcherPayload struct {
	NewValue    float64
	OldValue    float64
	Topic       string
	ControlName string
}

func (c *ValueControl) GetValue() float64 {
	return c.decode(c.control.GetValue())
}

func (c *ValueControl) AddWatcher(f func(payload ValueControlWatcherPayload)) {
	c.control.AddWatcher(func(p ControlWatcherPayload) {
		f(ValueControlWatcherPayload{
			NewValue:    c.decode(p.NewValue),
			OldValue:    c.decode(p.OldValue),
			Topic:       p.Topic,
			ControlName: p.ControlName,
		})
	})
}

func (c *ValueControl) decode(value string) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return 0
	}

	return v
}

func NewValueControl(client *wb.Client, topic string) *ValueControl {
	control := NewControl(client, topic)
	return &ValueControl{control: control}
}
