package main

import (
	"os"
	"time"

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
		cli.StringFlag{
			Name: "endpoint",
			Value: "http://127.0.0.1:2379,http://localhost:4001",
			Usage: "etcd cluster endpoints",
			EnvVar: "GRUCTL_ENDPOINT",
		},
		cli.StringFlag{
			Name: "username",
			Value: "",
			Usage: "username to use for authentication",
			EnvVar: "GRUCTL_USERNAME",
		},
		cli.StringFlag{
			Name: "password",
			Value: "",
			Usage: "password to use for authentication",
			EnvVar: "GRUCTL_PASSWORD",
		},
		cli.DurationFlag{
			Name: "timeout",
			Value: time.Second,
			Usage: "connection timeout per request",
		},
	}

	app.Commands = []cli.Command{
		command.NewListCommand(),
		command.NewInfoCommand(),
		command.NewServeCommand(),
		command.NewRunCommand(),
		command.NewClassifierCommand(),
		command.NewReportCommand(),
		command.NewQueueCommand(),
	}

	app.Run(os.Args)
}
