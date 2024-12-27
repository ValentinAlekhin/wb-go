package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
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
	name     string
	device   string
	address  string
	Controls *Dooya0X0104Controls
}

func (w *Dooya0X0104) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
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
		device := "dooya"
		address := "0x0104"
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &Dooya0X0104Controls{
			Position: controls.NewRangeControl(client, name, "Position", controls.Meta{
				Type: "range",

				Max: 100,

				Order:    1,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Позиция`},
			}),
			Open: controls.NewPushbuttonControl(client, name, "Open", controls.Meta{
				Type: "pushbutton",

				Order:    2,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Открыть`},
			}),
			Close: controls.NewPushbuttonControl(client, name, "Close", controls.Meta{
				Type: "pushbutton",

				Order:    3,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Закрыть`},
			}),
			Stop: controls.NewPushbuttonControl(client, name, "Stop", controls.Meta{
				Type: "pushbutton",

				Order:    4,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Остановить`},
			}),
			FactoryDefault: controls.NewPushbuttonControl(client, name, "Factory Default", controls.Meta{
				Type: "pushbutton",

				Order:    5,
				ReadOnly: false,
				Title:    controls.MultilingualText{"ru": `Сбросить настройки`},
			}),
		}

		instanceDooya0X0104 = &Dooya0X0104{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceDooya0X0104
}
