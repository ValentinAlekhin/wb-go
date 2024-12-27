package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbMdm381Controls struct {
	Input1                   *controls.SwitchControl
	Input1Counter            *controls.ValueControl
	Input1SinglePressCounter *controls.ValueControl
	Input1LongPressCounter   *controls.ValueControl
	Input2                   *controls.SwitchControl
	Input2Counter            *controls.ValueControl
	Input2SinglePressCounter *controls.ValueControl
	Input2LongPressCounter   *controls.ValueControl
	Input3                   *controls.SwitchControl
	Input3Counter            *controls.ValueControl
	Input3SinglePressCounter *controls.ValueControl
	Input3LongPressCounter   *controls.ValueControl
	Input4                   *controls.SwitchControl
	Input4Counter            *controls.ValueControl
	Input4SinglePressCounter *controls.ValueControl
	Input4LongPressCounter   *controls.ValueControl
	Input5                   *controls.SwitchControl
	Input5Counter            *controls.ValueControl
	Input5SinglePressCounter *controls.ValueControl
	Input5LongPressCounter   *controls.ValueControl
	Input6                   *controls.SwitchControl
	Input6Counter            *controls.ValueControl
	Input6SinglePressCounter *controls.ValueControl
	Input6LongPressCounter   *controls.ValueControl
	K1                       *controls.SwitchControl
	Channel1                 *controls.RangeControl
	K2                       *controls.SwitchControl
	Channel2                 *controls.RangeControl
	K3                       *controls.SwitchControl
	Channel3                 *controls.RangeControl
	Serial                   *controls.TextControl
	AcOnLN                   *controls.SwitchControl
	Overcurrent              *controls.SwitchControl
}

type WbMdm381 struct {
	name     string
	device   string
	address  string
	Controls *WbMdm381Controls
}

func (w *WbMdm381) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbMdm381) GetControlsInfo() []controls.ControlInfo {
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
	onceWbMdm381     sync.Once
	instanceWbMdm381 *WbMdm381
)

