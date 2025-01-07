package virualcontrol

import (
	"encoding/json"
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/atomic"
	"time"
)

type VirtualControl struct {
	name         string
	meta         control.Meta
	client       *wb.Client
	value        atomic.String
	valueTopic   string
	commandTopic string
	metaTopic    string
	addChan      chan func(payload control.ControlWatcherPayload)
	eventChan    chan control.ControlWatcherPayload
	onChan       chan string
	onHandler    OnHandler
	stopChan     chan struct{}
}

type OnHandler func(payload OnHandlerPayload)

type OnHandlerPayload struct {
	Set   func(value string)
	Value string
}

func (c *VirtualControl) GetValue() string {
	return c.value.Load()
}

func (c *VirtualControl) SetValue(value string) {
	oldValue := c.value.Load()
	c.value.Swap(value)

	if oldValue == value {
		return
	}

	payload := control.ControlWatcherPayload{
		NewValue: value,
		OldValue: oldValue,
		Topic:    c.valueTopic,
	}
	c.eventChan <- payload

	c.client.Publish(wb.PublishPayload{
		Topic:    c.valueTopic,
		Value:    value,
		QOS:      1,
		Retained: true,
	})
}

func (c *VirtualControl) GetInfo() control.ControlInfo {
	return control.ControlInfo{
		Name:         c.name,
		ValueTopic:   c.valueTopic,
		CommandTopic: c.commandTopic,
		Meta:         c.meta,
	}
}

func (c *VirtualControl) AddWatcher(f func(payload control.ControlWatcherPayload)) {
	c.addChan <- f
}

func (c *VirtualControl) runWatchHandler() {
	listeners := make([]func(p control.ControlWatcherPayload), 0)

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

func (c *VirtualControl) runOnHandler() {
	for {
		select {
		case newValue := <-c.onChan:
			c.onHandler(OnHandlerPayload{
				Set:   c.SetValue,
				Value: newValue,
			})
		case <-c.stopChan:
			return
		}
	}
}

func (c *VirtualControl) subscribeToOnTopic() {
	callback := func(client mqtt.Client, msg mqtt.Message) {
		c.onChan <- string(msg.Payload())
	}
	c.client.Subscribe(c.commandTopic, callback)
}

func (c *VirtualControl) setMeta() {
	byteMeta, err := json.Marshal(c.meta)
	if err != nil {
		fmt.Println(err)
	}

	c.client.Publish(wb.PublishPayload{
		Topic:    c.metaTopic,
		Value:    string(byteMeta),
		QOS:      1,
		Retained: true,
	})
}

func (c *VirtualControl) loadPrevValue() {
	eventChannel := make(chan struct{})
	done := make(chan bool)
	timeoutDuration := 500 * time.Millisecond

	go func() {
		timer := time.NewTimer(timeoutDuration)
		defer timer.Stop()

		for {
			select {
			case <-eventChannel:
				done <- true
				return

			case <-timer.C:
				done <- true
				return
			}
		}
	}()

	callback := func(client mqtt.Client, msg mqtt.Message) {
		value := string(msg.Payload())
		c.value.Swap(value)
		eventChannel <- struct{}{}
	}
	c.client.Subscribe(c.valueTopic, callback)

	<-done

	c.client.Unsubscribe(c.valueTopic)
}

func NewVirtualControl(client *wb.Client, device, controlName string, meta control.Meta, onHandler OnHandler) *VirtualControl {
	vc := &VirtualControl{
		name:         controlName,
		meta:         meta,
		client:       client,
		valueTopic:   fmt.Sprintf(conventions.CONV_CONTROL_VALUE_FMT, device, controlName),
		commandTopic: fmt.Sprintf(conventions.CONV_CONTROL_ON_VALUE_FMT, device, controlName),
		metaTopic:    fmt.Sprintf(conventions.CONV_CONTROL_META_V2_FMT, device, controlName),
		value:        atomic.String{},
		stopChan:     make(chan struct{}),
		onChan:       make(chan string),
		addChan:      make(chan func(payload control.ControlWatcherPayload)),
		eventChan:    make(chan control.ControlWatcherPayload),
		onHandler:    func(payload OnHandlerPayload) {},
	}

	if onHandler != nil {
		vc.onHandler = onHandler
	}

	vc.loadPrevValue()
	vc.subscribeToOnTopic()
	vc.setMeta()
	go vc.runWatchHandler()
	go vc.runOnHandler()

	return vc
}
