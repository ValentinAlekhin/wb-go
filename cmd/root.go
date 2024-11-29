/*
Copyright © 2024 Valentin Alekhin <alekhin.dev@yandex.ru>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "wb-go",
	Short: "CLI утилита для генерации Golang-кода управления устройствами Wiren Board.",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}
