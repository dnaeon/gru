package command

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/dnaeon/gru/minion"
)

// NewServeCommand creates a new sub-command for starting a
// minion and its services
func NewServeCommand() cli.Command {
	cmd := cli.Command{
		Name:   "serve",
		Usage:  "start minion",
		Action: execServeCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: "set minion name",
				Value: "",
			},
		},
	}

	return cmd
}

// Executes the "serve" command
func execServeCommand(c *cli.Context) {
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

	// Channel on which the shutdown signal is sent
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start minion
	err = m.Serve()
	if err != nil {
		displayError(err, 1)
	}

	// Block until a shutdown signal is received
	<-quit
	m.Stop()
}
