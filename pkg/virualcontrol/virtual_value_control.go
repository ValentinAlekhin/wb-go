package virualcontrol

import (
	"context"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
)

type VirtualValueControl struct {
	converter control.ValueConverter
	control   *VirtualControl
}

type ValueHandler = OnHandler[float64]
type ValueHandlerPayload = OnHandlerPayload[float64]

type ValueOptions struct {
	BaseOptions
	OnHandler    ValueHandler
	DefaultValue float64
}

func (c *VirtualValueControl) GetValue() float64 {
	value, _ := c.converter.Decode(c.control.GetValue())
	return value
}

func (c *VirtualValueControl) SetValue(v float64) {
	c.control.SetValue(c.converter.Encode(v))
}

func (c *VirtualValueControl) AddWatcher(f func(payload control.WatcherPayloadFloat64)) {
	c.control.AddWatcher(func(p control.WatcherPayloadString) {
		newValue, _ := c.converter.Decode(p.NewValue)
		oldValue, _ := c.converter.Decode(p.OldValue)

		f(control.WatcherPayloadFloat64{
			NewValue: newValue,
			OldValue: oldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *VirtualValueControl) GetInfo() control.Info {
	return c.control.GetInfo()
}

func NewVirtualValueControl(ctx context.Context, opt ValueOptions) *VirtualValueControl {
	vc := &VirtualValueControl{}
	onHandler := func(payload OnHandlerPayload[string]) {
		value, _ := vc.converter.Decode(payload.Value)

		newPayload := ValueHandlerPayload{
			Set:   vc.SetValue,
			Value: value,
		}

		if opt.OnHandler != nil {
			opt.OnHandler(newPayload)
		}
	}
	opt.Meta.Type = "value"

	vOpt := Options{BaseOptions: opt.BaseOptions, OnHandler: onHandler, DefaultValue: vc.converter.Encode(opt.DefaultValue)}

	vc.control = NewVirtualControl(ctx, vOpt)
	return vc
}
