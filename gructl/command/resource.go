package command

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/dnaeon/gru/resource"
	"github.com/gosuri/uitable"
)

// NewResourceCommand creates a new sub-command for
// displaying the list of registered resource types
func NewResourceCommand() cli.Command {
	cmd := cli.Command{
		Name:   "resource",
		Usage:  "display registered resources",
		Action: execResourceCommand,
	}

	return cmd
}

// Executes the "resource" command
func execResourceCommand(c *cli.Context) {
	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("RESOURCE", "DESCRIPTION")

	for _, item := range resource.Registry {
		table.AddRow(item.Name, item.Description)
	}

	fmt.Println(table)
}
