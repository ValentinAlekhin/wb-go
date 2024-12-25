package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type Dooya0X0104Controls struct {
	Position       *controls.RangeControl
	Open           *controls.PushbuttonControl
	Close          *controls.PushbuttonControl
	Stop           *controls.PushbuttonControl
	FactoryDefault *controls.PushbuttonControl
}

type Dooya0X0104 struct {
	Name     string
	Address  string
	Controls *Dooya0X0104Controls
}

func (w *Dooya0X0104) GetControlsInfo() []controls.ControlInfo {
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
	onceDooya0X0104     sync.Once
	instanceDooya0X0104 *Dooya0X0104
)

func NewDooya0X0104(client *mqtt.Client) *Dooya0X0104 {
	onceDooya0X0104.Do(func() {
		deviceName := fmt.Sprintf("%s_%s", "dooya", "0x0104")
		controlList := &Dooya0X0104Controls{
			Position:       controls.NewRangeControl(client, deviceName, "Position"),
			Open:           controls.NewPushbuttonControl(client, deviceName, "Open"),
			Close:          controls.NewPushbuttonControl(client, deviceName, "Close"),
			Stop:           controls.NewPushbuttonControl(client, deviceName, "Stop"),
			FactoryDefault: controls.NewPushbuttonControl(client, deviceName, "Factory Default"),
		}

		instanceDooya0X0104 = &Dooya0X0104{
			Name:     deviceName,
			Address:  "0x0104",
			Controls: controlList,
		}
	})

	return instanceDooya0X0104
}
