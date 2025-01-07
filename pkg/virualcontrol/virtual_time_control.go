package virualcontrol

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/dromara/carbon/v2"
	"time"
)

type VirtualTimeControl struct {
	control *VirtualControl
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
	c.control.AddWatcher(func(p control.ControlWatcherPayload) {
		newValue, _ := c.decode(p.NewValue)
		oldValue, _ := c.decode(p.OldValue)

		f(TimeControlWatcherPayload{
			NewValue: newValue,
			OldValue: oldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *VirtualTimeControl) GetInfo() control.ControlInfo {
	return c.control.GetInfo()
}

func NewVirtualTimeValueControl(client *wb.Client, device, control string, meta control.Meta, onTimeHandler OnTimeHandler) *VirtualTimeControl {
	vc := &VirtualTimeControl{}
	onHandler := func(payload OnHandlerPayload) {
		value, err := vc.decode(payload.Value)

		newPayload := OnTimeHandlerPayload{
			Set:   vc.SetValue,
			Value: value,
			Error: err,
		}

		if onTimeHandler != nil {
			onTimeHandler(newPayload)
		}
	}
	meta.Type = "text"
	vc.control = NewVirtualControl(client, device, control, meta, onHandler)
	return vc
}
