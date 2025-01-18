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
	"sync"
	"text/template"
	"time"
	"unicode"
)

type GenerateService struct {
	client      wb.ClientInterface
	outputDir   string
	packageName string
}

type deviceMeta struct {
	Driver string `json:"driver"`
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
	Controls                 []deviceControlTemplateData
	PackageName              string
}

type watchDeviceItem struct {
	Meta deviceMeta
	Name string
}

type watchControlResultItem struct {
	DeviceName string
	Control    string
	Meta       control.Meta
}

//go:embed templates/*
var embedFs embed.FS

func NewGenerateService(client wb.ClientInterface, output string, packageName string) *GenerateService {
	service := &GenerateService{
		client:      client,
		outputDir:   output,
		packageName: packageName,
	}
	return service
}

func (g *GenerateService) Run() {
	deviceData := g.collectDevicesData()
	watchResults := g.collectControlsData()
	filteredResults := g.filter(watchResults, deviceData)
	templatesData := g.generateTemplates(filteredResults)
	g.generateFiles(templatesData)
}

func (g *GenerateService) collectDevicesData() []watchDeviceItem {
	list := make([]watchDeviceItem, 0)
	devChan := make(chan watchDeviceItem)

	duration := 100 * time.Millisecond
	timer := time.NewTimer(duration)
	defer timer.Stop()

	watcher := g.getDeviceMetaWatcher(devChan)
	_ = g.client.Subscribe(deviceMetaTopic, watcher)

	for {
		select {
		case item := <-devChan:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(duration)
			list = append(list, item)
		case <-timer.C:
			return list
		}
	}
}

func (g *GenerateService) getDeviceMetaWatcher(ch chan<- watchDeviceItem) func(client mqtt.Client, msg mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		metaStr := string(msg.Payload())

		topicParts := strings.Split(topic, "/")
		name := topicParts[2]

		meta := deviceMeta{}
		err := json.Unmarshal([]byte(metaStr), &meta)
		if err != nil {
			panic(err)
		}

		ch <- watchDeviceItem{
			Meta: meta,
			Name: name,
		}
	}
}

func (g *GenerateService) collectControlsData() []watchControlResultItem {
	var list []watchControlResultItem
	controlCh := make(chan watchControlResultItem)

	duration := 100 * time.Millisecond
	timer := time.NewTimer(duration)
	defer timer.Stop()

	watcher := g.getControlMetaWatcher(controlCh)
	_ = g.client.Subscribe(controlMetaTopic, watcher)

	for {
		select {
		case item := <-controlCh:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(duration)
			list = append(list, item)
		case <-timer.C:
			return list
		}
	}
}

func (g *GenerateService) getControlMetaWatcher(ch chan<- watchControlResultItem) func(client mqtt.Client, msg mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		meta := string(msg.Payload())

		topicParts := strings.Split(topic, "/")
		device := topicParts[2]
		controlName := topicParts[4]
		controlMeta := control.Meta{}

		err := json.Unmarshal([]byte(meta), &controlMeta)
		if err != nil {
			panic(err)
		}

		ch <- watchControlResultItem{
			DeviceName: device,
			Control:    controlName,
			Meta:       controlMeta,
		}

	}
}

func (g *GenerateService) generateTemplates(list []watchControlResultItem) map[string]deviceTemplateData {
	deviceMap := map[string]deviceTemplateData{}

	for _, item := range list {
		key := item.DeviceName
		if val, ok := deviceMap[key]; !ok {
			deviceStructName := strcase.ToCamel(item.DeviceName)
			deviceControlsStructName := deviceStructName + "Controls"

			filename := item.DeviceName

			controlTemplate := g.getControlTemplate(item)
			newVal := deviceTemplateData{
				DeviceName:               item.DeviceName,
				DeviceStructName:         deviceStructName,
				DeviceControlsStructName: deviceControlsStructName,
				Filename:                 filename,
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

func (g *GenerateService) getControlTemplate(control watchControlResultItem) deviceControlTemplateData {
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

func (g *GenerateService) filter(list []watchControlResultItem, devices []watchDeviceItem) []watchControlResultItem {
	uniqueMap := make(map[string]struct{})
	var result []watchControlResultItem

	devicesToIgnore := make(map[string]struct{})
	for _, device := range devices {
		if slices.Contains(driversToExclude, device.Meta.Driver) {
			devicesToIgnore[device.Name] = struct{}{}
		}
	}

	for _, item := range list {
		if _, ok := devicesToIgnore[item.DeviceName]; ok {
			continue
		}

		key := item.DeviceName + "|" + item.Control
		if _, exists := uniqueMap[key]; !exists {
			uniqueMap[key] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func (g *GenerateService) generateFiles(data map[string]deviceTemplateData) {
	outputDir := g.getOutputDir()

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(data))

	for _, v := range data {
		go func() {
			g.generateFile(v, outputDir)
			wg.Done()
		}()
	}

	wg.Wait()
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

func (g *GenerateService) generateFile(data deviceTemplateData, outputDir string) {
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
		fmt.Println("Formating error:", err)
		fmt.Println("Code: \n", buf.String())
	}

	outputPath := fmt.Sprintf("%s/%s.go", outputDir, data.Filename)
	err = os.WriteFile(outputPath, formattedCode, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("File generated %s\n", outputPath)
}
