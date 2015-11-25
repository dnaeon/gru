package command

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/dnaeon/gru/minion"
)

func NewServeCommand() cli.Command {
	cmd := cli.Command{
		Name:   "serve",
		Usage:  "start minion",
		Action: execServeCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "name",
				Usage: "set minion name",
				Value: "",
			},
		},
	}

	return cmd
}

// Executes the "serve" command
func execServeCommand(c *cli.Context) {
	var name string

	name, err := os.Hostname()
	if err != nil {
		displayError(err, 1)
	}

	nameFlag := c.String("name")
	if nameFlag != "" {
		name = nameFlag
	}

	cfg := etcdConfigFromFlags(c)
	m := minion.NewEtcdMinion(name, cfg)
	m.Serve()
}
