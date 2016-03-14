package gructl

import (
	"os"
	"time"

	"github.com/dnaeon/gru/gructl/command"
	"github.com/dnaeon/gru/version"

	"github.com/codegangsta/cli"
)

// Main is the entry point of gructl
func Main() {
	app := cli.NewApp()
	app.Name = "gructl"
	app.Version = version.Version
	app.Usage = "command line tool for managing minions"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "endpoint",
			Value:  "http://127.0.0.1:2379,http://localhost:4001",
			Usage:  "etcd cluster endpoints",
			EnvVar: "GRUCTL_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "username",
			Value:  "",
			Usage:  "username to use for authentication",
			EnvVar: "GRUCTL_USERNAME",
		},
		cli.StringFlag{
			Name:   "password",
			Value:  "",
			Usage:  "password to use for authentication",
			EnvVar: "GRUCTL_PASSWORD",
		},
		cli.StringFlag{
			Name:   "modulepath",
			Value:  "",
			Usage:  "path to modules",
			EnvVar: "GRU_MODULEPATH",
		},
		cli.DurationFlag{
			Name:  "timeout",
			Value: time.Second,
			Usage: "connection timeout per request",
		},
	}

	app.Commands = []cli.Command{
		command.NewApplyCommand(),
		command.NewListCommand(),
		command.NewInfoCommand(),
		command.NewServeCommand(),
		command.NewRunCommand(),
		command.NewClassifierCommand(),
		command.NewReportCommand(),
		command.NewQueueCommand(),
		command.NewLogCommand(),
		command.NewLastseenCommand(),
		command.NewResultCommand(),
		command.NewGraphCommand(),
		command.NewValidateCommand(),
	}

	app.Run(os.Args)
}
