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
				Name:   "sitedir",
				Value:  "",
				Usage:  "specify path to the site directory",
				EnvVar: "GRU_SITEDIR",
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

	main := c.Args()[0]
	modulePath := filepath.Join(c.String("sitedir"), "modules")

	katalog, err := catalog.Load(main, modulePath)
	if err != nil {
		displayError(err, 1)
	}

	collection, err := module.ResourceCollection(katalog.Modules)
	if err != nil {
		displayError(err, 1)
	}

	opts := &resource.Options{
		SiteDir: c.String("sitedir"),
	}

	fmt.Printf("Loaded %d resources from %d modules\n", len(collection), len(katalog.Modules))
	for _, m := range katalog.Modules {
		for _, key := range m.UnknownKeys {
			fmt.Printf("Uknown key '%s' in module '%s'\n", key, m.Name)
		}
	}

	for _, r := range collection {
		if _, err := r.Evaluate(os.Stdout, opts); err != nil {
			fmt.Printf("Resource %s: %s\n", r.ResourceID(), err)
		}
	}
}
