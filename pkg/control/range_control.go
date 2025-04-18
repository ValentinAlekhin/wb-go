package control

import (
	"context"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type RangeControl struct {
	converter Converter[int]
	control   ControlInterface
}

func (c *RangeControl) GetValue() int {
	val, _ := c.converter.Decode(c.control.GetValue())
	return val
}

func (c *RangeControl) AddWatcher(f func(payload WatcherPayloadInt)) {
	c.control.AddWatcher(func(p WatcherPayloadString) {
		newValue, _ := c.converter.Decode(p.NewValue)
		oldValue, _ := c.converter.Decode(p.OldValue)

		f(WatcherPayloadInt{
			NewValue: newValue,
			OldValue: oldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *RangeControl) SetValue(value int) {
	c.control.SetValue(c.converter.Encode(value))
}

func (c *RangeControl) GetInfo() Info {
	return c.control.GetInfo()
}

func NewRangeControl(ctx context.Context, client wb.ClientInterface, device, control string, meta Meta) *RangeControl {
	c := NewControl(ctx, client, device, control, meta)
	return &RangeControl{control: c}
}
