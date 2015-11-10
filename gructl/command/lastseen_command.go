package command

import (
	"fmt"
	"time"

	"github.com/codegangsta/cli"
)

func NewLastseenCommand() cli.Command {
	cmd := cli.Command{
		Name: "lastseen",
		Usage: "show when minion(s) were last seen",
		Action: execLastseenCommand,
	}

	return cmd
}

// Executes the "lastseen" command
func execLastseenCommand(c *cli.Context) {
	client := newEtcdMinionClientFromFlags(c)
	minions, err := client.MinionList()

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
