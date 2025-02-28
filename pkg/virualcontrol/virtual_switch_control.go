package virualcontrol

import (
	"github.com/ValentinAlekhin/wb-go/pkg/control"
)

type VirtualSwitchControl struct {
	converter control.SwitchConverter
	control   *VirtualControl
}

type SwitchHandler = OnHandler[bool]
type SwitchHandlerPayload = OnHandlerPayload[bool]

type SwitchOptions struct {
	BaseOptions
	OnHandler    SwitchHandler
	DefaultValue bool
}

func (c *VirtualSwitchControl) GetValue() bool {
	value, _ := c.converter.Decode(c.control.GetValue())
	return value
}

func (c *VirtualSwitchControl) SetValue(value bool) {
	c.control.SetValue(c.converter.Encode(value))
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

func (c *VirtualSwitchControl) AddWatcher(f func(payload control.WatcherPayloadBool)) {
	c.control.AddWatcher(func(p control.WatcherPayloadString) {
		newValue, _ := c.converter.Decode(p.NewValue)
		oldValue, _ := c.converter.Decode(p.OldValue)

		f(control.WatcherPayloadBool{
			NewValue: newValue,
			OldValue: oldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *VirtualSwitchControl) GetInfo() control.Info {
	return c.control.GetInfo()
}

func NewVirtualSwitchControl(opt SwitchOptions) *VirtualSwitchControl {
	vc := &VirtualSwitchControl{}
	onHandler := func(payload OnHandlerPayload[string]) {
		value, _ := vc.converter.Decode(payload.Value)

		newPayload := SwitchHandlerPayload{
			Set:   vc.SetValue,
			Value: value,
		}

		if opt.OnHandler != nil {
			opt.OnHandler(newPayload)
		}
	}
	opt.Meta.Type = "switch"

	vOpt := Options{BaseOptions: opt.BaseOptions, OnHandler: onHandler, DefaultValue: vc.converter.Encode(opt.DefaultValue)}

	vc.control = NewVirtualControl(vOpt)
	return vc
}
