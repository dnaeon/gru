package command

import (
	"log"

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

	log.Printf("Loaded %d resource(s) in catalog", katalog.Len())
	err = katalog.Run()
	if err != nil {
		displayError(err, 1)
	}
}
