package command

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/dnaeon/gru/catalog"
	"github.com/dnaeon/gru/module"
	"github.com/dnaeon/gru/resource"
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
				Name:   "siterepo",
				Value:  "",
				Usage:  "path/url to the site repo",
				EnvVar: "GRU_SITEREPO",
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
	modulePath := filepath.Join(c.String("siterepo"), "modules")

	// Create DOT representation for module imports
	if c.Bool("modules") {
		if err := module.ImportGraphAsDot(main, modulePath, os.Stdout); err != nil {
			displayError(err, 1)
		}
	} else if c.Bool("resources") {
		// Create DOT representation for resources from catalog
		config := &catalog.Config{
			Main:   main,
			DryRun: true,
			ModuleConfig: &module.Config{
				Path: modulePath,
				ResourceConfig: &resource.Config{
					SiteRepo: c.String("siterepo"),
					Writer:   os.Stdout,
				},
			},
		}

		katalog, err := catalog.Load(config)
		if err != nil {
			displayError(err, 1)
		}

		collection, err := module.ResourceCollection(katalog.Modules)
		if err != nil {
			displayError(err, 1)
		}

		if err := collection.DependencyGraphAsDot(os.Stdout); err != nil {
			displayError(err, 1)
		}
	}
}
