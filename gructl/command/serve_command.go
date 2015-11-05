package command

import (
	"os"

	"github.com/codegangsta/cli"
)

func NewServeCommand() cli.Command {
	cmd := cli.Command{
		Name: "serve",
		Usage: "start minion",
		Action: runServeCommand,
	}

	return cmd
}

// Executes the "serve" command
func runServeCommand(c *cli.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		displayError(err, 1)
	}

	cfg := etcdConfigFromFlags(c)
	m := minion.NewEtcdMinion(hostname, cfg)
	m.Serve()
}
