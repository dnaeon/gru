package command

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/dnaeon/gru/catalog"
)

// NewGraphCommand creates a new sub-command for
// generating the resource DAG graph
func NewGraphCommand() cli.Command {
	cmd := cli.Command{
		Name:   "graph",
		Usage:  "generate a DOT file of the resources graph",
		Action: execGraphCommand,
	}

	return cmd
}

// Executes the "graph" command
func execGraphCommand(c *cli.Context) {
	if len(c.Args()) < 1 {
		displayError(errNoModuleName, 64)
	}

	main := c.Args()[0]
	katalog, err := catalog.Load(main, c.GlobalString("modulepath"))
	if err != nil {
		displayError(err, 1)
	}

	err = katalog.GenerateCatalogDOT(os.Stdout)
	if err != nil {
		displayError(err, 1)
	}
}
