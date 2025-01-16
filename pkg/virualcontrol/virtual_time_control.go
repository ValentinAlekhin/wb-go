package virualcontrol

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/dromara/carbon/v2"
	"time"
)

type VirtualTimeControl struct {
	control *VirtualControl
}

type TimeOptions struct {
	BaseOptions
	OnHandler    OnTimeHandler
	DefaultValue time.Time
}

type OnTimeHandler func(payload OnTimeHandlerPayload)

type OnTimeHandlerPayload struct {
	Set   func(time.Time)
	Value time.Time
	Error error
}

type TimeControlWatcherPayload struct {
	NewValue time.Time
	OldValue time.Time
	Topic    string
}

func (c *VirtualTimeControl) GetValue() time.Time {
	value := c.control.GetValue()
	timeValue, _ := c.decode(value)
	return timeValue
}

func (c *VirtualTimeControl) SetValue(v time.Time) {
	c.control.SetValue(c.encode(v))
}

func (c *VirtualTimeControl) decode(value string) (time.Time, error) {
	ca := carbon.ParseByFormat(value, "H:i:s")
	valid := ca.IsValid()
	if !valid {
		return time.Time{}, fmt.Errorf("неправильный формат времени: %s", value)
	}

	return ca.StdTime(), nil
}

func (c *VirtualTimeControl) encode(value time.Time) string {
	return carbon.CreateFromStdTime(value).ToTimeString()
}

func (c *VirtualTimeControl) AddWatcher(f func(payload TimeControlWatcherPayload)) {
	c.control.AddWatcher(func(p control.WatcherPayload) {
		newValue, _ := c.decode(p.NewValue)
		oldValue, _ := c.decode(p.OldValue)

		f(TimeControlWatcherPayload{
			NewValue: newValue,
			OldValue: oldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *VirtualTimeControl) GetInfo() control.Info {
	return c.control.GetInfo()
}

func NewVirtualTimeControl(opt TimeOptions) *VirtualTimeControl {
	vc := &VirtualTimeControl{}
	onHandler := func(payload OnHandlerPayload) {
		value, err := vc.decode(payload.Value)

		newPayload := OnTimeHandlerPayload{
			Set:   vc.SetValue,
			Value: value,
			Error: err,
		}

		if opt.OnHandler != nil {
			opt.OnHandler(newPayload)
		}
	}
	opt.Meta.Type = "text"

	vOpt := Options{BaseOptions: opt.BaseOptions, OnHandler: onHandler, DefaultValue: vc.encode(opt.DefaultValue)}

	vc.control = NewVirtualControl(vOpt)
	return vc
}
