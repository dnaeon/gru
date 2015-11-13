package command

import (
	"fmt"
	"time"

	"github.com/codegangsta/cli"
)

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
func execLastseenCommand(c *cli.Context) {
	client := newEtcdMinionClientFromFlags(c)

	cFlag := c.String("with-classifier")
	minions, err := parseClassifierPattern(client, cFlag)

	if err != nil {
		displayError(err, 1)
	}

	for _, minion := range minions {
		lastseen, err := client.MinionLastseen(minion)
		if err != nil {
			displayError(err, 1)
		}

		fmt.Printf("%s\t%s\n", minion, time.Unix(lastseen, 0))
	}
}
