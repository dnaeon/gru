package command

import (
	"fmt"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/cli"
	"github.com/gosuri/uitable"
)

func NewResultCommand() cli.Command {
	cmd := cli.Command{
		Name: "result",
		Usage: "get task results",
		Action: execResultCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "minion",
				Usage: "get task result for given minion only",
			},
		},
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

	// If --minion flag was specified parse the
	// minion uuid and get the task result only
	// from the specified minion, otherwise find
	// all minions which contain the given
	// task and get their results
	var minionWithTask []uuid.UUID

	mFlag := c.String("minion")
	if mFlag == "" {
		// No minion was specified, get all minions
		// with the given task uuid
		m, err := client.MinionWithTaskResult(taskId)
		if err != nil {
			displayError(err, 1)
		}
		minionWithTask = m
	} else {
		// Minion was specified, get task result
		// from the given minion only
		minion := uuid.Parse(mFlag)
		if minion == nil {
			displayError(errInvalidUUID, 64)
		}
		minionWithTask = append(minionWithTask, minion)
	}

	if len(minionWithTask) == 0 {
		displayError(errNoMinionFound, 1)
	}

	// Create table for the task results
	table := uitable.New()
	table.MaxColWidth = 60
	table.AddRow("MINION", "RESULT")

	for _, minion := range minionWithTask {
		task, err := client.MinionTaskResult(minion, taskId)
		if err != nil {
			displayError(err, 1)
		}

		table.AddRow(minion, task.Result)
	}

	fmt.Println(table)
}
