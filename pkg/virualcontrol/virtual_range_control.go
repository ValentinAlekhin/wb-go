package virualcontrol

import (
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"strconv"
)

type VirtualRangeControl struct {
	control *VirtualControl
}

type RangeOptions struct {
	BaseOptions
	OnHandler    OnRangeHandler
	DefaultValue int
}

type OnRangeHandler func(payload OnRangeHandlerPayload)

type OnRangeHandlerPayload struct {
	Set   func(value int)
	Value int
}

func (c *VirtualRangeControl) GetValue() int {
	return c.decode(c.control.GetValue())
}

func (c *VirtualRangeControl) SetValue(v int) {
	c.control.SetValue(c.encode(v))
}

func (c *VirtualRangeControl) encode(value int) string {
	return strconv.Itoa(value)
}

func (c *VirtualRangeControl) decode(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return v
}

func (c *VirtualRangeControl) AddWatcher(f func(payload control.RangeControlWatcherPayload)) {
	c.control.AddWatcher(func(p control.WatcherPayload) {
		f(control.RangeControlWatcherPayload{
			NewValue: c.decode(p.NewValue),
			OldValue: c.decode(p.OldValue),
			Topic:    p.Topic,
		})
	})
}

func (c *VirtualRangeControl) GetInfo() control.Info {
	return c.control.GetInfo()
}

func NewVirtualRangeControl(opt RangeOptions) *VirtualRangeControl {
	vc := &VirtualRangeControl{}
	onHandler := func(payload OnHandlerPayload) {
		value := vc.decode(payload.Value)

		newPayload := OnRangeHandlerPayload{
			Set:   vc.SetValue,
			Value: value,
		}

		if opt.OnHandler != nil {
			opt.OnHandler(newPayload)
		}
	}
	opt.Meta.Type = "range"

	vOpt := Options{BaseOptions: opt.BaseOptions, OnHandler: onHandler, DefaultValue: vc.encode(opt.DefaultValue)}

	vc.control = NewVirtualControl(vOpt)
	return vc
}
