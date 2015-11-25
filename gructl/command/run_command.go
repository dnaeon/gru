package command

import (
	"fmt"
	"time"

	"github.com/dnaeon/gru/task"

	"github.com/codegangsta/cli"
	"github.com/gosuri/uitable"
	"github.com/gosuri/uiprogress"
)

func NewRunCommand() cli.Command {
	cmd := cli.Command{
		Name:   "run",
		Usage:  "send task to minion(s)",
		Action: execRunCommand,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "is-concurrent",
				Usage: "flag task as concurrent",
			},
			cli.StringFlag{
				Name:  "with-classifier",
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
		displayError(errMissingTask, 64)
	}

	client := newEtcdMinionClientFromFlags(c)

	cFlag := c.String("with-classifier")
	minions, err := parseClassifierPattern(client, cFlag)

	if err != nil {
		displayError(err, 1)
	}

	numMinions := len(minions)
	if numMinions == 0 {
		displayError(errNoMinionFound, 1)
	}

	fmt.Printf("Found %d minion(s) for task processing\n\n", numMinions)

	// Create the task that we send to our minions
	// The first argument is the command and anything else
	// that follows is considered task arguments
	args := c.Args()
	isConcurrent := c.Bool("is-concurrent")
	taskCommand := args[0]
	taskArgs := args[1:]
	t := task.New(taskCommand, taskArgs...)
	t.IsConcurrent = isConcurrent

	// Progress bar to display while submitting task
	progress := uiprogress.New()
	bar := progress.AddBar(numMinions)
	bar.AppendCompleted()
	bar.PrependElapsed()
	progress.Start()

	// Number of minions to which submitting the task has failed
	failed := 0

	// Submit task to minions
	fmt.Println("Submitting task to minion(s) ...")
	for _, minion := range minions {
		err = client.MinionSubmitTask(minion, t)
		if err != nil {
			fmt.Printf("Failed to submit task to %s: %s\n", minion, err)
			failed += 1
		}
		bar.Incr()
	}

	// Stop progress bar and sleep for a bit to make sure the
	// progress bar gets updated if we were too fast for it
	progress.Stop()
	time.Sleep(time.Millisecond * 100)

	// Display task report
	fmt.Println()
	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true
	table.AddRow("TASK", "SUBMITTED", "FAILED", "TOTAL")
	table.AddRow(t.TaskID, numMinions-failed, failed, numMinions)
	fmt.Println(table)
}
