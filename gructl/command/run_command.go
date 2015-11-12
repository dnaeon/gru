package command

import (
	"fmt"
	"errors"

	"github.com/dnaeon/gru/task"

	"github.com/codegangsta/cli"
)

func NewRunCommand() cli.Command {
	cmd := cli.Command{
		Name: "run",
		Usage: "send task to minion(s)",
		Action: execRunCommand,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "is-concurrent",
				Usage: "flag task as concurrent",
			},
			cli.StringFlag{
				Name: "with-classifier",
				Value: "",
				Usage: "match minions with given classifier pattern",
			},
		},
	}

	return cmd
}

// Executes the "run" command
func execRunCommand(c *cli.Context) {
	if len(c.Args()) < 1 {
		displayError(errors.New("Must provide command to run"), 64)
	}

	client := newEtcdMinionClientFromFlags(c)

	cFlag := c.String("with-classifier")
	minions, err := parseClassifierPattern(client, cFlag)

	if err != nil {
		displayError(err, 1)
	}

	numMinions := len(minions)
	if numMinions == 0 {
		displayError(errors.New("No minion(s) found"), 1)
	}

	fmt.Printf("Found %d minion(s) for task processing\n", numMinions)

	// The first argument is the command and anything else
	// that follows is considered as task arguments
	args := c.Args()
	isConcurrent := c.Bool("is-concurrent")
	taskCommand := args[0]
	taskArgs := args[1:]
	t := task.New(taskCommand, taskArgs...)
	t.IsConcurrent = isConcurrent

	failed := 0
	for i, minion := range minions {
		fmt.Printf("[%d/%d] Submitting task to minion %s\r", i + 1, numMinions, minion)
		err = client.MinionSubmitTask(minion, t)
		if err != nil {
			failed += 1
			fmt.Printf("\nFailed to submit task to %s: %s\n", minion, err)
		}
	}
	fmt.Println()

	fmt.Printf("Task submitted to %d minion(s), %d of which has failed\n", numMinions, failed)
	fmt.Printf("Task results can be retrieved by using this task id: %s\n", t.TaskID)
}
