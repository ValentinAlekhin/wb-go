package control

import (
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type TextControl struct {
	control *Control
}

func (c *TextControl) GetValue() string {
	return c.control.GetValue()
}

func (c *TextControl) AddWatcher(f func(payload WatcherPayload)) {
	c.control.AddWatcher(func(p WatcherPayload) {
		f(WatcherPayload{
			NewValue: p.NewValue,
			OldValue: p.OldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *TextControl) GetInfo() Info {
	return c.control.GetInfo()
}

func (c *TextControl) SetValue(value string) {
	c.control.SetValue(value)
}

func NewTextControl(client wb.ClientInterface, device, control string, meta Meta) *TextControl {
	c := NewControl(client, device, control, meta)
	return &TextControl{control: c}
}
