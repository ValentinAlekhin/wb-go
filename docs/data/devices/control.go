package devices

import (
	"fmt"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/atomic"
)

type ControlWatcherPayload struct {
	NewValue    string
	OldValue    string
	Topic       string
	ControlName string
}

type Control struct {
	client    *wb.Client
	value     atomic.String
	topic     string
	addChan   chan func(payload ControlWatcherPayload) // Канал для добавления новых подписчиков
	eventChan chan ControlWatcherPayload
	setChan   chan string
	stopChan  chan struct{}
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
	commandTopic := fmt.Sprintf("%s/on", c.topic)
	c.client.Publish(commandTopic, value)
}

func (c *Control) subscribe() {
	callback := func(client mqtt.Client, msg mqtt.Message) {
		newValue := string(msg.Payload())
		c.handleValueUpdate(newValue)
	}

	c.client.Subscribe(c.topic, callback)
}

func (c *Control) handleValueUpdate(value string) {
	oldValue := c.value.Load()
	c.value.Swap(value)

	payload := ControlWatcherPayload{
		NewValue:    value,
		OldValue:    oldValue,
		Topic:       c.topic,
		ControlName: "",
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

func NewControl(client *wb.Client, topic string) *Control {
	sw := &Control{
		client:    client,
		topic:     topic,
		value:     atomic.String{},
		addChan:   make(chan func(payload ControlWatcherPayload)),
		eventChan: make(chan ControlWatcherPayload),
		setChan:   make(chan string),
		stopChan:  make(chan struct{}),
	}

	sw.value.Store("")
	go sw.runWatchHandler()
	go sw.runSetValueHandler()
	sw.subscribe()

	return sw
}
