package virualcontrol

import (
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"strconv"
)

type VirtualSwitchControl struct {
	control *VirtualControl
}

type OnSwitchHandler func(payload OnSwitchHandlerPayload)

type OnSwitchHandlerPayload struct {
	Set   func(bool)
	Value bool
}

func (c *VirtualSwitchControl) GetValue() bool {
	return c.decode(c.control.GetValue())
}

func (c *VirtualSwitchControl) SetValue(value bool) {
	c.control.SetValue(c.encode(value))
}

func (c *VirtualSwitchControl) Toggle() {
	if c.GetValue() {
		c.SetValue(false)
	} else {
		c.SetValue(true)
	}
}

func (c *VirtualSwitchControl) TurnOff() {
	c.SetValue(false)
}

func (c *VirtualSwitchControl) TurnOn() {
	c.SetValue(true)
}

func (c *VirtualSwitchControl) encode(value bool) string {
	if value {
		return conventions.CONV_SWITCH_VALUE_TRUE
	} else {
		return conventions.CONV_SWITCH_VALUE_FALSE
	}
}

func (c *VirtualSwitchControl) decode(value string) bool {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return v
}

func (c *VirtualSwitchControl) AddWatcher(f func(payload control.SwitchControlWatcherPayload)) {
	c.control.AddWatcher(func(p control.ControlWatcherPayload) {
		f(control.SwitchControlWatcherPayload{
			NewValue: c.decode(p.NewValue),
			OldValue: c.decode(p.OldValue),
			Topic:    p.Topic,
		})
	})
}

func NewVirtualSwitchControl(client *wb.Client, device, control string, meta control.Meta, onSwitchHandler OnSwitchHandler) *VirtualSwitchControl {
	vc := &VirtualSwitchControl{}
	onHandler := func(payload OnHandlerPayload) {
		value := vc.decode(payload.Value)

		newPayload := OnSwitchHandlerPayload{
			Set:   vc.SetValue,
			Value: value,
		}

		if onSwitchHandler != nil {
			onSwitchHandler(newPayload)
		} else {
			vc.SetValue(value)
		}
	}
	meta.Type = "switch"
	vc.control = NewVirtualControl(client, device, control, meta, onHandler)
	return vc
}
