package command

import (
	"fmt"
	"time"

	"github.com/gosuri/uitable"
	"github.com/urfave/cli"
)

// NewLastseenCommand creates a new sub-command for
// retrieving the last time minions were seen
func NewLastseenCommand() cli.Command {
	cmd := cli.Command{
		Name:   "lastseen",
		Usage:  "show when minion(s) were last seen",
		Action: execLastseenCommand,
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

// Executes the "lastseen" command
func execLastseenCommand(c *cli.Context) error {
	client := newEtcdMinionClientFromFlags(c)

	cFlag := c.String("with-classifier")
	minions, err := parseClassifierPattern(client, cFlag)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("MINION", "LASTSEEN")
	for _, minion := range minions {
		lastseen, err := client.MinionLastseen(minion)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		table.AddRow(minion, time.Unix(lastseen, 0))
	}

	fmt.Println(table)

	return nil
}
