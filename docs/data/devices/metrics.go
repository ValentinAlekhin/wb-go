package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"sync"
)

type MetricsControls struct {
	LoadAverage1Min   *ValueControl
	LoadAverage5Min   *ValueControl
	LoadAverage15Min  *ValueControl
	RamAvailable      *ValueControl
	RamUsed           *ValueControl
	RamTotal          *ValueControl
	SwapTotal         *ValueControl
	SwapUsed          *ValueControl
	DevRootUsedSpace  *ValueControl
	DevRootTotalSpace *ValueControl
	DevRootLinkedOn   *TextControl
	DataUsedSpace     *ValueControl
	DataTotalSpace    *ValueControl
}

type Metrics struct {
	Name     string
	Controls *MetricsControls
}

var (
	onceMetrics     sync.Once
	instanceMetrics *Metrics
)

func NewMetrics(client *mqtt.Client) *Metrics {
	onceMetrics.Do(func() {
		name := "metrics"
		deviceTopic := fmt.Sprintf("/devices/%s_%s", name, "")
		controls := &MetricsControls{
			LoadAverage1Min:   NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "load_average_1min")),
			LoadAverage5Min:   NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "load_average_5min")),
			LoadAverage15Min:  NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "load_average_15min")),
			RamAvailable:      NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "ram_available")),
			RamUsed:           NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "ram_used")),
			RamTotal:          NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "ram_total")),
			SwapTotal:         NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "swap_total")),
			SwapUsed:          NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "swap_used")),
			DevRootUsedSpace:  NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "dev_root_used_space")),
			DevRootTotalSpace: NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "dev_root_total_space")),
			DevRootLinkedOn:   NewTextControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "dev_root_linked_on")),
			DataUsedSpace:     NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "data_used_space")),
			DataTotalSpace:    NewValueControl(client, fmt.Sprintf("%s/controls/%s", deviceTopic, "data_total_space")),
		}

		instanceMetrics = &Metrics{
			Name:     name,
			Controls: controls,
		}
	})

	return instanceMetrics
}
