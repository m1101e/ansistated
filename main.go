package main

import (
	"github.com/m1101e/ansistated/internal/cli"
	"github.com/m1101e/ansistated/internal/config"
	_ "github.com/m1101e/ansistated/internal/plugins"
)

func main() {
	cfg := &config.Config{}

	command := cli.NewRootCmd(cfg)

	command.Execute()
}
