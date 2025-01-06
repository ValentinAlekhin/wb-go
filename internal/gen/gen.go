package gen

import (
	"bytes"
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iancoleman/strcase"
	"go/format"
	"os"
	"slices"
	"strings"
	"text/template"
	"time"
	"unicode"
)

type GenerateService struct {
	client      *wb.Client
	outputDir   string
	packageName string
}

type deviceControlTemplateData struct {
	Name       string
	Mqtt       string
	ReadOnly   bool
	Type       string
	StructName string
	Meta       control.Meta
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
	Meta          control.Meta
}

//go:embed templates/*
var embedFs embed.FS

func NewGenerateService(client *wb.Client, output string, packageName string) *GenerateService {
	service := &GenerateService{
		client:      client,
		outputDir:   output,
		packageName: packageName,
	}
	return service
}

func (g *GenerateService) Run() {
	watchResults := g.collectData()
	filteredResults := g.filterUnique(watchResults)
	templatesData := g.generateTemplates(filteredResults)
	g.generateFiles(templatesData)
}

func (g *GenerateService) collectData() []watchResultItem {
	var list []watchResultItem

	watcher := g.getTopicWatcher(&list)
	g.client.Subscribe(mqttTopic, watcher)
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

		controlMeta := control.Meta{}

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

func (g *GenerateService) generateTemplates(list []watchResultItem) map[string]*deviceTemplateData {
	deviceMap := map[string]*deviceTemplateData{}

	for _, item := range list {
		key := fmt.Sprintf("%s_%s", item.DeviceName, item.ModbusAddress)
		if val, ok := deviceMap[key]; !ok {
			deviceStructName := strcase.ToCamel(item.DeviceName + item.ModbusAddress)
			deviceControlsStructName := deviceStructName + "controls"

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

func (g *GenerateService) getControlTemplate(control watchResultItem) deviceControlTemplateData {
	typeName := control.Meta.Type
	for key, val := range controlValueTypeMap {
		if !slices.Contains(val, control.Meta.Type) {
			continue
		}

		typeName = key
	}

	name := strcase.ToCamel(control.Control)
	if unicode.IsDigit(rune(name[0])) {
		name = "C" + name
	}

	return deviceControlTemplateData{
		Name:       name,
		Mqtt:       control.Control,
		ReadOnly:   control.Meta.ReadOnly,
		Type:       typeName,
		StructName: strcase.ToCamel(typeName) + "Control",
		Meta:       control.Meta,
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

func (g *GenerateService) generateFiles(data map[string]*deviceTemplateData) {
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
	tmpl, err := template.ParseFS(embedFs, deviceTemplateFile)
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
	}

	outputPath := fmt.Sprintf("%s/%s.go", outputDir, data.Filename)
	err = os.WriteFile(outputPath, formattedCode, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Сгененирован файл %s\n", outputPath)
}
