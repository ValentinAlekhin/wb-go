package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
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
	Controls *MetricsControls
}

func (w *Metrics) GetInfo() basedevice.Info {
	return basedevice.Info{
		Name:         w.name,
		ControlsInfo: basedevice.GetControlsInfo(w.Controls),
	}
}

var (
	onceMetrics     sync.Once
	instanceMetrics *Metrics
)

func NewMetrics(client mqtt.ClientInterface) *Metrics {
	onceMetrics.Do(func() {
		name := "metrics"

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
			Controls: controlList,
		}
	})

	return instanceMetrics
}
