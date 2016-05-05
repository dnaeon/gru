package command

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/dnaeon/gru/catalog"
	"github.com/dnaeon/gru/module"
)

// NewGraphCommand creates a new sub-command for
// generating the resource DAG graph
func NewGraphCommand() cli.Command {
	cmd := cli.Command{
		Name:   "graph",
		Usage:  "create DOT representation for modules and resources",
		Action: execGraphCommand,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "modules",
				Usage: "create DOT for module imports",
			},
			cli.BoolFlag{
				Name:  "resources",
				Usage: "create DOT for resources in catalog",
			},
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

// Executes the "graph" command
func execGraphCommand(c *cli.Context) {
	if len(c.Args()) < 1 {
		displayError(errNoModuleName, 64)
	}

	if !c.Bool("modules") && !c.Bool("resources") {
		displayError(errors.New("Must specify either --modules or --resources flag"), 64)
	}

	if c.Bool("modules") && c.Bool("resources") {
		displayError(errors.New("Only one of --modules or --resources can be specified"), 64)
	}

	main := c.Args()[0]
	modulePath := filepath.Join(c.String("sitedir"), "modules")

	// Create DOT representation for module imports
	if c.Bool("modules") {
		if err := module.ImportGraphAsDot(main, modulePath, os.Stdout); err != nil {
			displayError(err, 1)
		}
	} else if c.Bool("resources") {
		// Create DOT representation for resources
		katalog, err := catalog.Load(main, modulePath)
		if err != nil {
			displayError(err, 1)
		}

		err = katalog.ResourcesAsDot(os.Stdout)
		if err != nil {
			displayError(err, 1)
		}
	}
}