func NewWbMdm381(client *mqtt.Client) *WbMdm381 {
	onceWbMdm381.Do(func() {
		device := "wb-mdm3"
		address := "81"
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &WbMdm381Controls{
			Input1: controls.NewSwitchControl(client, name, "Input 1", controls.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 1`},
			}),
			Input1Counter: controls.NewValueControl(client, name, "Input 1 counter", controls.Meta{
				Type: "value",

				Order:    2,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 1`},
			}),
			Input1SinglePressCounter: controls.NewValueControl(client, name, "Input 1 Single Press Counter", controls.Meta{
				Type: "value",

				Order:    3,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик коротких нажатий входа 1`},
			}),
			Input1LongPressCounter: controls.NewValueControl(client, name, "Input 1 Long Press Counter", controls.Meta{
				Type: "value",

				Order:    4,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик длинных нажатий входа 1`},
			}),
			Input2: controls.NewSwitchControl(client, name, "Input 2", controls.Meta{
				Type: "switch",

				Order:    5,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 2`},
			}),
			Input2Counter: controls.NewValueControl(client, name, "Input 2 counter", controls.Meta{
				Type: "value",

				Order:    6,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 2`},
			}),
			Input2SinglePressCounter: controls.NewValueControl(client, name, "Input 2 Single Press Counter", controls.Meta{
				Type: "value",

				Order:    7,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик коротких нажатий входа 2`},
			}),
			Input2LongPressCounter: controls.NewValueControl(client, name, "Input 2 Long Press Counter", controls.Meta{
				Type: "value",

				Order:    8,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик длинных нажатий входа 2`},
			}),
			Input3: controls.NewSwitchControl(client, name, "Input 3", controls.Meta{
				Type: "switch",

				Order:    9,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 3`},
			}),
			Input3Counter: controls.NewValueControl(client, name, "Input 3 counter", controls.Meta{
				Type: "value",

				Order:    10,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 3`},
			}),
			Input3SinglePressCounter: controls.NewValueControl(client, name, "Input 3 Single Press Counter", controls.Meta{
				Type: "value",

				Order:    11,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик коротких нажатий входа 3`},
			}),
			Input3LongPressCounter: controls.NewValueControl(client, name, "Input 3 Long Press Counter", controls.Meta{
				Type: "value",

				Order:    12,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик длинных нажатий входа 3`},
			}),
			Input4: controls.NewSwitchControl(client, name, "Input 4", controls.Meta{
				Type: "switch",

				Order:    13,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 4`},
			}),
			Input4Counter: controls.NewValueControl(client, name, "Input 4 counter", controls.Meta{
				Type: "value",

				Order:    14,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 4`},
			}),
			Input4SinglePressCounter: controls.NewValueControl(client, name, "Input 4 Single Press Counter", controls.Meta{
				Type: "value",

				Order:    15,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик коротких нажатий входа 4`},
			}),
			Input4LongPressCounter: controls.NewValueControl(client, name, "Input 4 Long Press Counter", controls.Meta{
				Type: "value",

				Order:    16,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик длинных нажатий входа 4`},
			}),
			Input5: controls.NewSwitchControl(client, name, "Input 5", controls.Meta{
				Type: "switch",

				Order:    17,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 5`},
			}),
			Input5Counter: controls.NewValueControl(client, name, "Input 5 counter", controls.Meta{
				Type: "value",

				Order:    18,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 5`},
			}),
			Input5SinglePressCounter: controls.NewValueControl(client, name, "Input 5 Single Press Counter", controls.Meta{
				Type: "value",

				Order:    19,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик коротких нажатий входа 5`},
			}),
			Input5LongPressCounter: controls.NewValueControl(client, name, "Input 5 Long Press Counter", controls.Meta{
				Type: "value",

				Order:    20,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик длинных нажатий входа 5`},
			}),
			Input6: controls.NewSwitchControl(client, name, "Input 6", controls.Meta{
				Type: "switch",

				Order:    21,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Вход 6`},
			}),
			Input6Counter: controls.NewValueControl(client, name, "Input 6 counter", controls.Meta{
				Type: "value",

				Order:    22,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик 6`},
			}),
			Input6SinglePressCounter: controls.NewValueControl(client, name, "Input 6 Single Press Counter", controls.Meta{
				Type: "value",

				Order:    23,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик коротких нажатий входа 6`},
			}),
			Input6LongPressCounter: controls.NewValueControl(client, name, "Input 6 Long Press Counter", controls.Meta{
				Type: "value",

				Order:    24,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Счетчик длинных нажатий входа 6`},
			}),
			K1: controls.NewSwitchControl(client, name, "K1", controls.Meta{
				Type: "switch",

				Order:    25,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			Channel1: controls.NewRangeControl(client, name, "Channel 1", controls.Meta{
				Type: "range",

				Max: 100,

				Order:    26,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Канал 1`},
			}),
			K2: controls.NewSwitchControl(client, name, "K2", controls.Meta{
				Type: "switch",

				Order:    27,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			Channel2: controls.NewRangeControl(client, name, "Channel 2", controls.Meta{
				Type: "range",

				Max: 100,

				Order:    28,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Канал 2`},
			}),
			K3: controls.NewSwitchControl(client, name, "K3", controls.Meta{
				Type: "switch",

				Order:    29,
				ReadOnly: false,
				Title:    controls.MultilingualText{},
			}),
			Channel3: controls.NewRangeControl(client, name, "Channel 3", controls.Meta{
				Type: "range",

				Max: 100,

				Order:    30,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Канал 3`},
			}),
			Serial: controls.NewTextControl(client, name, "Serial", controls.Meta{
				Type: "text",

				Order:    31,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Серийный номер`},
			}),
			AcOnLN: controls.NewSwitchControl(client, name, "AC on L-N", controls.Meta{
				Type: "switch",

				Order:    32,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Напряжение на клеммах L-N`},
			}),
			Overcurrent: controls.NewSwitchControl(client, name, "Overcurrent", controls.Meta{
				Type: "switch",

				Order:    33,
				ReadOnly: true,
				Title:    controls.MultilingualText{"ru": `Перегрузка по току`},
			}),
		}

		instanceWbMdm381 = &WbMdm381{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceWbMdm381
}
