package controls

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/atomic"
)

type ControlWatcherPayload struct {
	NewValue string
	OldValue string
	Topic    string
}

type Control struct {
	client       *wb.Client
	value        atomic.String
	valueTopic   string
	commandTopic string
	addChan      chan func(payload ControlWatcherPayload) // Канал для добавления новых подписчиков
	eventChan    chan ControlWatcherPayload
	setChan      chan string
	stopChan     chan struct{}
}

func (c *Control) GetValue() string {
	return c.value.Load()
}

func (c *Control) AddWatcher(f func(payload ControlWatcherPayload)) {
	c.addChan <- f
}

func (c *Control) SetValue(value string) {
	c.setChan <- value
}

func (c *Control) publish(value string) {
	c.client.Publish(c.commandTopic, value)
}

func (c *Control) subscribe() {
	callback := func(client mqtt.Client, msg mqtt.Message) {
		newValue := string(msg.Payload())
		c.handleValueUpdate(newValue)
	}

	c.client.Subscribe(c.valueTopic, callback)
}

func (c *Control) handleValueUpdate(value string) {
	oldValue := c.value.Load()
	c.value.Swap(value)

	payload := ControlWatcherPayload{
		NewValue: value,
		OldValue: oldValue,
		Topic:    c.valueTopic,
	}

	c.eventChan <- payload
}

func (c *Control) runWatchHandler() {
	listeners := make([]func(p ControlWatcherPayload), 0)

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

func NewControl(client *wb.Client, device, control string) *Control {
	sw := &Control{
		client:       client,
		valueTopic:   fmt.Sprintf(conventions.CONV_CONTROL_VALUE_FMT, device, control),
		commandTopic: fmt.Sprintf(conventions.CONV_CONTROL_ON_VALUE_FMT, device, control),
		value:        atomic.String{},
		addChan:      make(chan func(payload ControlWatcherPayload)),
		eventChan:    make(chan ControlWatcherPayload),
		setChan:      make(chan string),
		stopChan:     make(chan struct{}),
	}

	sw.value.Store("")
	go sw.runWatchHandler()
	go sw.runSetValueHandler()
	sw.subscribe()

	return sw
}
