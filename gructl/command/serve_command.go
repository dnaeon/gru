package command

import (
	"os"

	"github.com/dnaeon/gru/minion"
	"github.com/codegangsta/cli"
)

func NewServeCommand() cli.Command {
	cmd := cli.Command{
		Name: "serve",
		Usage: "start minion",
		Action: execServeCommand,
	}

	return cmd
}

// Executes the "serve" command
func execServeCommand(c *cli.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		displayError(err, 1)
	}

	cfg := etcdConfigFromFlags(c)
	m := minion.NewEtcdMinion(hostname, cfg)
	m.Serve()
}
