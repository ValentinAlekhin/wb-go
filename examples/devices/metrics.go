package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
	"github.com/ValentinAlekhin/wb-go/pkg/deviceInfo"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"reflect"
	"sync"
)

type MetricsControls struct {
	LoadAverage1Min   *controls.ValueControl
	LoadAverage5Min   *controls.ValueControl
	LoadAverage15Min  *controls.ValueControl
	RamAvailable      *controls.ValueControl
	RamUsed           *controls.ValueControl
	RamTotal          *controls.ValueControl
	SwapTotal         *controls.ValueControl
	SwapUsed          *controls.ValueControl
	DevRootUsedSpace  *controls.ValueControl
	DevRootTotalSpace *controls.ValueControl
	DevRootLinkedOn   *controls.TextControl
	DataUsedSpace     *controls.ValueControl
	DataTotalSpace    *controls.ValueControl
}

type Metrics struct {
	name     string
	device   string
	address  string
	Controls *MetricsControls
}

func (w *Metrics) GetInfo() deviceInfo.DeviceInfo {
	controlsInfo := w.GetControlsInfo()

	return deviceInfo.DeviceInfo{
		Name:         w.name,
		Device:       w.device,
		Address:      w.address,
		ControlsInfo: controlsInfo,
	}
}

func (w *Metrics) GetControlsInfo() []controls.ControlInfo {
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
	onceMetrics     sync.Once
	instanceMetrics *Metrics
)

func NewMetrics(client *mqtt.Client) *Metrics {
	onceMetrics.Do(func() {
		device := "metrics"
		address := ""
		name := fmt.Sprintf("%s_%s", device, address)
		controlList := &MetricsControls{
			LoadAverage1Min: controls.NewValueControl(client, name, "load_average_1min", controls.Meta{
				Type:  "value",
				Units: "tasks",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			LoadAverage5Min: controls.NewValueControl(client, name, "load_average_5min", controls.Meta{
				Type:  "value",
				Units: "tasks",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			LoadAverage15Min: controls.NewValueControl(client, name, "load_average_15min", controls.Meta{
				Type:  "value",
				Units: "tasks",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			RamAvailable: controls.NewValueControl(client, name, "ram_available", controls.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			RamUsed: controls.NewValueControl(client, name, "ram_used", controls.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			RamTotal: controls.NewValueControl(client, name, "ram_total", controls.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			SwapTotal: controls.NewValueControl(client, name, "swap_total", controls.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			SwapUsed: controls.NewValueControl(client, name, "swap_used", controls.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			DevRootUsedSpace: controls.NewValueControl(client, name, "dev_root_used_space", controls.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			DevRootTotalSpace: controls.NewValueControl(client, name, "dev_root_total_space", controls.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			DevRootLinkedOn: controls.NewTextControl(client, name, "dev_root_linked_on", controls.Meta{
				Type: "text",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			DataUsedSpace: controls.NewValueControl(client, name, "data_used_space", controls.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
			}),
			DataTotalSpace: controls.NewValueControl(client, name, "data_total_space", controls.Meta{
				Type:  "value",
				Units: "MiB",

				Order:    0,
				ReadOnly: true,
				Title:    controls.MultilingualText{},
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
