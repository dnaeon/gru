package command

import (
	"encoding/json"

	"github.com/codegangsta/cli"
	"github.com/dnaeon/gru/catalog"
)

// NewValidateCommand creates a new sub-command for
// validating module files
func NewValidateCommand() cli.Command {
	cmd := cli.Command{
		Name:   "validate",
		Usage:  "validate module file",
		Action: execValidateCommand,
	}

	return cmd
}

// Executes the "validate" command
func execValidateCommand(c *cli.Context) {
	if len(c.Args()) < 1 {
		displayError(errNoModuleFile, 64)
	}

	moduleFile := c.Args()[0]
	katalog, err := catalog.Load(moduleFile)
	if err != nil {
		displayError(err, 1)
	}

	_, err = json.Marshal(katalog)
	if err != nil {
		displayError(err, 1)
	}
}
