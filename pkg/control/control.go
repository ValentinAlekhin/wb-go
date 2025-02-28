package control

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/atomic"
)

// Control represents a control entity that interacts with MQTT to manage values and notify about changes.
type Control struct {
	name         string
	meta         Meta
	client       wb.ClientInterface
	value        atomic.String
	valueTopic   string
	commandTopic string
	addChan      chan func(payload WatcherPayloadString)
	eventChan    chan WatcherPayloadString
	setChan      chan string
	stopChan     chan struct{}
}

// GetValue returns the current value of the control.
func (c *Control) GetValue() string {
	return c.value.Load()
}

// AddWatcher adds a watcher function that will be called when the control's value changes.
func (c *Control) AddWatcher(f func(payload WatcherPayloadString)) {
	c.addChan <- f
}

// SetValue sets a new value for the control.
func (c *Control) SetValue(value string) {
	c.setChan <- value
}

// GetInfo returns information about the control, including its name, topics, and metadata.
func (c *Control) GetInfo() Info {
	return Info{
		Name:         c.name,
		ValueTopic:   c.valueTopic,
		CommandTopic: c.commandTopic,
		Meta:         c.meta,
	}
}

// publish sends a new value to the MQTT command topic.
func (c *Control) publish(value string) {
	_ = c.client.Publish(wb.PublishPayload{
		Value: value,
		Topic: c.commandTopic,
		QOS:   1,
	})
}

// subscribe subscribes to the MQTT value topic to receive updates.
func (c *Control) subscribe() {
	callback := func(client mqtt.Client, msg mqtt.Message) {
		newValue := string(msg.Payload())
		c.handleValueUpdate(newValue)
	}

	_ = c.client.Subscribe(c.valueTopic, callback)
}

// handleValueUpdate processes a new value received from the MQTT topic.
func (c *Control) handleValueUpdate(value string) {
	oldValue := c.value.Load()
	c.value.Swap(value)

	payload := WatcherPayloadString{
		NewValue: value,
		OldValue: oldValue,
		Topic:    c.valueTopic,
	}

	c.eventChan <- payload
}

// runWatchHandler manages the list of watchers and notifies them about value changes.
func (c *Control) runWatchHandler() {
	listeners := make([]func(p WatcherPayloadString), 0)

	for {
		select {
		case callback := <-c.addChan:
			listeners = append(listeners, callback)
		case event := <-c.eventChan:
			for _, callback := range listeners {
				go callback(event)
			}
		case <-c.stopChan:
			return
		}
	}
}

// runSetValueHandler processes new values and publishes them to the MQTT topic.
func (c *Control) runSetValueHandler() {
	var valueToSet string

	for {
		select {
		case newValue := <-c.setChan:
			if newValue != valueToSet {
				c.publish(newValue)
				valueToSet = newValue
			}
		case <-c.stopChan:
			return
		}
	}
}

// NewControl creates a new Control instance with the specified MQTT client, device, control name, and metadata.
func NewControl(client wb.ClientInterface, device, control string, meta Meta) *Control {
	c := &Control{
		name:         control,
		meta:         meta,
		client:       client,
		valueTopic:   fmt.Sprintf(conventions.CONV_CONTROL_VALUE_FMT, device, control),
		commandTopic: fmt.Sprintf(conventions.CONV_CONTROL_ON_VALUE_FMT, device, control),
		value:        atomic.String{},
		addChan:      make(chan func(payload WatcherPayloadString)),
		eventChan:    make(chan WatcherPayloadString),
		setChan:      make(chan string),
		stopChan:     make(chan struct{}),
	}

	c.value.Store("")
	go c.runWatchHandler()
	go c.runSetValueHandler()
	c.subscribe()

	return c
}
