package virualcontrol

import (
	"context"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
)

type VirtualRangeControl struct {
	converter control.RangeConverter
	control   *VirtualControl
}

type RangeHandler = OnHandler[int]
type RangeHandlerPayload = OnHandlerPayload[int]

type RangeOptions struct {
	BaseOptions
	OnHandler    RangeHandler
	DefaultValue int
}

func (c *VirtualRangeControl) GetValue() int {
	value, _ := c.converter.Decode(c.control.GetValue()) // Обработать ошибку
	return value
}

func (c *VirtualRangeControl) SetValue(v int) {
	c.control.SetValue(c.converter.Encode(v))
}

func (c *VirtualRangeControl) AddWatcher(f func(payload control.WatcherPayloadInt)) {
	c.control.AddWatcher(func(p control.WatcherPayloadString) {
		newValue, _ := c.converter.Decode(p.NewValue) // Обработать ошибку
		oldValue, _ := c.converter.Decode(p.OldValue) // Обработать ошибку

		f(control.WatcherPayloadInt{
			NewValue: newValue,
			OldValue: oldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *VirtualRangeControl) GetInfo() control.Info {
	return c.control.GetInfo()
}

func NewVirtualRangeControl(ctx context.Context, opt RangeOptions) *VirtualRangeControl {
	vc := &VirtualRangeControl{}
	onHandler := func(payload OnHandlerPayload[string]) {
		value, err := vc.converter.Decode(payload.Value)
		if err != nil {
			// Логировать ошибку
			return
		}

		newPayload := RangeHandlerPayload{
			Set:   vc.SetValue,
			Value: value,
		}

		if opt.OnHandler != nil {
			opt.OnHandler(newPayload)
		}
	}
	opt.Meta.Type = "range"

	vOpt := Options{
		BaseOptions:  opt.BaseOptions,
		OnHandler:    onHandler,
		DefaultValue: vc.converter.Encode(opt.DefaultValue),
	}

	vc.control = NewVirtualControl(ctx, vOpt)
	return vc
}
