package command

import (
	"fmt"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/cli"
)

func NewResultCommand() cli.Command {
	cmd := cli.Command{
		Name: "result",
		Usage: "get task results",
		Action: execResultCommand,
	}

	return cmd
}

// Executes the "result" command
func execResultCommand(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayError(errMissingTask, 64)
	}

	arg := c.Args()[0]
	taskId := uuid.Parse(arg)
	if taskId == nil {
		displayError(errInvalidUUID, 64)
	}

	client := newEtcdMinionClientFromFlags(c)
	minionWithTask, err := client.MinionWithTaskResult(taskId)
	if err != nil {
		displayError(err, 1)
	}

	for _, minion := range minionWithTask {
		task, err := client.MinionTaskResult(minion, taskId)
		if err != nil {
			displayError(err, 1)
		}

		fmt.Printf("Minion: %s\n", minion)
		fmt.Printf("Result: %s\n\n", task.Result)
	}
}
