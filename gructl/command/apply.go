package command

import (
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/dnaeon/gru/catalog"
	"github.com/dnaeon/gru/resource"
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
			cli.BoolFlag{
				Name:  "dry-run",
				Usage: "just report what would be done, instead of doing it",
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

	opts := &resource.Options{
		SiteDir: c.String("sitedir"),
		DryRun:  c.Bool("dry-run"),
	}

	err = katalog.Run(os.Stdout, opts)
	if err != nil {
		displayError(err, 1)
	}
}
