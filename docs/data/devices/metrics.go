package devices

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/controls"
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
	Name     string
	Address  string
	Controls *MetricsControls
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
		deviceName := fmt.Sprintf("%s_%s", "metrics", "")
		controlList := &MetricsControls{
			LoadAverage1Min:   controls.NewValueControl(client, deviceName, "load_average_1min"),
			LoadAverage5Min:   controls.NewValueControl(client, deviceName, "load_average_5min"),
			LoadAverage15Min:  controls.NewValueControl(client, deviceName, "load_average_15min"),
			RamAvailable:      controls.NewValueControl(client, deviceName, "ram_available"),
			RamUsed:           controls.NewValueControl(client, deviceName, "ram_used"),
			RamTotal:          controls.NewValueControl(client, deviceName, "ram_total"),
			SwapTotal:         controls.NewValueControl(client, deviceName, "swap_total"),
			SwapUsed:          controls.NewValueControl(client, deviceName, "swap_used"),
			DevRootUsedSpace:  controls.NewValueControl(client, deviceName, "dev_root_used_space"),
			DevRootTotalSpace: controls.NewValueControl(client, deviceName, "dev_root_total_space"),
			DevRootLinkedOn:   controls.NewTextControl(client, deviceName, "dev_root_linked_on"),
			DataUsedSpace:     controls.NewValueControl(client, deviceName, "data_used_space"),
			DataTotalSpace:    controls.NewValueControl(client, deviceName, "data_total_space"),
		}

		instanceMetrics = &Metrics{
			Name:     deviceName,
			Address:  "",
			Controls: controlList,
		}
	})

	return instanceMetrics
}
