package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type WbMdm381Controls struct {
	Input1                   *control.SwitchControl
	Input1Counter            *control.ValueControl
	Input1SinglePressCounter *control.ValueControl
	Input1LongPressCounter   *control.ValueControl
	Input2                   *control.SwitchControl
	Input2Counter            *control.ValueControl
	Input2SinglePressCounter *control.ValueControl
	Input2LongPressCounter   *control.ValueControl
	Input3                   *control.SwitchControl
	Input3Counter            *control.ValueControl
	Input3SinglePressCounter *control.ValueControl
	Input3LongPressCounter   *control.ValueControl
	Input4                   *control.SwitchControl
	Input4Counter            *control.ValueControl
	Input4SinglePressCounter *control.ValueControl
	Input4LongPressCounter   *control.ValueControl
	Input5                   *control.SwitchControl
	Input5Counter            *control.ValueControl
	Input5SinglePressCounter *control.ValueControl
	Input5LongPressCounter   *control.ValueControl
	Input6                   *control.SwitchControl
	Input6Counter            *control.ValueControl
	Input6SinglePressCounter *control.ValueControl
	Input6LongPressCounter   *control.ValueControl
	K1                       *control.SwitchControl
	Channel1                 *control.RangeControl
	K2                       *control.SwitchControl
	Channel2                 *control.RangeControl
	K3                       *control.SwitchControl
	Channel3                 *control.RangeControl
	Serial                   *control.TextControl
	AcOnLN                   *control.SwitchControl
	Overcurrent              *control.SwitchControl
}

type WbMdm381 struct {
	name     string
	Controls *WbMdm381Controls
}

func (w *WbMdm381) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		ControlsInfo: controlsInfo,
	}
}

