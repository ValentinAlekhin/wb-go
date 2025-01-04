package control

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

func (c *RangeControl) GetValue() int {
	return c.decode(c.control.GetValue())
}

func (c *RangeControl) AddWatcher(f func(payload RangeControlWatcherPayload)) {
	c.control.AddWatcher(func(p ControlWatcherPayload) {
		f(RangeControlWatcherPayload{
			NewValue: c.decode(p.NewValue),
			OldValue: c.decode(p.OldValue),
			Topic:    p.Topic,
		})
	})
}

func (c *RangeControl) SetValue(value int) {
	c.control.SetValue(c.encode(value))
}

func (c *RangeControl) GetInfo() ControlInfo {
	return c.control.GetInfo()
}

func (c *RangeControl) encode(value int) string {
	return strconv.Itoa(value)
}

func (c *RangeControl) decode(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return v
}

func NewRangeControl(client *wb.Client, device, control string, meta Meta) *RangeControl {
	c := NewControl(client, device, control, meta)
	return &RangeControl{control: c}
}
