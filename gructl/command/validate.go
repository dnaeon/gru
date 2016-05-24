package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/dnaeon/gru/catalog"
	"github.com/dnaeon/gru/module"
	"github.com/dnaeon/gru/resource"
)

// NewValidateCommand creates a new sub-command for
// validating module files
func NewValidateCommand() cli.Command {
	cmd := cli.Command{
		Name:   "validate",
		Usage:  "validate module file",
		Action: execValidateCommand,
		Flags: []cli.Flag{
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

// Executes the "validate" command
func execValidateCommand(c *cli.Context) {
	if len(c.Args()) < 1 {
		displayError(errNoModuleName, 64)
	}

	config := &catalog.Config{
		Main:   c.Args()[0],
		DryRun: true,
		ModuleConfig: &module.Config{
			Path: filepath.Join(c.String("siterepo"), "modules"),
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

	fmt.Printf("Loaded %d resources from %d modules\n", len(collection), len(katalog.Modules))
	for _, m := range katalog.Modules {
		for _, key := range m.UnknownKeys {
			fmt.Printf("Uknown key '%s' in module '%s'\n", key, m.Name)
		}
	}
}
