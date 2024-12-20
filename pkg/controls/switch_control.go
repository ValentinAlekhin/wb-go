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

func (sw *SwitchControl) GetValue() bool {
	return sw.decode(sw.control.GetValue())
}

func (sw *SwitchControl) AddWatcher(f func(payload SwitchControlWatcherPayload)) {
	sw.control.AddWatcher(func(p ControlWatcherPayload) {
		f(SwitchControlWatcherPayload{
			NewValue: sw.decode(p.NewValue),
			OldValue: sw.decode(p.OldValue),
			Topic:    p.Topic,
		})
	})
}
func (sw *SwitchControl) SetValue(value bool) {
	sw.control.SetValue(sw.encode(value))
}

func (sw *SwitchControl) Toggle() {
	if sw.GetValue() {
		sw.SetValue(false)
	} else {
		sw.SetValue(true)
	}
}

func (sw *SwitchControl) TurnOff() {
	sw.SetValue(false)
}

func (sw *SwitchControl) TurnOn() {
	sw.SetValue(true)
}

func (sw *SwitchControl) encode(value bool) string {
	if value {
		return conventions.CONV_SWITCH_VALUE_TRUE
	} else {
		return conventions.CONV_SWITCH_VALUE_FALSE
	}
}

func (sw *SwitchControl) decode(value string) bool {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return v
}

func NewSwitchControl(client *wb.Client, device, control string) *SwitchControl {
	c := NewControl(client, device, control)
	return &SwitchControl{control: c}
}
