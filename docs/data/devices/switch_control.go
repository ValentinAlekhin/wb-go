package devices

import (
	"strconv"
	wb "wb-go/pkg/mqtt"
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
			NewValue:    sw.decode(p.NewValue),
			OldValue:    sw.decode(p.OldValue),
			Topic:       p.Topic,
			ControlName: p.ControlName,
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
		return "1"
	} else {
		return "0"
	}
}

func (sw *SwitchControl) decode(value string) bool {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return v
}

func NewSwitchControl(client *wb.Client, topic string) *SwitchControl {
	control := NewControl(client, topic)
	return &SwitchControl{control: control}
}
