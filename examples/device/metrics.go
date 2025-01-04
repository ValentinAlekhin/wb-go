package device

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceinfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type MetricsControls struct {
	LoadAverage1Min   *control.ValueControl
	LoadAverage5Min   *control.ValueControl
	LoadAverage15Min  *control.ValueControl
	RamAvailable      *control.ValueControl
	RamUsed           *control.ValueControl
	RamTotal          *control.ValueControl
	SwapTotal         *control.ValueControl
	SwapUsed          *control.ValueControl
	DevRootUsedSpace  *control.ValueControl
	DevRootTotalSpace *control.ValueControl
	DevRootLinkedOn   *control.TextControl
	DataUsedSpace     *control.ValueControl
	DataTotalSpace    *control.ValueControl
}

type Metrics struct {
	name     string
	device   string
	address  string
	Controls *MetricsControls
}

func (w *Metrics) GetInfo() deviceinfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceinfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Metrics) GetControlsInfo() []control.ControlInfo {
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
	onceMetrics     sync.Once
	instanceMetrics *Metrics
)

func NewMetrics(client *mqtt.Client) *Metrics {
	onceMetrics.Do(func() {
		device := "metrics"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &MetricsControls{
			LoadAverage1Min: control.NewValueControl(client, name, "load_average_1min", control.Meta{
				Type:  "value",
				Units: "tasks",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			LoadAverage5Min: control.NewValueControl(client, name, "load_average_5min", control.Meta{
				Type:  "value",
				Units: "tasks",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			LoadAverage15Min: control.NewValueControl(client, name, "load_average_15min", control.Meta{
				Type:  "value",
				Units: "tasks",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			RamAvailable: control.NewValueControl(client, name, "ram_available", control.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			RamUsed: control.NewValueControl(client, name, "ram_used", control.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			RamTotal: control.NewValueControl(client, name, "ram_total", control.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			SwapTotal: control.NewValueControl(client, name, "swap_total", control.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			SwapUsed: control.NewValueControl(client, name, "swap_used", control.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			DevRootUsedSpace: control.NewValueControl(client, name, "dev_root_used_space", control.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			DevRootTotalSpace: control.NewValueControl(client, name, "dev_root_total_space", control.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			DevRootLinkedOn: control.NewTextControl(client, name, "dev_root_linked_on", control.Meta{
				Type: "text",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			DataUsedSpace: control.NewValueControl(client, name, "data_used_space", control.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
			DataTotalSpace: control.NewValueControl(client, name, "data_total_space", control.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    control.MultilingualText{},
			}),
		}

		instanceMetrics = &Metrics{
			name:     name,
			device:   device,
			address:  address,
			Controls: controlList,
		}
	})

	return instanceMetrics
}
