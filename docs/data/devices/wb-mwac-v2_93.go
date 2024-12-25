package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbMwacV293Controls struct {
	P1Volume       *controls.ValueControl
	P2Volume       *controls.ValueControl
	InputF1        *controls.SwitchControl
	InputF1Counter *controls.ValueControl
	InputF2        *controls.SwitchControl
	InputF2Counter *controls.ValueControl
	InputF3        *controls.SwitchControl
	InputF3Counter *controls.ValueControl
	InputF4        *controls.SwitchControl
	InputF4Counter *controls.ValueControl
	InputF5        *controls.SwitchControl
	InputF5Counter *controls.ValueControl
	InputS6        *controls.SwitchControl
	InputS6Counter *controls.ValueControl
	OutputK1       *controls.SwitchControl
	OutputK2       *controls.SwitchControl
	LeakageMode    *controls.SwitchControl
	CleaningMode   *controls.SwitchControl
	Serial         *controls.TextControl
}

type WbMwacV293 struct {
	Name     string
	Address  string
	Controls *WbMwacV293Controls
}

func (w *WbMwacV293) GetControlsInfo() []controls.ControlInfo {
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
	onceWbMwacV293     sync.Once
	instanceWbMwacV293 *WbMwacV293
)

func NewWbMwacV293(client *mqtt.Client) *WbMwacV293 {
	onceWbMwacV293.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "wb-mwac-v2", "93")
		controlList := &WbMwacV293Controls{
			P1Volume:       controls.NewValueControl(client, deviceName, "P1 Volume"),
			P2Volume:       controls.NewValueControl(client, deviceName, "P2 Volume"),
			InputF1:        controls.NewSwitchControl(client, deviceName, "Input F1"),
			InputF1Counter: controls.NewValueControl(client, deviceName, "Input F1 Counter"),
			InputF2:        controls.NewSwitchControl(client, deviceName, "Input F2"),
			InputF2Counter: controls.NewValueControl(client, deviceName, "Input F2 Counter"),
			InputF3:        controls.NewSwitchControl(client, deviceName, "Input F3"),
			InputF3Counter: controls.NewValueControl(client, deviceName, "Input F3 Counter"),
			InputF4:        controls.NewSwitchControl(client, deviceName, "Input F4"),
			InputF4Counter: controls.NewValueControl(client, deviceName, "Input F4 Counter"),
			InputF5:        controls.NewSwitchControl(client, deviceName, "Input F5"),
			InputF5Counter: controls.NewValueControl(client, deviceName, "Input F5 Counter"),
			InputS6:        controls.NewSwitchControl(client, deviceName, "Input S6"),
			InputS6Counter: controls.NewValueControl(client, deviceName, "Input S6 Counter"),
			OutputK1:       controls.NewSwitchControl(client, deviceName, "Output K1"),
			OutputK2:       controls.NewSwitchControl(client, deviceName, "Output K2"),
			LeakageMode:    controls.NewSwitchControl(client, deviceName, "Leakage Mode"),
			CleaningMode:   controls.NewSwitchControl(client, deviceName, "Cleaning Mode"),
			Serial:         controls.NewTextControl(client, deviceName, "Serial"),
		}

		instanceWbMwacV293 = &WbMwacV293{
			Name:     deviceName,
			Address:  "93",
			Controls: controlList,
		}
	})

	return instanceWbMwacV293
}
