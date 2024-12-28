package homeassistant

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/iancoleman/strcase"
	"regexp"
	"strings"
)

var DeviceConfigMap = map[string]*DeviceConfig{
	"wb-led": {
		ignoreRegexpStr: []string{"Brightness", "Temperature", "RGB Palette"},
		getters: []*ConfigGetter{
			{
				regexpStr: `^CCT\d+$`,
				getter:    getWbLedCctConfig,
				domain:    "light",
			},
			{
				regexpStr: `^Channels? [0-9_]{1,3}`,
				getter:    getWbLedDimConfig,
				domain:    "light",
			},
			{
				regexpStr: `^RGB Strip`,
				getter:    getWbLEdRgbConfig,
				domain:    "light",
			},
		},
	},
	"wb-mdm3": {
		ignoreRegexpStr: []string{"^Channel [0-9]"},
		getters: []*ConfigGetter{
			{
				regexpStr: `^K\d+$`,
				getter:    getWbMdm3Config,
				domain:    "light",
			},
		},
	},
}

type DeviceConfig struct {
	getters         []*ConfigGetter
	ignoreRegexpStr []string
	ignoreRegexp    []*regexp.Regexp
}

type ConfigGetter struct {
	regexpStr string
	regexp    *regexp.Regexp
	getter    ConfigGetterFn
	domain    string
}

type ConfigGetterFn func(deviceInfo deviceInfo.DeviceInfo, controlInfo controls.ControlInfo) MqttDiscoveryConfig

func init() {
	for _, devConfig := range DeviceConfigMap {
		devConfig.ignoreRegexp = make([]*regexp.Regexp, len(devConfig.ignoreRegexpStr))

		for i, regexpSrt := range devConfig.ignoreRegexpStr {
			devConfig.ignoreRegexp[i] = regexp.MustCompile(regexpSrt)
		}

		for _, handler := range devConfig.getters {
			handler.regexp = regexp.MustCompile(handler.regexpStr)
		}
	}
}

func getConfigAndDomain(deviceInfo deviceInfo.DeviceInfo, controlInfo controls.ControlInfo) (config MqttDiscoveryConfig, domain string, ignore bool) {
	config = getAnyControlConfig(deviceInfo, controlInfo)
	domain = getAnyDomain(controlInfo)
	ignore = false

	for dev, devConfig := range DeviceConfigMap {
		if dev != deviceInfo.Device {
			continue
		}

		for _, ignoreRegexp := range devConfig.ignoreRegexp {
			if ignoreRegexp.MatchString(controlInfo.Name) {
				ignore = true
				return
			}
		}

		for _, getter := range devConfig.getters {
			if !getter.regexp.MatchString(controlInfo.Name) {
				continue
			}

			config = getter.getter(deviceInfo, controlInfo)
			domain = getter.domain
			break
		}
	}

	return
}

func getAnyDomain(info controls.ControlInfo) string {
	domain := "sensor"

	switch info.Meta.Type {
	case "rgb":
		domain = "light"
	case "range":
		domain = "number"
	case "switch":
		if info.Meta.ReadOnly {
			domain = "binary_sensor"
		} else {
			domain = "switch"
		}
	case "pushbutton":
		domain = "button"
	}

	return domain
}

func getAnyControlConfig(deviceInfo deviceInfo.DeviceInfo, controlInfo controls.ControlInfo) MqttDiscoveryConfig {
	id := getControlId(deviceInfo.Name, controlInfo.Name)

	wbUnit := controlInfo.Meta.Units
	if wbUnit == "" {
		wbUnit, _ = ConvertMetaTypeToUnit(controlInfo.Meta.Type)
	}

	haUnit, _ := ConvertWBUnitToHA(wbUnit)
	class, _ := ConvertUnitToDeviceClass(wbUnit)

	config := GetConfig(MqttDiscoveryConfig{
		Device:            getDevice(deviceInfo),
		Name:              controlInfo.Name,
		UniqueID:          id,
		ObjectID:          id,
		StateTopic:        fmt.Sprintf("/devices/%s/controls/%s", deviceInfo.Name, controlInfo.Name),
		CommandTopic:      fmt.Sprintf("/devices/%s/controls/%s/on", deviceInfo.Name, controlInfo.Name),
		UnitOfMeasurement: haUnit,
		DeviceClass:       class,
	})

	if controlInfo.Meta.Max != 0 || controlInfo.Meta.Min != 0 {
		config.Min = controlInfo.Meta.Min
		config.Max = controlInfo.Meta.Max
	}

	return config
}

