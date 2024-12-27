package controls

import (
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"strconv"
)

type SwitchControl struct {
	control *Control
}

type SwitchControlWatcherPayload struct {
	NewValue    bool
	OldValue    bool
	Topic       string
	ControlName string
}

func (c *SwitchControl) GetValue() bool {
	return c.decode(c.control.GetValue())
}

func (c *SwitchControl) AddWatcher(f func(payload SwitchControlWatcherPayload)) {
	c.control.AddWatcher(func(p ControlWatcherPayload) {
		f(SwitchControlWatcherPayload{
			NewValue: c.decode(p.NewValue),
			OldValue: c.decode(p.OldValue),
			Topic:    p.Topic,
		})
	})
}

func (c *SwitchControl) SetValue(value bool) {
	c.control.SetValue(c.encode(value))
}

func (c *SwitchControl) Toggle() {
	if c.GetValue() {
		c.SetValue(false)
	} else {
		c.SetValue(true)
	}
}

func (c *SwitchControl) TurnOff() {
	c.SetValue(false)
}

func (c *SwitchControl) TurnOn() {
	c.SetValue(true)
}

func (c *SwitchControl) GetInfo() ControlInfo {
	return c.control.GetInfo()
}

func (c *SwitchControl) encode(value bool) string {
	if value {
		return conventions.CONV_SWITCH_VALUE_TRUE
	} else {
		return conventions.CONV_SWITCH_VALUE_FALSE
	}
}

func (c *SwitchControl) decode(value string) bool {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return v
}

func NewSwitchControl(client *wb.Client, device, control string, meta Meta) *SwitchControl {
	c := NewControl(client, device, control, meta)
	return &SwitchControl{control: c}
}
