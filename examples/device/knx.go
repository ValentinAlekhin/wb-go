package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type Knxcontrols struct {
	Data *control.TextControl
}

type Knx struct {
	name     string
	device   string
	address  string
	Controls *Knxcontrols
}

func (w *Knx) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Knx) GetControlsInfo() []control.ControlInfo {
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
	onceKnx     sync.Once
	instanceKnx *Knx
)

func NewKnx(client *mqtt.Client) *Knx {
	onceKnx.Do(func() {
		device := "knx"
		address := ""
		name := device

		controlList := &Knxcontrols{
			Data: control.NewTextControl(client, name, "data", control.Meta{
				Type: "text",

				Order:    0,
				ReadOnly: false,
				Title:    control.MultilingualText{"en": `Message`, "ru": `Сообщение`},
			}),
		}

		instanceKnx = &Knx{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceKnx
}
