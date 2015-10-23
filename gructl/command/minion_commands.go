package command

import (
	"os"
	"fmt"

	"github.com/codegangsta/cli"
)

func NewMinionCommands() cli.Command {
	cmd := cli.Command{
		Name: "minion",
		Usage: "manage minions",
		Subcommands: []cli.Command{
			{
				Name: "list",
				Usage: "list registered minions",
				Action: minionListCommand,
			},
		},
	}

	return cmd
}

func minionListCommand(c *cli.Context) {
	client := newEtcdMinionClientFromFlags(c)
	minions, err := client.MinionList()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, minion := range minions {
		fmt.Println(minion)
	}
}
