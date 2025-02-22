package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type Dooya0X0104Controls struct {
	Position       *control.RangeControl
	Open           *control.PushbuttonControl
	Close          *control.PushbuttonControl
	Stop           *control.PushbuttonControl
	FactoryDefault *control.PushbuttonControl
}

type Dooya0X0104 struct {
	name     string
	Controls *Dooya0X0104Controls
}

func (w *Dooya0X0104) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceDooya0X0104     sync.Once
	instanceDooya0X0104 *Dooya0X0104
)

func NewDooya0X0104(client mqtt.ClientInterface) *Dooya0X0104 {
	onceDooya0X0104.Do(func() {
		name := "dooya_0x0104"

		controlList := &Dooya0X0104Controls{
			Position: control.NewRangeControl(client, name, "Position", control.Meta{
				Type: "range",

				Max: 100,

				Order:    1,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Позиция`},
			}),
			Open: control.NewPushbuttonControl(client, name, "Open", control.Meta{
				Type: "pushbutton",

				Order:    2,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Открыть`},
			}),
			Close: control.NewPushbuttonControl(client, name, "Close", control.Meta{
				Type: "pushbutton",

				Order:    3,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Закрыть`},
			}),
			Stop: control.NewPushbuttonControl(client, name, "Stop", control.Meta{
				Type: "pushbutton",

				Order:    4,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Остановить`},
			}),
			FactoryDefault: control.NewPushbuttonControl(client, name, "Factory Default", control.Meta{
				Type: "pushbutton",

				Order:    5,
				ReadOnly: false,
				Title:    control.MultilingualText{"ru": `Сбросить настройки`},
			}),
		}

		instanceDooya0X0104 = &Dooya0X0104{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceDooya0X0104
}
