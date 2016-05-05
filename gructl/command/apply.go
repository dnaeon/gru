package command

import (
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/dnaeon/gru/catalog"
)

// NewApplyCommand creates a new sub-command for
// applying configurations on the local system
func NewApplyCommand() cli.Command {
	cmd := cli.Command{
		Name:   "apply",
		Usage:  "apply configuration on the local system",
		Action: execApplyCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "sitedir",
				Value:  "",
				Usage:  "specify path to the site directory",
				EnvVar: "GRU_SITEDIR",
			},
		},
	}

	return cmd
}

// Executes the "apply" command
func execApplyCommand(c *cli.Context) {
	if len(c.Args()) < 1 {
		displayError(errNoModuleName, 64)
	}

	main := c.Args()[0]
	modulePath := filepath.Join(c.String("sitedir"), "modules")
	katalog, err := catalog.Load(main, modulePath)
	if err != nil {
		displayError(err, 1)
	}

	err = katalog.Run(os.Stdout)
	if err != nil {
		displayError(err, 1)
	}
}
