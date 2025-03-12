package virualcontrol

import (
	"context"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/timeonly"
)

type VirtualTimeControl struct {
	converter control.TimeOnlyConverter
	control   *VirtualControl
}

type TimeHandler = OnHandler[timeonly.Time]
type TimeHandlerPayload = OnHandlerPayload[timeonly.Time]

type TimeOptions struct {
	BaseOptions
	OnHandler    TimeHandler
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
	timeValue, _ := c.converter.Decode(value)
	return timeValue
}

func (c *VirtualTimeControl) SetValue(v timeonly.Time) {
	c.control.SetValue(c.converter.Encode(v))
}

func (c *VirtualTimeControl) AddWatcher(f func(payload TimeControlWatcherPayload)) {
	c.control.AddWatcher(func(p control.WatcherPayloadString) {
		newValue, _ := c.converter.Decode(p.NewValue)
		oldValue, _ := c.converter.Decode(p.OldValue)

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

func NewVirtualTimeControl(ctx context.Context, opt TimeOptions) *VirtualTimeControl {
	vc := &VirtualTimeControl{}
	onHandler := func(payload OnHandlerPayload[string]) {
		value, err := vc.converter.Decode(payload.Value)

		newPayload := TimeHandlerPayload{
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
		DefaultValue: vc.converter.Encode(opt.DefaultValue),
	}

	vc.control = NewVirtualControl(ctx, vOpt)
	return vc
}
