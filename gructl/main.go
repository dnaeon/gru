package main

import (
	"os"

	"github.com/dnaeon/gru/version"
	"github.com/dnaeon/gru/gructl/command"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gructl"
	app.Version = version.Version
	app.Usage = "command line tool for managing minions"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "endpoint", Value: "http://127.0.0.1:2379,http://localhost:4001", Usage: "etcd cluster endpoints"},
	}
	app.Commands = []cli.Command{
		command.NewLsCommand(),
	}

	app.Run(os.Args)
}
