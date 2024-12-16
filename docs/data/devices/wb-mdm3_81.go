package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type WbMdm381Controls struct {
	Input1                   *SwitchControl
	Input1Counter            *ValueControl
	Input1SinglePressCounter *ValueControl
	Input1LongPressCounter   *ValueControl
	Input2                   *SwitchControl
	Input2Counter            *ValueControl
	Input2SinglePressCounter *ValueControl
	Input2LongPressCounter   *ValueControl
	Input3                   *SwitchControl
	Input3Counter            *ValueControl
	Input3SinglePressCounter *ValueControl
	Input3LongPressCounter   *ValueControl
	Input4                   *SwitchControl
	Input4Counter            *ValueControl
	Input4SinglePressCounter *ValueControl
	Input4LongPressCounter   *ValueControl
	Input5                   *SwitchControl
	Input5Counter            *ValueControl
	Input5SinglePressCounter *ValueControl
	Input5LongPressCounter   *ValueControl
	Input6                   *SwitchControl
	Input6Counter            *ValueControl
	Input6SinglePressCounter *ValueControl
	Input6LongPressCounter   *ValueControl
	K1                       *SwitchControl
	Channel1                 *RangeControl
	K2                       *SwitchControl
	Channel2                 *RangeControl
	K3                       *SwitchControl
	Channel3                 *RangeControl
	Serial                   *TextControl
	AcOnLN                   *SwitchControl
	Overcurrent              *SwitchControl
}

type WbMdm381 struct {
	Name          string
	ModbusAddress int32
	Controls      *WbMdm381Controls
}

var (
	onceWbMdm381     sync.Once
	instanceWbMdm381 *WbMdm381
)

func NewWbMdm381(client *mqtt.Client) *WbMdm381 {
	onceWbMdm381.Do(func() {
		name := "wb-mdm3"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "81")
		controls := &WbMdm381Controls{
			Input1:                   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 1")),
			Input1Counter:            NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 1 counter")),
			Input1SinglePressCounter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 1 Single Press Counter")),
			Input1LongPressCounter:   NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 1 Long Press Counter")),
			Input2:                   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 2")),
			Input2Counter:            NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 2 counter")),
			Input2SinglePressCounter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 2 Single Press Counter")),
			Input2LongPressCounter:   NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 2 Long Press Counter")),
			Input3:                   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 3")),
			Input3Counter:            NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 3 counter")),
			Input3SinglePressCounter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 3 Single Press Counter")),
			Input3LongPressCounter:   NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 3 Long Press Counter")),
			Input4:                   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 4")),
			Input4Counter:            NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 4 counter")),
			Input4SinglePressCounter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 4 Single Press Counter")),
			Input4LongPressCounter:   NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 4 Long Press Counter")),
			Input5:                   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 5")),
			Input5Counter:            NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 5 counter")),
			Input5SinglePressCounter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 5 Single Press Counter")),
			Input5LongPressCounter:   NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 5 Long Press Counter")),
			Input6:                   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 6")),
			Input6Counter:            NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 6 counter")),
			Input6SinglePressCounter: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 6 Single Press Counter")),
			Input6LongPressCounter:   NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Input 6 Long Press Counter")),
			K1:                       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "K1")),
			Channel1:                 NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Channel 1")),
			K2:                       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "K2")),
			Channel2:                 NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Channel 2")),
			K3:                       NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "K3")),
			Channel3:                 NewRangeControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Channel 3")),
			Serial:                   NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Serial")),
			AcOnLN:                   NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "AC on L-N")),
			Overcurrent:              NewSwitchControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "Overcurrent")),
		}

		instanceWbMdm381 = &WbMdm381{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceWbMdm381
}