func (w *WbMdm381) GetControlsInfo() []control.ControlInfo {
	var infoList []control.ControlInfo

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
				info := method.Call(nil)[0].Interface().(control.ControlInfo)
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
		name := "wb-mdm3_81"

		controlList := &WbMdm381Controls{
			Input1: control.NewSwitchControl(client, name, "Input 1", control.Meta{
				Type: "switch",

				Order:    1,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 1`},
			}),
			Input1Counter: control.NewValueControl(client, name, "Input 1 counter", control.Meta{
				Type: "value",

				Order:    2,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 1`},
			}),
			Input1SinglePressCounter: control.NewValueControl(client, name, "Input 1 Single Press Counter", control.Meta{
				Type: "value",

				Order:    3,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик коротких нажатий входа 1`},
			}),
			Input1LongPressCounter: control.NewValueControl(client, name, "Input 1 Long Press Counter", control.Meta{
				Type: "value",

				Order:    4,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик длинных нажатий входа 1`},
			}),
			Input2: control.NewSwitchControl(client, name, "Input 2", control.Meta{
				Type: "switch",

				Order:    5,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 2`},
			}),
			Input2Counter: control.NewValueControl(client, name, "Input 2 counter", control.Meta{
				Type: "value",

				Order:    6,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 2`},
			}),
			Input2SinglePressCounter: control.NewValueControl(client, name, "Input 2 Single Press Counter", control.Meta{
				Type: "value",

				Order:    7,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик коротких нажатий входа 2`},
			}),
			Input2LongPressCounter: control.NewValueControl(client, name, "Input 2 Long Press Counter", control.Meta{
				Type: "value",

				Order:    8,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик длинных нажатий входа 2`},
			}),
			Input3: control.NewSwitchControl(client, name, "Input 3", control.Meta{
				Type: "switch",

				Order:    9,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 3`},
			}),
			Input3Counter: control.NewValueControl(client, name, "Input 3 counter", control.Meta{
				Type: "value",

				Order:    10,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 3`},
			}),
			Input3SinglePressCounter: control.NewValueControl(client, name, "Input 3 Single Press Counter", control.Meta{
				Type: "value",

				Order:    11,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик коротких нажатий входа 3`},
			}),
			Input3LongPressCounter: control.NewValueControl(client, name, "Input 3 Long Press Counter", control.Meta{
				Type: "value",

				Order:    12,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик длинных нажатий входа 3`},
			}),
			Input4: control.NewSwitchControl(client, name, "Input 4", control.Meta{
				Type: "switch",

				Order:    13,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 4`},
			}),
			Input4Counter: control.NewValueControl(client, name, "Input 4 counter", control.Meta{
				Type: "value",

				Order:    14,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 4`},
			}),
			Input4SinglePressCounter: control.NewValueControl(client, name, "Input 4 Single Press Counter", control.Meta{
				Type: "value",

				Order:    15,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик коротких нажатий входа 4`},
			}),
			Input4LongPressCounter: control.NewValueControl(client, name, "Input 4 Long Press Counter", control.Meta{
				Type: "value",

				Order:    16,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик длинных нажатий входа 4`},
			}),
			Input5: control.NewSwitchControl(client, name, "Input 5", control.Meta{
				Type: "switch",

				Order:    17,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 5`},
			}),
			Input5Counter: control.NewValueControl(client, name, "Input 5 counter", control.Meta{
				Type: "value",

				Order:    18,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 5`},
			}),
			Input5SinglePressCounter: control.NewValueControl(client, name, "Input 5 Single Press Counter", control.Meta{
				Type: "value",

				Order:    19,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик коротких нажатий входа 5`},
			}),
			Input5LongPressCounter: control.NewValueControl(client, name, "Input 5 Long Press Counter", control.Meta{
				Type: "value",

				Order:    20,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик длинных нажатий входа 5`},
			}),
			Input6: control.NewSwitchControl(client, name, "Input 6", control.Meta{
				Type: "switch",

				Order:    21,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Вход 6`},
			}),
			Input6Counter: control.NewValueControl(client, name, "Input 6 counter", control.Meta{
				Type: "value",

				Order:    22,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик 6`},
			}),
			Input6SinglePressCounter: control.NewValueControl(client, name, "Input 6 Single Press Counter", control.Meta{
				Type: "value",

				Order:    23,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик коротких нажатий входа 6`},
			}),
			Input6LongPressCounter: control.NewValueControl(client, name, "Input 6 Long Press Counter", control.Meta{
				Type: "value",

				Order:    24,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Счетчик длинных нажатий входа 6`},
			}),
			K1: control.NewSwitchControl(client, name, "K1", control.Meta{
				Type: "switch",

				Order:    25,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			Channel1: control.NewRangeControl(client, name, "Channel 1", control.Meta{
				Type: "range",

				Max: 100,

				Order:    26,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Канал 1`},
			}),
			K2: control.NewSwitchControl(client, name, "K2", control.Meta{
				Type: "switch",

				Order:    27,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			Channel2: control.NewRangeControl(client, name, "Channel 2", control.Meta{
				Type: "range",

				Max: 100,

				Order:    28,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Канал 2`},
			}),
			K3: control.NewSwitchControl(client, name, "K3", control.Meta{
				Type: "switch",

				Order:    29,
				ReadOnly: false,
				Title:    control.MultilingualText{},
			}),
			Channel3: control.NewRangeControl(client, name, "Channel 3", control.Meta{
				Type: "range",

				Max: 100,

				Order:    30,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Канал 3`},
			}),
			Serial: control.NewTextControl(client, name, "Serial", control.Meta{
				Type: "text",

				Order:    31,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Серийный номер`},
			}),
			AcOnLN: control.NewSwitchControl(client, name, "AC on L-N", control.Meta{
				Type: "switch",

				Order:    32,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Напряжение на клеммах L-N`},
			}),
			Overcurrent: control.NewSwitchControl(client, name, "Overcurrent", control.Meta{
				Type: "switch",

				Order:    33,
				ReadOnly: true,
				Title:    control.MultilingualText{"ru": `Перегрузка по току`},
			}),
		}

		instanceWbMdm381 = &WbMdm381{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceWbMdm381
}
