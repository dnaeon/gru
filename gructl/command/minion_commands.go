package command

import (
	"fmt"
	"time"
	"errors"

	"code.google.com/p/go-uuid/uuid"
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
			{
				Name: "info",
				Usage: "get info about a minion",
				Action: minionInfoCommand,
			},
		},
	}

	return cmd
}

// Executes the "minion list" command
func minionListCommand(c *cli.Context) {
	client := newEtcdMinionClientFromFlags(c)
	minions, err := client.MinionList()

	if err != nil {
		displayError(err, 1)
	}

	for _, minion := range minions {
		fmt.Println(minion)
	}
}

// The "minion info" command
func minionInfoCommand(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayError(errors.New("Must provide a minion uuid"), 64)
	}

	arg := c.Args()[0]
	minion := uuid.Parse(arg)
	if minion == nil {
		displayError(errors.New("Bad minion uuid given"), 64)
	}

	client := newEtcdMinionClientFromFlags(c)
	name, err := client.MinionName(minion)
	if err != nil {
		displayError(err, 1)
	}

	lastseen, err := client.MinionLastseen(minion)
	if err != nil {
		displayError(err, 1)
	}

	taskQueue, err := client.MinionTaskQueue(minion)
	if err != nil {
		displayError(err, 1)
	}

	taskLog, err := client.MinionTaskLog(minion)
	if err != nil {
		displayError(err, 1)
	}

	fmt.Printf("%-15s: %s\n", "Minion", minion)
	fmt.Printf("%-15s: %s\n", "Name", name)
	fmt.Printf("%-15s: %s\n", "Lastseen", time.Unix(lastseen, 0))
	fmt.Printf("%-15s: %d task(s)\n", "Queue", len(taskQueue))
	fmt.Printf("%-15s: %d task(s)\n", "Processed", len(taskLog))
}
