package device

import (
	"github.com/ValentinAlekhin/wb-go/pkg/basedevice"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type MetricsControls struct {
	LoadAverage1Min *control.ValueControl
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
		}

		instanceMetrics = &Metrics{
			name:     name,
			Controls: controlList,
		}
	})

	return instanceMetrics
}
