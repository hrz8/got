package main

import (
	"os"

	"github.com/hrz8/got/cmd"
)

func main() {
	root := cmd.NewRootCommand()

	version := cmd.NewVersionCommand()
	root.AddCommand(version)

	serve := cmd.NewServeCommand()
	root.AddCommand(serve)

	migrate := cmd.NewMigrateCommand()
	root.AddCommand(migrate)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
