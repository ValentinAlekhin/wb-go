package homeassistant

import (
	"encoding/json"
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iancoleman/strcase"
	"strings"
	"sync"
	"time"
)

func (d *Discovery) AddDevice(info basedevice.Info) {
	d.AddDeviceWithMiddleware(info, nil)
}

func (d *Discovery) AddDeviceWithMiddleware(info basedevice.Info, middleware ConfigMiddleware) {
	wg := &sync.WaitGroup{}
	wg.Add(len(info.ControlsInfo))

	for _, controlInfo := range info.ControlsInfo {
		go func() {
			defer wg.Done()

			config, domain, ignore := GetConfigAndDomain(info, controlInfo)
			if ignore {
				return
			}

			configPointer := &config
			domainPointer := &domain
			if middleware != nil {
				middleware(domainPointer, configPointer, info, controlInfo)
			}

			byteConfig, _ := json.Marshal(configPointer)
			haControlTopic := d.getHaControlTopic(controlInfo.Name)
			baseTopic := fmt.Sprintf("%s/%s/%s/%s", d.prefix, *domainPointer, info.Name, haControlTopic)
			configTopic := baseTopic + "/config"
			metaTopic := fmt.Sprintf("%s/%s", baseTopic, DiscoveryMetaTopic)
			meta := DiscoveryMeta{
				ClientName: d.name,
				CreatedAt:  time.Now(),
			}
			jsonMeta, _ := json.Marshal(meta)

			_ = d.client.Publish(wb.PublishPayload{
				Topic:    configTopic,
				Value:    string(byteConfig),
				Retained: true,
				QOS:      2,
			})

			_ = d.client.Publish(wb.PublishPayload{
				Topic:    metaTopic,
				Value:    string(jsonMeta),
				Retained: true,
				QOS:      2,
			})
		}()
	}

	wg.Wait()
	fmt.Printf("Устроство %s добавлено в Home Assistant\n", info.Name)
}

// Clear all devices, added by wb-go
func (d *Discovery) Clear() {
	eventChannel := make(chan struct{})
	timeoutDuration := 1 * time.Second
	done := make(chan bool)

	go func() {
		timer := time.NewTimer(timeoutDuration)
		defer timer.Stop()

		for {
			select {
			case <-eventChannel:
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(timeoutDuration)

			case <-timer.C:
				done <- true
				return
			}
		}
	}()

	result := map[string]DiscoveryMeta{}

	topicName := fmt.Sprintf("%s/+/+/+/%s", d.prefix, DiscoveryMetaTopic)
	handler := func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := msg.Payload()
		var meta DiscoveryMeta
		err := json.Unmarshal(payload, &meta)
		if err != nil {
			return
		}

		result[topic] = meta
		eventChannel <- struct{}{}
	}
	_ = d.client.Subscribe(topicName, handler)

	<-done

	_ = d.client.Unsubscribe(topicName)

	for topic, meta := range result {
		if meta.ClientName != d.name {
			continue
		}

		configTopic := strings.Replace(topic, DiscoveryMetaTopic, "config", 1)
		_ = d.client.Publish(wb.PublishPayload{
			Topic:    configTopic,
			Value:    "",
			QOS:      1,
			Retained: true,
		})

		_ = d.client.Publish(wb.PublishPayload{
			Topic:    topic,
			Value:    "",
			QOS:      1,
			Retained: true,
		})
	}
}

func (d *Discovery) getHaControlTopic(name string) string {
	topic := strcase.ToSnake(name)
	topic = clearControlName(topic)

	return topic
}

func (d *Discovery) applyDefaultOptions() {
	if d.name == "" {
		d.name = DefaultDiscoveryName
	}
}

func NewDiscovery(opt DiscoveryOptions) *Discovery {

	discovery := &Discovery{
		prefix: opt.Prefix,
		client: opt.Client,
		name:   opt.Name,
	}
	discovery.applyDefaultOptions()

	return discovery
}
