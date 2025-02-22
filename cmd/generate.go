package cmd

import (
	"github.com/ValentinAlekhin/wb-go/internal/gen"
	"github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"github.com/spf13/cobra"
	"log"
)

var broker string
var username string
var password string
var output string
var packageName string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Генерирует код для работы с устройствами",
	Long:  `Команда обращается к MQTT-топикам и автоматически генерирует код на Golang для взаимодействия с устройствами.`,
	Run: func(cmd *cobra.Command, args []string) {
		opt := mqtt.Options{
			Broker:   broker,
			ClientId: "wb-go-generator",
			Username: username,
			Password: password,
		}
		client, err := mqtt.NewClient(opt)
		if err != nil {
			log.Fatal(err)
		}

		generateService := gen.NewGenerateService(client, output, packageName)
		generateService.Run()
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&broker, "broker", "b", "", "Адрес MQTT брокера")
	generateCmd.Flags().StringVarP(&username, "username", "u", "", "Имя пользователя")
	generateCmd.Flags().StringVarP(&password, "password", "p", "", "Пароль пользователя")
	generateCmd.Flags().StringVarP(&output, "output", "o", "", "Директория, куда будут сгенерированы файлы")
	generateCmd.Flags().StringVarP(&packageName, "package", "n", "device", "Имя пакета сгенерированных файлов")
	err := generateCmd.MarkFlagRequired("broker")
	if err != nil {
		panic(err)
	}
}
