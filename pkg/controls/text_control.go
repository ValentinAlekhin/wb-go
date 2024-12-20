package controls

import (
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type TextControl struct {
	control *Control
}

func (tc *TextControl) GetValue() string {
	return tc.control.GetValue()
}

func (tc *TextControl) AddWatcher(f func(payload ControlWatcherPayload)) {
	tc.control.AddWatcher(func(p ControlWatcherPayload) {
		f(ControlWatcherPayload{
			NewValue: p.NewValue,
			OldValue: p.OldValue,
			Topic:    p.Topic,
		})
	})
}

func (tc *TextControl) SetValue(value string) {
	tc.control.SetValue(value)
}

func NewTextControl(client *wb.Client, device, control string) *TextControl {
	c := NewControl(client, device, control)
	return &TextControl{control: c}
}
