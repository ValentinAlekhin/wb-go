package virualcontrol

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/timeonly"
)

type VirtualTimeControl struct {
	control *VirtualControl
}

type TimeOptions struct {
	BaseOptions
	OnHandler    OnTimeHandler
	DefaultValue timeonly.Time
}

type OnTimeHandler func(payload OnTimeHandlerPayload)

type OnTimeHandlerPayload struct {
	Set   func(timeonly.Time)
	Value timeonly.Time
	Error error
}

type TimeControlWatcherPayload struct {
	NewValue timeonly.Time
	OldValue timeonly.Time
	Topic    string
}

func (c *VirtualTimeControl) GetValue() timeonly.Time {
	value := c.control.GetValue()
	timeValue, _ := c.decode(value)
	return timeValue
}

func (c *VirtualTimeControl) SetValue(v timeonly.Time) {
	c.control.SetValue(c.encode(v))
}

func (c *VirtualTimeControl) decode(value string) (timeonly.Time, error) {
	t, err := timeonly.ParseString(value)
	if err != nil {
		return timeonly.Time{}, fmt.Errorf("invalid time format: %s", value)
	}
	return t, nil
}

func (c *VirtualTimeControl) encode(value timeonly.Time) string {
	return value.String()
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

	vOpt := Options{
		BaseOptions:  opt.BaseOptions,
		OnHandler:    onHandler,
		DefaultValue: vc.encode(opt.DefaultValue),
	}

	vc.control = NewVirtualControl(vOpt)
	return vc
}
