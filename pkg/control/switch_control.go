package control

import (
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type SwitchControl struct {
	converter SwitchConverter
	control   *Control
}

func (c *SwitchControl) GetValue() bool {
	v, _ := c.converter.Decode(c.control.GetValue())
	return v
}

func (c *SwitchControl) AddWatcher(f func(payload WatcherPayloadBool)) {
	c.control.AddWatcher(func(p WatcherPayloadString) {
		newValue, _ := c.converter.Decode(p.NewValue)
		oldValue, _ := c.converter.Decode(p.OldValue)

		f(WatcherPayloadBool{
			NewValue: newValue,
			OldValue: oldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *SwitchControl) SetValue(value bool) {
	c.control.SetValue(c.converter.Encode(value))
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

func (c *SwitchControl) GetInfo() Info {
	return c.control.GetInfo()
}

func NewSwitchControl(client wb.ClientInterface, device, control string, meta Meta) *SwitchControl {
	c := NewControl(client, device, control, meta)
	return &SwitchControl{control: c}
}
