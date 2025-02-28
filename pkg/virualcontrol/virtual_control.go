package virualcontrol

import (
	"encoding/json"
	"fmt"
	"github.com/ValentinAlekhin/wb-go/internal/db"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/atomic"
	"gorm.io/gorm"
)

type VirtualControl struct {
	name         string
	meta         control.Meta
	db           *gorm.DB
	client       wb.ClientInterface
	value        atomic.String
	valueTopic   string
	commandTopic string
	metaTopic    string
	addChan      chan func(payload control.WatcherPayloadString)
	eventChan    chan control.WatcherPayloadString
	onChan       chan string
	onHandler    OnHandler[string]
	stopChan     chan struct{}
}

type Options struct {
	BaseOptions
	OnHandler    OnHandler[string]
	DefaultValue string
}

type OnHandler[T comparable] func(payload OnHandlerPayload[T])

type OnHandlerPayload[T comparable] struct {
	Set   func(value T)
	Value T
	Error error
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

	c.db.Model(&db.ControlModel{}).Where("topic = ?", c.valueTopic).Update("value", value)

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
	c.addChan <- f
}

func (c *VirtualControl) runWatchHandler() {
	listeners := make([]func(p control.WatcherPayloadString), 0)

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
			c.onHandler(OnHandlerPayload[string]{
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
	model := db.ControlModel{Topic: c.valueTopic}
	result := c.db.First(&model)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	if result.RowsAffected == 0 {
		model.Value = defaultValue
		c.db.Create(&model)
	}

	c.value.Swap(model.Value)

	_ = c.client.Publish(wb.PublishPayload{
		Topic:    c.valueTopic,
		Value:    model.Value,
		QOS:      1,
		Retained: true,
	})
}

func NewVirtualControl(opt Options) *VirtualControl {
	vc := &VirtualControl{
		name:         opt.Name,
		meta:         opt.Meta,
		db:           opt.DB,
		client:       opt.Client,
		valueTopic:   fmt.Sprintf(conventions.CONV_CONTROL_VALUE_FMT, opt.Device, opt.Name),
		commandTopic: fmt.Sprintf(conventions.CONV_CONTROL_ON_VALUE_FMT, opt.Device, opt.Name),
		metaTopic:    fmt.Sprintf(conventions.CONV_CONTROL_META_V2_FMT, opt.Device, opt.Name),
		value:        atomic.String{},
		stopChan:     make(chan struct{}),
		onChan:       make(chan string),
		addChan:      make(chan func(payload control.WatcherPayloadString)),
		eventChan:    make(chan control.WatcherPayloadString),
		onHandler:    func(payload OnHandlerPayload[string]) {},
	}

	if opt.OnHandler != nil {
		vc.onHandler = opt.OnHandler
	}

	go vc.runWatchHandler()

	vc.loadPrevValue(opt.DefaultValue)
	vc.subscribeToOnTopic()
	vc.setMeta()

	go vc.runOnHandler()

	return vc
}
