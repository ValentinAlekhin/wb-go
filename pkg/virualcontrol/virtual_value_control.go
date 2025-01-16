package virualcontrol

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"strconv"
	"strings"
)

type VirtualValueControl struct {
	control *VirtualControl
}

type ValueOptions struct {
	BaseOptions
	OnHandler    OnValueHandler
	DefaultValue float64
}

type OnValueHandler func(payload OnValueHandlerPayload)

type OnValueHandlerPayload struct {
	Set   func(float64)
	Value float64
}

func (c *VirtualValueControl) GetValue() float64 {
	return c.decode(c.control.GetValue())
}

func (c *VirtualValueControl) SetValue(v float64) {
	c.control.SetValue(c.encode(v))
}

func (c *VirtualValueControl) decode(value string) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return 0
	}

	return v
}

func (c *VirtualValueControl) encode(value float64) string {
	return fmt.Sprintf("%f", value)
}

func (c *VirtualValueControl) AddWatcher(f func(payload control.ValueControlWatcherPayload)) {
	c.control.AddWatcher(func(p control.WatcherPayload) {
		f(control.ValueControlWatcherPayload{
			NewValue: c.decode(p.NewValue),
			OldValue: c.decode(p.OldValue),
			Topic:    p.Topic,
		})
	})
}

func (c *VirtualValueControl) GetInfo() control.Info {
	return c.control.GetInfo()
}

func NewVirtualValueControl(opt ValueOptions) *VirtualValueControl {
	vc := &VirtualValueControl{}
	onHandler := func(payload OnHandlerPayload) {
		value := vc.decode(payload.Value)

		newPayload := OnValueHandlerPayload{
			Set:   vc.SetValue,
			Value: value,
		}

		if opt.OnHandler != nil {
			opt.OnHandler(newPayload)
		}
	}
	opt.Meta.Type = "value"

	vOpt := Options{BaseOptions: opt.BaseOptions, OnHandler: onHandler, DefaultValue: vc.encode(opt.DefaultValue)}

	vc.control = NewVirtualControl(vOpt)
	return vc
}
