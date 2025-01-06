package control

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/atomic"
)

type Control struct {
	name         string
	meta         Meta
	client       *wb.Client
	value        atomic.String
	valueTopic   string
	commandTopic string
	addChan      chan func(payload ControlWatcherPayload)
	eventChan    chan ControlWatcherPayload
	setChan      chan string
	stopChan     chan struct{}
}

type Meta struct {
	Type      string                      `json:"type,omitempty"`      // Тип контроля
	Units     string                      `json:"units,omitempty"`     // Единицы измерения (только для type="value")
	Max       float64                     `json:"max,omitempty"`       // Максимальное значение
	Min       float64                     `json:"min,omitempty"`       // Минимальное значение
	Precision float64                     `json:"precision,omitempty"` // Точность
	Order     int                         `json:"order"`               // Порядок отображения
	ReadOnly  bool                        `json:"readonly"`            // Только для чтения
	Title     MultilingualText            `json:"title"`               // Название (разные языки)
	Enum      map[string]MultilingualEnum `json:"enum,omitempty"`      // Заголовки для enum
}

type MultilingualEnum struct {
	Title MultilingualText `json:"title"` // Название enum на разных языках
}

// MultilingualText хранит текстовые значения на разных языках
type MultilingualText map[string]string

type ControlWatcherPayload struct {
	NewValue string
	OldValue string
	Topic    string
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

func (c *Control) GetInfo() ControlInfo {
	return ControlInfo{
		Name:         c.name,
		ValueTopic:   c.valueTopic,
		CommandTopic: c.commandTopic,
		Meta:         c.meta,
	}
}

func (c *Control) publish(value string) {
	c.client.Publish(wb.PublishPayload{
		Value: value,
		Topic: c.commandTopic,
		QOS:   1,
	})
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

func NewControl(client *wb.Client, device, control string, meta Meta) *Control {
	c := &Control{
		name:         control,
		meta:         meta,
		client:       client,
		valueTopic:   fmt.Sprintf(conventions.CONV_CONTROL_VALUE_FMT, device, control),
		commandTopic: fmt.Sprintf(conventions.CONV_CONTROL_ON_VALUE_FMT, device, control),
		value:        atomic.String{},
		addChan:      make(chan func(payload ControlWatcherPayload)),
		eventChan:    make(chan ControlWatcherPayload),
		setChan:      make(chan string),
		stopChan:     make(chan struct{}),
	}

	c.value.Store("")
	go c.runWatchHandler()
	go c.runSetValueHandler()
	c.subscribe()

	return c
}
