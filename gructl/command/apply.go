package command

import (
	"fmt"
	"os"

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
	}

	return cmd
}

// Executes the "apply" command
func execApplyCommand(c *cli.Context) {
	if len(c.Args()) < 1 {
		displayError(errNoModuleFile, 64)
	}

	main := c.Args()[0]
	katalog, err := catalog.Load(main, c.GlobalString("modulepath"))
	if err != nil {
		displayError(err, 1)
	}

	fmt.Printf("Loaded %d resource(s) in catalog\n", katalog.Len())
	err = katalog.Run(os.Stdout)
	if err != nil {
		displayError(err, 1)
	}
}
