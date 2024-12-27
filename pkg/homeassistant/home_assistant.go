package homeassistant

import (
	"encoding/json"
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/iancoleman/strcase"
	"strings"
	"sync"
)

type Discovery struct {
	client    *wb.Client
	baseTopic string
}

func (d *Discovery) AddDevice(info deviceInfo.DeviceInfo) {
	wg := &sync.WaitGroup{}
	wg.Add(len(info.ControlsInfo))

	for _, controlInfo := range info.ControlsInfo {
		go func() {
			defer wg.Done()

			var config MqttDiscoveryConfig
			domain := getAnyDomain(controlInfo)

			switch controlInfo.Name {
			case "RGB Strip":
				config = getRgbLightConfig(info, controlInfo)
				domain = "light"

			case "Channel 1":
				config = getDimLightConfig(info, controlInfo)
				domain = "light"

			case "Channel 2":
				config = getDimLightConfig(info, controlInfo)
				domain = "light"

			case "Channel 3":
				config = getDimLightConfig(info, controlInfo)
				domain = "light"

			case "Channel 4":
				config = getDimLightConfig(info, controlInfo)
				domain = "light"

			case "CCT1":
				config = getCctLightConfig(info, controlInfo)
				domain = "light"

			case "CCT2":
				config = getCctLightConfig(info, controlInfo)
				domain = "light"

			default:
				if d.skipControl(controlInfo.Name) {
					return
				}
				config = getAnyControlConfig(info, controlInfo)
			}

			byteConfig, _ := json.Marshal(config)
			haControlTopic := d.getHaControlTopic(controlInfo.Name)
			topic := fmt.Sprintf("%s/%s/%s/%s/config", d.baseTopic, domain, info.Name, haControlTopic)

			d.client.Publish(wb.PublishPayload{
				Topic:    topic,
				Value:    string(byteConfig),
				Retained: true,
				QOS:      2,
			})
		}()
	}

	wg.Wait()
	fmt.Printf("Устроство %s добавлено в Home Assistant\n", info.Name)
}

func (d *Discovery) skipControl(name string) bool {
	subStrings := []string{"rgb", "cct", "channel"}

	lowName := strings.ToLower(name)

	for _, subString := range subStrings {
		if strings.Contains(lowName, subString) {
			return true
		}
	}

	return false
}

func (d *Discovery) getHaControlTopic(name string) string {
	topic := strcase.ToSnake(name)
	topic = clearControlName(topic)

	return topic
}

func NewDiscovery(baseTopic string, client *wb.Client) *Discovery {
	return &Discovery{baseTopic: baseTopic, client: client}
}
