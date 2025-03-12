package control

import (
	"context"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type ValueControl struct {
	converter ValueConverter
	control   *Control
}

func (c *ValueControl) GetValue() float64 {
	v, _ := c.converter.Decode(c.control.GetValue())
	return v
}

func (c *ValueControl) AddWatcher(f func(payload WatcherPayloadFloat64)) {
	c.control.AddWatcher(func(p WatcherPayloadString) {
		newValue, _ := c.converter.Decode(p.NewValue)
		oldValue, _ := c.converter.Decode(p.OldValue)

		f(WatcherPayloadFloat64{
			NewValue: newValue,
			OldValue: oldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *ValueControl) GetInfo() Info {
	return c.control.GetInfo()
}

func NewValueControl(ctx context.Context, client wb.ClientInterface, device, control string, meta Meta) *ValueControl {
	c := NewControl(ctx, client, device, control, meta)
	return &ValueControl{control: c}
}
