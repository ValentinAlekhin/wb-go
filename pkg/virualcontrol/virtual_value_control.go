package virualcontrol

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"strconv"
	"strings"
)

type VirtualValueControl struct {
	control *VirtualControl
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
	c.control.AddWatcher(func(p control.ControlWatcherPayload) {
		f(control.ValueControlWatcherPayload{
			NewValue: c.decode(p.NewValue),
			OldValue: c.decode(p.OldValue),
			Topic:    p.Topic,
		})
	})
}

func NewVirtualValueControl(client *wb.Client, device, control string, meta control.Meta, onValueHandler OnValueHandler) *VirtualValueControl {
	vc := &VirtualValueControl{}
	onHandler := func(payload OnHandlerPayload) {
		value := vc.decode(payload.Value)

		newPayload := OnValueHandlerPayload{
			Set:   vc.SetValue,
			Value: value,
		}

		if onValueHandler != nil {
			onValueHandler(newPayload)
		} else {
			vc.SetValue(value)
		}
	}
	meta.Type = "value"
	vc.control = NewVirtualControl(client, device, control, meta, onHandler)
	return vc
}
