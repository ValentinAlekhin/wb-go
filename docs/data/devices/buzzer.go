package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type BuzzerControls struct {
	Enabled   *controls.SwitchControl
	Frequency *controls.RangeControl
	Volume    *controls.RangeControl
}

type Buzzer struct {
	Name     string
	Address  string
	Controls *BuzzerControls
}

func (w *Buzzer) GetControlsInfo() []controls.ControlInfo {
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
	onceBuzzer     sync.Once
	instanceBuzzer *Buzzer
)

func NewBuzzer(client *mqtt.Client) *Buzzer {
	onceBuzzer.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "buzzer", "")
		controlList := &BuzzerControls{
			Enabled:   controls.NewSwitchControl(client, deviceName, "enabled"),
			Frequency: controls.NewRangeControl(client, deviceName, "frequency"),
			Volume:    controls.NewRangeControl(client, deviceName, "volume"),
		}

		instanceBuzzer = &Buzzer{
			Name:     deviceName,
			Address:  "",
			Controls: controlList,
		}
	})

	return instanceBuzzer
}
