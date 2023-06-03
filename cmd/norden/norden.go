package main

import (
	"os"

	"github.com/kris-nova/logger"

	"github.com/xrelkd/norden/pkg/commands/root"
)

func main() {
	cmd := root.Command()

	if err := cmd.Execute(); err != nil {
		logger.Debug("%s", err.Error())
		os.Exit(1)
	}
}
