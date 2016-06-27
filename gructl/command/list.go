package command

import (
	"fmt"

	"github.com/coreos/etcd/client"
	"github.com/gosuri/uitable"
	"github.com/urfave/cli"
)

// NewListCommand creates a new sub-command for retrieving the
// list of registered minions
func NewListCommand() cli.Command {
	cmd := cli.Command{
		Name:   "list",
		Usage:  "list registered minions",
		Action: execListCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "with-classifier",
				Value: "",
				Usage: "match minions with given classifier pattern",
			},
		},
	}

	return cmd
}

// Executes the "list" command
func execListCommand(c *cli.Context) error {
	klient := newEtcdMinionClientFromFlags(c)

	cFlag := c.String("with-classifier")
	minions, err := parseClassifierPattern(klient, cFlag)

	// Ignore errors about missing minion directory
	if err != nil {
		if eerr, ok := err.(client.Error); !ok || eerr.Code != client.ErrorCodeKeyNotFound {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	if len(minions) == 0 {
		return nil
	}

	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("MINION", "NAME")
	for _, minion := range minions {
		name, err := klient.MinionName(minion)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		table.AddRow(minion, name)
	}

	fmt.Println(table)

	return nil
}
