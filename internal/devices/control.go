package devices

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type Control struct {
	publish   func(topic string, value string)
	value     ControlValue
	Topic     string
	valueLock sync.RWMutex
	watchers  []func(payload ControlWatcherPayload)
}

type ControlValue struct {
	Raw   string
	Bool  bool
	Int   int
	Float float64
}

type ControlWatcherPayload struct {
	NewValue    ControlValue
	OldValue    ControlValue
	Topic       string
	ControlName string
}

func (c *Control) handleValueUpdate(value string) {
	defer c.valueLock.Unlock()

	c.valueLock.Lock()
	oldValue := c.value
	newValue := c.parseValue(value)

	c.value = newValue

	payload := ControlWatcherPayload{
		NewValue:    newValue,
		OldValue:    oldValue,
		Topic:       c.Topic,
		ControlName: "",
	}

	for _, handler := range c.watchers {
		go handler(payload)
	}
}

func (c *Control) parseValue(value string) ControlValue {
	boolValue, _ := strconv.ParseBool(value)
	intValue, _ := strconv.Atoi(value)
	floatValue, _ := strconv.ParseFloat(strings.TrimSpace(value), 64)

	return ControlValue{
		Raw:   value,
		Bool:  boolValue,
		Int:   intValue,
		Float: floatValue,
	}
}

func (c *Control) GetValue() ControlValue {
	defer c.valueLock.RUnlock()
	c.valueLock.RLock()
	return c.value
}

func (c *Control) AddWatcher(f func(payload ControlWatcherPayload)) {
	c.watchers = append(c.watchers, f)
}

func (c *Control) SetValue(value string) {
	commandTopic := fmt.Sprintf("%s/on", c.Topic)
	c.publish(commandTopic, value)
}

func (c *Control) SetBoolValue(value bool) {
	finalValue := "0"
	if value {
		finalValue = "1"
	}

	c.SetValue(finalValue)
}

func (c *Control) SetIntValue(value int) {
	c.SetValue(strconv.Itoa(value))
}

func (c *Control) SetFloatValue(value float64) {
	c.SetValue(strconv.FormatFloat(value, 'f', 6, 64))
}

func (c *Control) Toggle() {
	if c.value.Bool {
		c.SetBoolValue(false)
	} else {
		c.SetBoolValue(true)
	}
}
