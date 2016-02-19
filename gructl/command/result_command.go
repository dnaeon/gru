package command

import (
	"fmt"
	"time"

	"github.com/codegangsta/cli"
	"github.com/gosuri/uitable"
	"github.com/pborman/uuid"
)

// NewResultCommand creates a new sub-command for retrieving
// results of previously executed tasks by minions
func NewResultCommand() cli.Command {
	cmd := cli.Command{
		Name:   "result",
		Usage:  "get task results",
		Action: execResultCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "minion",
				Usage: "get task result for given minion only",
			},
			cli.BoolFlag{
				Name:  "details",
				Usage: "provide more details about the tasks",
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
		return
	}

	// Create table for the task results
	// If the --details flag is specified, then
	// create a table that holds all details about the
	// tasks, otherwise use a simple summary table
	table := uitable.New()
	if c.Bool("details") {
		table.MaxColWidth = 80
		table.Wrap = true
	} else {
		table.MaxColWidth = 40
		table.AddRow("MINION", "RESULT", "STATE")
	}

	for _, minion := range minionWithTask {
		task, err := client.MinionTaskResult(minion, taskId)
		if err != nil {
			displayError(err, 1)
		}

		if c.Bool("details") {
			table.AddRow("Minion:", minion)
			table.AddRow("Task ID:", task.TaskID)
			table.AddRow("State:", task.State)
			table.AddRow("Command:", task.Command)
			table.AddRow("Args:", task.Args)
			table.AddRow("Concurrent:", task.IsConcurrent)
			table.AddRow("Received:", time.Unix(task.TimeReceived, 0))
			table.AddRow("Processed:", time.Unix(task.TimeProcessed, 0))
			table.AddRow("Result:", task.Result)
			table.AddRow("Error:", task.Error)
		} else {
			table.AddRow(minion, task.Result, task.State)
		}
	}

	fmt.Println(table)
}
