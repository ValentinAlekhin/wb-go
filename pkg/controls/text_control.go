package controls

import (
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type TextControl struct {
	control *Control
}

func (c *TextControl) GetValue() string {
	return c.control.GetValue()
}

func (c *TextControl) AddWatcher(f func(payload ControlWatcherPayload)) {
	c.control.AddWatcher(func(p ControlWatcherPayload) {
		f(ControlWatcherPayload{
			NewValue: p.NewValue,
			OldValue: p.OldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *TextControl) GetInfo() ControlInfo {
	return c.control.GetInfo()
}

func (c *TextControl) SetValue(value string) {
	c.control.SetValue(value)
}

func NewTextControl(client *wb.Client, device, control string) *TextControl {
	c := NewControl(client, device, control)
	return &TextControl{control: c}
}
