package main

import (
	"github.com/spf13/cobra/doc"
	"os"
	"wb-go/cmd"
)

func main() {
	dir := "docs/data/cli"
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = doc.GenYamlTree(cmd.RootCmd, dir)
	if err != nil {
		panic(err)
	}
}
