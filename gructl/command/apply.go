package command

import (
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/dnaeon/gru/catalog"
	"github.com/dnaeon/gru/module"
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

	config := &catalog.Config{
		Main:   c.Args()[0],
		DryRun: c.Bool("dry-run"),
		ModuleConfig: &module.Config{
			Path: filepath.Join(c.String("sitedir"), "modules"),
			ResourceConfig: &resource.Config{
				SiteDir: c.String("sitedir"),
				Writer:  os.Stdout,
			},
		},
	}
	katalog, err := catalog.Load(config)
	if err != nil {
		displayError(err, 1)
	}

	err = katalog.Run()
	if err != nil {
		displayError(err, 1)
	}
}
