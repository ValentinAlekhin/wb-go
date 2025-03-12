package control

import (
	"context"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type TextControl struct {
	control *Control
}

func (c *TextControl) GetValue() string {
	return c.control.GetValue()
}

func (c *TextControl) AddWatcher(f func(payload WatcherPayloadString)) {
	c.control.AddWatcher(func(p WatcherPayloadString) {
		f(p)
	})
}

func (c *TextControl) GetInfo() Info {
	return c.control.GetInfo()
}

func (c *TextControl) SetValue(value string) {
	c.control.SetValue(value)
}

func NewTextControl(ctx context.Context, client wb.ClientInterface, device, control string, meta Meta) *TextControl {
	c := NewControl(ctx, client, device, control, meta)
	return &TextControl{control: c}
}
