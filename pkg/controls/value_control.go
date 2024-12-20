package controls

import (
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"strconv"
	"strings"
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
			NewValue: c.decode(p.NewValue),
			OldValue: c.decode(p.OldValue),
			Topic:    p.Topic,
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

func NewValueControl(client *wb.Client, device, control string) *ValueControl {
	c := NewControl(client, device, control)
	return &ValueControl{c}
}
