package gructl

import (
	"os"
	"time"

	"github.com/dnaeon/gru/gructl/command"
	"github.com/dnaeon/gru/version"
	"github.com/urfave/cli"
)

// Main is the entry point of gructl
func Main() {
	app := cli.NewApp()
	app.Name = "gructl"
	app.Version = version.Version
	app.EnableBashCompletion = true
	app.Usage = "command line tool for managing minions"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "endpoint",
			Value:  "http://127.0.0.1:2379,http://localhost:4001",
			Usage:  "etcd cluster endpoints",
			EnvVar: "GRU_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "username",
			Value:  "",
			Usage:  "username to use for authentication",
			EnvVar: "GRU_USERNAME",
		},
		cli.StringFlag{
			Name:   "password",
			Value:  "",
			Usage:  "password to use for authentication",
			EnvVar: "GRU_PASSWORD",
		},
		cli.DurationFlag{
			Name:   "timeout",
			Value:  time.Second,
			Usage:  "connection timeout per request",
			EnvVar: "GRU_TIMEOUT",
		},
	}

	app.Commands = []cli.Command{
		command.NewApplyCommand(),
		command.NewListCommand(),
		command.NewInfoCommand(),
		command.NewServeCommand(),
		command.NewPushCommand(),
		command.NewClassifierCommand(),
		command.NewReportCommand(),
		command.NewQueueCommand(),
		command.NewLogCommand(),
		command.NewLastseenCommand(),
		command.NewResultCommand(),
		command.NewGraphCommand(),
	}

	app.Run(os.Args)
}
