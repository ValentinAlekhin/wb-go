package virualcontrol

import (
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type VirtualTextControl struct {
	control *VirtualControl
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

func (c *VirtualTextControl) AddWatcher(f func(payload control.WatcherPayload)) {
	c.control.AddWatcher(func(p control.WatcherPayload) {
		f(control.WatcherPayload{
			NewValue: p.NewValue,
			OldValue: p.OldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *VirtualTextControl) GetInfo() control.Info {
	return c.control.GetInfo()
}

func NewVirtualTextControl(client wb.ClientInterface, device, control string, meta control.Meta, onTextHandler OnTextHandler) *VirtualTextControl {
	vc := &VirtualTextControl{}
	onHandler := func(payload OnHandlerPayload) {
		newPayload := OnTextHandlerPayload{
			Set:   payload.Set,
			Value: payload.Value,
		}

		if onTextHandler != nil {
			onTextHandler(newPayload)
		}
	}
	meta.Type = "text"
	vc.control = NewVirtualControl(client, device, control, meta, onHandler)
	return vc
}
