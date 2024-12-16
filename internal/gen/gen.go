package gen

import (
	"bytes"
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iancoleman/strcase"
	"go/format"
	"os"
	"slices"
	"strings"
	"text/template"
	"time"
)

type GenerateService struct {
	client            *wb.Client
	outputDir         string
	devicesToGenerate []string
	packageName       string
}

type deviceControlTemplateData struct {
	Name       string
	Mqtt       string
	ReadOnly   bool
	Type       string
	StructName string
}

type deviceTemplateData struct {
	DeviceName               string
	DeviceStructName         string
	DeviceControlsStructName string
	Filename                 string
	ModbusAddress            string
	Controls                 []deviceControlTemplateData
	PackageName              string
}

type watchResultItem struct {
	DeviceName    string
	ModbusAddress string
	Control       string
	Meta          ControlMeta
}

type ControlMeta struct {
	Type      string           `json:"type"`      // Тип контроля
	Units     string           `json:"units"`     // Единицы измерения (только для type="value")
	Max       float64          `json:"max"`       // Максимальное значение
	Min       float64          `json:"min"`       // Минимальное значение
	Precision float64          `json:"precision"` // Точность
	Order     int              `json:"order"`     // Порядок отображения
	ReadOnly  bool             `json:"readonly"`  // Только для чтения
	Title     MultilingualText `json:"title"`     // Название (разные языки)
	Enum      map[string]Enum  `json:"enum"`      // Заголовки для enum
}

// MultilingualText хранит текстовые значения на разных языках
type MultilingualText map[string]string

// Enum хранит значения для enum (вложенные текстовые описания)
type Enum struct {
	Title MultilingualText `json:"title"` // Название enum на разных языках
}

//go:embed templates/*
var embedFs embed.FS

func NewGenerateService(client *wb.Client, output string, devicesToGenerate []string, packageName string) *GenerateService {
	service := &GenerateService{
		client:            client,
		outputDir:         output,
		devicesToGenerate: devicesToGenerate,
		packageName:       packageName,
	}
	return service
}

func (g *GenerateService) Run() {
	watchResults := g.collectData()
	filteredResults := g.filterUnique(watchResults)
	templatesData := g.generateTemplates(filteredResults)
	g.generateFiles(templatesData)
	g.copyControls()
}

func (g GenerateService) collectData() []watchResultItem {
	var list []watchResultItem

	watcher := g.getTopicWatcher(&list)
	g.client.Subscribe("/devices/+/controls/+/meta", watcher)
	time.Sleep(1 * time.Second)

	return list
}

func (g *GenerateService) getTopicWatcher(list *[]watchResultItem) func(client mqtt.Client, msg mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		meta := string(msg.Payload())

		topicParts := strings.Split(topic, "/")
		if len(topicParts) < 5 {
			return
		}

		modbusAddress := ""
		device := topicParts[2]
		deviceParts := strings.Split(device, "_")
		if len(deviceParts) > 1 {
			device = deviceParts[0]
			modbusAddress = deviceParts[1]
		}

		controlName := topicParts[4]

		if slices.Contains(g.devicesToGenerate, device) {
			controlMeta := ControlMeta{}

			err := json.Unmarshal([]byte(meta), &controlMeta)
			if err != nil {
				panic(err)
			}

			item := watchResultItem{
				DeviceName:    device,
				ModbusAddress: modbusAddress,
				Control:       controlName,
				Meta:          controlMeta,
			}

			*list = append(*list, item)
		}
	}
}

func (g *GenerateService) generateTemplates(list []watchResultItem) map[string]*deviceTemplateData {
	deviceMap := map[string]*deviceTemplateData{}

	for _, item := range list {
		key := fmt.Sprintf("%s_%s", item.DeviceName, item.ModbusAddress)
		if val, ok := deviceMap[key]; !ok {
			deviceStructName := strcase.ToCamel(item.DeviceName + item.ModbusAddress)
			deviceControlsStructName := deviceStructName + "Controls"

			filename := item.DeviceName
			if item.ModbusAddress != "" {
				filename = fmt.Sprintf("%s_%s", item.DeviceName, item.ModbusAddress)
			}

			controlTemplate := g.getControlTemplate(item)
			newVal := &deviceTemplateData{
				DeviceName:               item.DeviceName,
				DeviceStructName:         deviceStructName,
				DeviceControlsStructName: deviceControlsStructName,
				Filename:                 filename,
				ModbusAddress:            item.ModbusAddress,
				PackageName:              g.packageName,
				Controls:                 []deviceControlTemplateData{controlTemplate},
			}

			deviceMap[key] = newVal
		} else {
			controlTemplate := g.getControlTemplate(item)
			val.Controls = append(val.Controls, controlTemplate)
		}
	}

	return deviceMap

}

func (g GenerateService) getControlTemplate(control watchResultItem) deviceControlTemplateData {
	typeName := control.Meta.Type
	for key, val := range ControlValueTypeMap {
		if !slices.Contains(val, control.Meta.Type) {
			continue
		}

		typeName = key
	}

	return deviceControlTemplateData{
		Name:       strcase.ToCamel(control.Control),
		Mqtt:       control.Control,
		ReadOnly:   control.Meta.ReadOnly,
		Type:       typeName,
		StructName: strcase.ToCamel(typeName) + "Control",
	}
}

func (g *GenerateService) filterUnique(list []watchResultItem) []watchResultItem {
	uniqueMap := make(map[string]struct{})
	var result []watchResultItem

	for _, item := range list {
		key := item.DeviceName + "|" + item.ModbusAddress + "|" + item.Control
		if _, exists := uniqueMap[key]; !exists {
			uniqueMap[key] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func (g GenerateService) generateFiles(data map[string]*deviceTemplateData) {
	outputDir := g.getOutputDir()

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, v := range data {
		g.generateFile(v, outputDir)
	}
}

func (g *GenerateService) getOutputDir() string {
	if g.outputDir != "" {
		return g.outputDir
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return dir
}

func (g *GenerateService) generateFile(data *deviceTemplateData, outputDir string) {
	tmpl, err := template.ParseFS(embedFs, "templates/device.txt")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}

	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("Ошибка форматирования:", err)
		fmt.Println("Сгенерированный код до форматирования:")
		fmt.Println(buf.String())
		return
	}

	outputPath := fmt.Sprintf("%s/%s.go", outputDir, data.Filename)
	err = os.WriteFile(outputPath, formattedCode, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Сгененирован файл %s\n", outputPath)
}

func (g *GenerateService) copyControls() {
	outputDir := g.getOutputDir()

	controls := []string{"control", "pushbutton_control", "range_control", "switch_control", "text_control", "value_control"}
	for _, control := range controls {
		src := fmt.Sprintf("templates/%s.go", control)
		data, err := embedFs.ReadFile(src)
		if err != nil {
			panic(err)
		}

		dst := fmt.Sprintf("%s/%s.go", outputDir, control)
		err = os.WriteFile(dst, data, 0644)
		if err != nil {
			panic(err)
		}
	}
}
