package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbLed106Controls struct {
	Overcurrent        *controls.SwitchControl
	RgbStrip           *controls.SwitchControl
	RgbPalette         *controls.RgbControl
	RgbStripHue        *controls.RangeControl
	RgbStripSaturation *controls.RangeControl
	RgbStripBrightness *controls.RangeControl
	HueChanging        *controls.SwitchControl
	HueChangingRate    *controls.ValueControl
	Channel4           *controls.SwitchControl
	Channel4Brightness *controls.RangeControl
	BoardTemperature   *controls.ValueControl
	AllowedPower       *controls.ValueControl
	Input1             *controls.SwitchControl
	Input1Counter      *controls.ValueControl
	Input2             *controls.SwitchControl
	Input2Counter      *controls.ValueControl
	Input3             *controls.SwitchControl
	Input3Counter      *controls.ValueControl
	Input4             *controls.SwitchControl
	Input4Counter      *controls.ValueControl
	Serial             *controls.TextControl
}

type WbLed106 struct {
	Name     string
	Address  string
	Controls *WbLed106Controls
}

func (w *WbLed106) GetControlsInfo() []controls.ControlInfo {
	var infoList []controls.ControlInfo

	// Получаем значение и тип структуры Controls
	controlsValue := reflect.ValueOf(w.Controls).Elem()
	controlsType := controlsValue.Type()

	// Проходимся по всем полям структуры Controls
	for i := 0; i < controlsValue.NumField(); i++ {
		field := controlsValue.Field(i)

		// Проверяем, что поле является указателем и не nil
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			// Проверяем, реализует ли поле метод GetInfo
			method := field.MethodByName("GetInfo")
			if method.IsValid() {
				// Вызываем метод GetInfo
				info := method.Call(nil)[0].Interface().(controls.ControlInfo)
				infoList = append(infoList, info)
			} else {
				fmt.Printf("Field %s does not implement GetInfo\n", controlsType.Field(i).Name)
			}
		}
	}

	return infoList
}

var (
	onceWbLed106     sync.Once
	instanceWbLed106 *WbLed106
)

func NewWbLed106(client *mqtt.Client) *WbLed106 {
	onceWbLed106.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "wb-led", "106")
		controlList := &WbLed106Controls{
			Overcurrent:        controls.NewSwitchControl(client, deviceName, "Overcurrent"),
			RgbStrip:           controls.NewSwitchControl(client, deviceName, "RGB Strip"),
			RgbPalette:         controls.NewRgbControl(client, deviceName, "RGB Palette"),
			RgbStripHue:        controls.NewRangeControl(client, deviceName, "RGB Strip Hue"),
			RgbStripSaturation: controls.NewRangeControl(client, deviceName, "RGB Strip Saturation"),
			RgbStripBrightness: controls.NewRangeControl(client, deviceName, "RGB Strip Brightness"),
			HueChanging:        controls.NewSwitchControl(client, deviceName, "Hue Changing"),
			HueChangingRate:    controls.NewValueControl(client, deviceName, "Hue Changing Rate"),
			Channel4:           controls.NewSwitchControl(client, deviceName, "Channel 4"),
			Channel4Brightness: controls.NewRangeControl(client, deviceName, "Channel 4 Brightness"),
			BoardTemperature:   controls.NewValueControl(client, deviceName, "Board Temperature"),
			AllowedPower:       controls.NewValueControl(client, deviceName, "Allowed Power"),
			Input1:             controls.NewSwitchControl(client, deviceName, "Input 1"),
			Input1Counter:      controls.NewValueControl(client, deviceName, "Input 1 Counter"),
			Input2:             controls.NewSwitchControl(client, deviceName, "Input 2"),
			Input2Counter:      controls.NewValueControl(client, deviceName, "Input 2 Counter"),
			Input3:             controls.NewSwitchControl(client, deviceName, "Input 3"),
			Input3Counter:      controls.NewValueControl(client, deviceName, "Input 3 Counter"),
			Input4:             controls.NewSwitchControl(client, deviceName, "Input 4"),
			Input4Counter:      controls.NewValueControl(client, deviceName, "Input 4 Counter"),
			Serial:             controls.NewTextControl(client, deviceName, "Serial"),
		}

		instanceWbLed106 = &WbLed106{
			Name:     deviceName,
			Address:  "106",
			Controls: controlList,
		}
	})

	return instanceWbLed106
}
