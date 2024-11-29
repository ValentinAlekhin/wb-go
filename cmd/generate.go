/*
Copyright © 2024 Valentin Alekhin <alekhin.dev@yandex.ru>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"wb-go/internal/gen"
	"wb-go/pkg/mqtt"
)

var broker string
var output string
var devices []string
var packageName string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Генерирует код для работы с устройствами",
	Long:  `Команда обращается к MQTT-топикам и автоматически генерирует код на Golang для взаимодействия с устройствами.`,
	Run: func(cmd *cobra.Command, args []string) {
		opt := mqtt.Options{
			Broker:   "192.168.1.150:1883",
			ClientId: "wb-go-generator",
		}
		client := mqtt.NewClient(opt)

		generateService := gen.NewGenerateService(client, output, devices, packageName)
		generateService.Run()
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&broker, "broker", "b", "", "Адрес MQTT брокера")
	generateCmd.Flags().StringVarP(&output, "output", "o", "", "Директория, куда будут сгенерированы файлы")
	generateCmd.Flags().StringArrayVarP(&devices, "devices", "d", gen.DevicesToGenerate, "Имена устройств для генерации")
	generateCmd.Flags().StringVarP(&packageName, "package", "p", "devices", "Имя пакета сгенерированных файлов")
	err := generateCmd.MarkFlagRequired("broker")
	if err != nil {
		panic(err)
	}
}