func getWbLEdRgbConfig(deviceInfo deviceInfo.DeviceInfo, controlInfo controls.ControlInfo) MqttDiscoveryConfig {
	id := getControlId(deviceInfo.Name, controlInfo.Name)

	return GetConfig(MqttDiscoveryConfig{
		Device:                 getDevice(deviceInfo),
		RGBStateTopic:          fmt.Sprintf("/devices/%s/controls/RGB Palette", deviceInfo.Name),
		RGBCommandTopic:        fmt.Sprintf("/devices/%s/controls/RGB Palette/on", deviceInfo.Name),
		BrightnessStateTopic:   fmt.Sprintf("/devices/%s/controls/RGB Strip Brightness", deviceInfo.Name),
		BrightnessCommandTopic: fmt.Sprintf("/devices/%s/controls/RGB Strip Brightness/on", deviceInfo.Name),
		RGBValueTemplate:       "{{ value.split(';') | join(',') }}",
		RGBCommandTemplate:     "{{ red }};{{ green }};{{ blue }}",
		Name:                   controlInfo.Name,
		UniqueID:               id,
		ObjectID:               id,
		StateTopic:             fmt.Sprintf("/devices/%s/controls/RGB Strip", deviceInfo.Name),
		CommandTopic:           fmt.Sprintf("/devices/%s/controls/RGB Strip/on", deviceInfo.Name),
	})

}

func getWbLedCctConfig(deviceInfo deviceInfo.DeviceInfo, controlInfo controls.ControlInfo) MqttDiscoveryConfig {
	id := getControlId(deviceInfo.Name, controlInfo.Name)

	return GetConfig(MqttDiscoveryConfig{
		Device:                   getDevice(deviceInfo),
		BrightnessStateTopic:     fmt.Sprintf("/devices/%s/controls/%s Brightness", deviceInfo.Name, controlInfo.Name),
		BrightnessCommandTopic:   fmt.Sprintf("/devices/%s/controls/%s Brightness/on", deviceInfo.Name, controlInfo.Name),
		ColorTempStateTopic:      fmt.Sprintf("/devices/%s/controls/%s Temperature", deviceInfo.Name, controlInfo.Name),
		ColorTempCommandTopic:    fmt.Sprintf("/devices/%s/controls/%s Temperature/on", deviceInfo.Name, controlInfo.Name),
		ColorTempValueTemplate:   "{{ ((((100 - value | float) / 100) * (this.attributes.max_mireds - this.attributes.min_mireds)) + this.attributes.min_mireds) | round(0) }}",
		ColorTempCommandTemplate: "{{ (100 - (((value - this.attributes.min_mireds) / (this.attributes.max_mireds - this.attributes.min_mireds)) * 100)) | round(0) }}",
		MaxMireds:                454,
		MinMireds:                154,
		Name:                     controlInfo.Name,
		UniqueID:                 id,
		ObjectID:                 id,
		StateTopic:               fmt.Sprintf("/devices/%s/controls/%s", deviceInfo.Name, controlInfo.Name),
		CommandTopic:             fmt.Sprintf("/devices/%s/controls/%s/on", deviceInfo.Name, controlInfo.Name),
	})

}

func getWbLedDimConfig(deviceInfo deviceInfo.DeviceInfo, controlInfo controls.ControlInfo) MqttDiscoveryConfig {
	id := getControlId(deviceInfo.Name, controlInfo.Name)

	return GetConfig(MqttDiscoveryConfig{
		Device:                 getDevice(deviceInfo),
		BrightnessStateTopic:   fmt.Sprintf("/devices/%s/controls/%s Brightness", deviceInfo.Name, controlInfo.Name),
		BrightnessCommandTopic: fmt.Sprintf("/devices/%s/controls/%s Brightness/on", deviceInfo.Name, controlInfo.Name),
		Name:                   controlInfo.Name,
		UniqueID:               id,
		ObjectID:               id,
		StateTopic:             fmt.Sprintf("/devices/%s/controls/%s", deviceInfo.Name, controlInfo.Name),
		CommandTopic:           fmt.Sprintf("/devices/%s/controls/%s/on", deviceInfo.Name, controlInfo.Name),
	})
}

func getWbMdm3Config(deviceInfo deviceInfo.DeviceInfo, controlInfo controls.ControlInfo) MqttDiscoveryConfig {
	id := getControlId(deviceInfo.Name, controlInfo.Name)
	channelNumber := strings.TrimPrefix(controlInfo.Name, "K")

	return GetConfig(MqttDiscoveryConfig{
		Device:                 getDevice(deviceInfo),
		BrightnessStateTopic:   fmt.Sprintf("/devices/%s/controls/Channel %s", deviceInfo.Name, channelNumber),
		BrightnessCommandTopic: fmt.Sprintf("/devices/%s/controls/Channel %s/on", deviceInfo.Name, channelNumber),
		Name:                   controlInfo.Name,
		UniqueID:               id,
		ObjectID:               id,
		StateTopic:             fmt.Sprintf("/devices/%s/controls/%s", deviceInfo.Name, controlInfo.Name),
		CommandTopic:           fmt.Sprintf("/devices/%s/controls/%s/on", deviceInfo.Name, controlInfo.Name),
	})
}

func getDevice(deviceInfo deviceInfo.DeviceInfo) MqttDiscoveryDevice {
	return MqttDiscoveryDevice{
		Identifiers: deviceInfo.Name,
		Model:       deviceInfo.Device,
		Name:        deviceInfo.Name,
	}
}

func getControlId(device string, control string) string {
	c := strcase.ToSnake(control)
	c = clearControlName(c)
	return fmt.Sprintf("%s_%s", device, c)
}

func clearControlName(name string) string {
	replaceMap := map[string]string{"(": "", ")": ""}

	for k, v := range replaceMap {
		name = strings.ReplaceAll(name, k, v)
	}

	return name
}
