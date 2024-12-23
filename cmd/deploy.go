package cmd

import (
	"github.com/ValentinAlekhin/wb-go/internal/deploy"
	"github.com/spf13/cobra"
	"log"
)

var config string

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Развертывает приложение на целевых устройствах",
	Long:  `Команда обращается к файлу конфигурации. Собирает бинарный файл, подготавливает systemd файл, загружаем все на устройства и запускает приложение`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := deploy.GetConfig(config)
		if err != nil {
			log.Fatal(err)
		}
		err = deploy.Run(cfg)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&config, "config", "c", "", "Путь к файлу конфигурации")
}
