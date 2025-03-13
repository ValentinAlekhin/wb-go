package virualcontrol

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ValentinAlekhin/wb-go/internal/db"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/atomic"
)

type VirtualControl struct {
	name         string
	meta         control.Meta
	queries      *db.Queries
	client       wb.ClientInterface
	value        atomic.String
	valueTopic   string
	commandTopic string
	metaTopic    string
	addChan      chan func(payload control.WatcherPayloadString)
	eventChan    chan control.WatcherPayloadString
	onChan       chan string
	onHandler    OnHandler[string]
	closed       atomic.Bool
}

func (c *VirtualControl) GetValue() string {
	return c.value.Load()
}

func (c *VirtualControl) SetValue(value string) {
	if c.closed.Load() {
		return
	}

	oldValue := c.value.Load()
	c.value.Swap(value)

	if oldValue == value {
		return
	}

	ctx := context.TODO()
	params := db.UpdateVirtualControlParams{
		Value: value,
		Topic: c.valueTopic,
	}
	_, err := c.queries.UpdateVirtualControl(ctx, params)
	if err != nil {
		fmt.Printf("Error updating virtual control: %v\n", err)
	}

	payload := control.WatcherPayloadString{
		NewValue: value,
		OldValue: oldValue,
		Topic:    c.valueTopic,
	}
	c.eventChan <- payload

	_ = c.client.Publish(wb.PublishPayload{
		Topic:    c.valueTopic,
		Value:    value,
		QOS:      1,
		Retained: true,
	})
}

func (c *VirtualControl) GetInfo() control.Info {
	return control.Info{
		Name:         c.name,
		ValueTopic:   c.valueTopic,
		CommandTopic: c.commandTopic,
		Meta:         c.meta,
	}
}

func (c *VirtualControl) AddWatcher(f func(payload control.WatcherPayloadString)) {
	if c.closed.Load() {
		return
	}

	c.addChan <- f
}

func (c *VirtualControl) runWatchHandler(ctx context.Context) {
	listeners := make([]func(p control.WatcherPayloadString), 0)
	defer c.close()

	for {
		select {
		case callback := <-c.addChan:
			listeners = append(listeners, callback)
		case event := <-c.eventChan:
			for _, callback := range listeners {
				go callback(event)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *VirtualControl) runOnHandler(ctx context.Context) {
	defer c.close()

	for {
		select {
		case newValue := <-c.onChan:
			c.onHandler(OnHandlerPayload[string]{
				Set:   c.SetValue,
				Value: newValue,
			})
		case <-ctx.Done():
			return
		}
	}
}

func (c *VirtualControl) subscribeToOnTopic() {
	callback := func(client mqtt.Client, msg mqtt.Message) {
		c.onChan <- string(msg.Payload())
	}
	_ = c.client.Subscribe(c.commandTopic, callback)
}

func (c *VirtualControl) setMeta() {
	byteMeta, err := json.Marshal(c.meta)
	if err != nil {
		fmt.Println(err)
	}

	_ = c.client.Publish(wb.PublishPayload{
		Topic:    c.metaTopic,
		Value:    string(byteMeta),
		QOS:      1,
		Retained: true,
	})
}

func (c *VirtualControl) loadPrevValue(defaultValue string) {
	ctx := context.TODO()

	virtualControl, err := c.queries.GetVirtualControl(ctx, c.valueTopic)
	if err != nil {
		params := db.CreateVirtualControlParams{
			Topic: c.valueTopic,
			Value: defaultValue,
		}
		virtualControl, err = c.queries.CreateVirtualControl(ctx, params)
		if err != nil {
			fmt.Printf("Insert control '%s' value error: %s\n", c.valueTopic, err)
		}
	}

	value := virtualControl.Value
	if value == "" {
		value = defaultValue
	}

	c.value.Swap(value)

	_ = c.client.Publish(wb.PublishPayload{
		Topic:    c.valueTopic,
		Value:    value,
		QOS:      1,
		Retained: true,
	})
}

func (c *VirtualControl) close() {
	if c.closed.CompareAndSwap(false, true) {
		_ = c.client.Unsubscribe(c.commandTopic)

		close(c.onChan)
		close(c.addChan)
		close(c.eventChan)
	}
}

func NewVirtualControl(ctx context.Context, opt Options) *VirtualControl {
	vc := &VirtualControl{
		name:         opt.Name,
		meta:         opt.Meta,
		queries:      opt.Queries,
		client:       opt.Client,
		valueTopic:   fmt.Sprintf(conventions.CONV_CONTROL_VALUE_FMT, opt.Device, opt.Name),
		commandTopic: fmt.Sprintf(conventions.CONV_CONTROL_ON_VALUE_FMT, opt.Device, opt.Name),
		metaTopic:    fmt.Sprintf(conventions.CONV_CONTROL_META_V2_FMT, opt.Device, opt.Name),
		value:        atomic.String{},
		onChan:       make(chan string),
		addChan:      make(chan func(payload control.WatcherPayloadString)),
		eventChan:    make(chan control.WatcherPayloadString),
		onHandler:    func(payload OnHandlerPayload[string]) {},
	}

	if opt.OnHandler != nil {
		vc.onHandler = opt.OnHandler
	}

	go vc.runWatchHandler(ctx)

	vc.loadPrevValue(opt.DefaultValue)
	vc.subscribeToOnTopic()
	vc.setMeta()

	go vc.runOnHandler(ctx)

	return vc
}
