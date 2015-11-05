package command

import (
	"github.com/codegangsta/cli"
)

func NewListCommand() cli.Command {
	cmd := cli.Command{
		Name: "list",
		Usage: "list registered minions",
		Action: runListCommand,
	}

	return cmd
}

// Executes the "list" command
func runListCommand(c *cli.Context) {
	client := newEtcdMinionClientFromFlags(c)
	minions, err := client.MinionList()

	if err != nil {
		displayError(err, 1)
	}

	for _, minion := range minions {
		fmt.Println(minion)
	}
}
