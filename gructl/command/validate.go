package command

import (
	"fmt"

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
		displayError(errNoModuleName, 64)
	}

	main := c.Args()[0]
	katalog, err := catalog.Load(main, c.GlobalString("modulepath"))
	if err != nil {
		displayError(err, 1)
	}

	// Validate() returns a slice of errors
	foundErrors := katalog.Validate()
	for _, err = range foundErrors {
		fmt.Println(err)
	}

	if len(foundErrors) > 0 {
		return
	}
}
