package virualcontrol

import (
	"context"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
)

type VirtualTextControl struct {
	control *VirtualControl
}

type TextHandler = OnHandler[string]
type TextHandlerPayload = OnHandlerPayload[string]

type TextOptions struct {
	BaseOptions
	OnHandler    TextHandler
	DefaultValue string
}

type OnTextHandler func(payload OnTextHandlerPayload)

type OnTextHandlerPayload struct {
	Set   func(string)
	Value string
}

func (c *VirtualTextControl) GetValue() string {
	return c.control.GetValue()
}

func (c *VirtualTextControl) SetValue(v string) {
	c.control.SetValue(v)
}

func (c *VirtualTextControl) AddWatcher(f func(payload control.WatcherPayloadString)) {
	c.control.AddWatcher(func(p control.WatcherPayloadString) {
		f(control.WatcherPayloadString{
			NewValue: p.NewValue,
			OldValue: p.OldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *VirtualTextControl) GetInfo() control.Info {
	return c.control.GetInfo()
}

func NewVirtualTextControl(ctx context.Context, opt TextOptions) *VirtualTextControl {
	vc := &VirtualTextControl{}
	onHandler := func(payload TextHandlerPayload) {
		newPayload := TextHandlerPayload{
			Set:   payload.Set,
			Value: payload.Value,
		}
		if opt.OnHandler != nil {
			opt.OnHandler(newPayload)
		}
	}
	opt.Meta.Type = "text"

	vOpt := Options{BaseOptions: opt.BaseOptions, OnHandler: onHandler, DefaultValue: opt.DefaultValue}

	vc.control = NewVirtualControl(ctx, vOpt)
	return vc
}
