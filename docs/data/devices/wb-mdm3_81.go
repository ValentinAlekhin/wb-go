package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
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
	Name     string
	Controls *WbMdm381Controls
}

var (
	onceWbMdm381     sync.Once
	instanceWbMdm381 *WbMdm381
)

func NewWbMdm381(client *mqtt.Client) *WbMdm381 {
	onceWbMdm381.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "wb-mdm3", "81")
		controlList := &WbMdm381Controls{
			Input1:                   controls.NewSwitchControl(client, deviceName, "Input 1"),
			Input1Counter:            controls.NewValueControl(client, deviceName, "Input 1 counter"),
			Input1SinglePressCounter: controls.NewValueControl(client, deviceName, "Input 1 Single Press Counter"),
			Input1LongPressCounter:   controls.NewValueControl(client, deviceName, "Input 1 Long Press Counter"),
			Input2:                   controls.NewSwitchControl(client, deviceName, "Input 2"),
			Input2Counter:            controls.NewValueControl(client, deviceName, "Input 2 counter"),
			Input2SinglePressCounter: controls.NewValueControl(client, deviceName, "Input 2 Single Press Counter"),
			Input2LongPressCounter:   controls.NewValueControl(client, deviceName, "Input 2 Long Press Counter"),
			Input3:                   controls.NewSwitchControl(client, deviceName, "Input 3"),
			Input3Counter:            controls.NewValueControl(client, deviceName, "Input 3 counter"),
			Input3SinglePressCounter: controls.NewValueControl(client, deviceName, "Input 3 Single Press Counter"),
			Input3LongPressCounter:   controls.NewValueControl(client, deviceName, "Input 3 Long Press Counter"),
			Input4:                   controls.NewSwitchControl(client, deviceName, "Input 4"),
			Input4Counter:            controls.NewValueControl(client, deviceName, "Input 4 counter"),
			Input4SinglePressCounter: controls.NewValueControl(client, deviceName, "Input 4 Single Press Counter"),
			Input4LongPressCounter:   controls.NewValueControl(client, deviceName, "Input 4 Long Press Counter"),
			Input5:                   controls.NewSwitchControl(client, deviceName, "Input 5"),
			Input5Counter:            controls.NewValueControl(client, deviceName, "Input 5 counter"),
			Input5SinglePressCounter: controls.NewValueControl(client, deviceName, "Input 5 Single Press Counter"),
			Input5LongPressCounter:   controls.NewValueControl(client, deviceName, "Input 5 Long Press Counter"),
			Input6:                   controls.NewSwitchControl(client, deviceName, "Input 6"),
			Input6Counter:            controls.NewValueControl(client, deviceName, "Input 6 counter"),
			Input6SinglePressCounter: controls.NewValueControl(client, deviceName, "Input 6 Single Press Counter"),
			Input6LongPressCounter:   controls.NewValueControl(client, deviceName, "Input 6 Long Press Counter"),
			K1:                       controls.NewSwitchControl(client, deviceName, "K1"),
			Channel1:                 controls.NewRangeControl(client, deviceName, "Channel 1"),
			K2:                       controls.NewSwitchControl(client, deviceName, "K2"),
			Channel2:                 controls.NewRangeControl(client, deviceName, "Channel 2"),
			K3:                       controls.NewSwitchControl(client, deviceName, "K3"),
			Channel3:                 controls.NewRangeControl(client, deviceName, "Channel 3"),
			Serial:                   controls.NewTextControl(client, deviceName, "Serial"),
			AcOnLN:                   controls.NewSwitchControl(client, deviceName, "AC on L-N"),
			Overcurrent:              controls.NewSwitchControl(client, deviceName, "Overcurrent"),
		}

		instanceWbMdm381 = &WbMdm381{
			Name:     deviceName,
			Controls: controlList,
		}
	})

	return instanceWbMdm381
}
