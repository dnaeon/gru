package command

import (
	"fmt"
	"time"

	"github.com/dnaeon/gru/task"

	"github.com/gosuri/uiprogress"
	"github.com/gosuri/uitable"
	"github.com/urfave/cli"
)

// NewPushCommand creates a new sub-command for submitting
// tasks to minions
func NewPushCommand() cli.Command {
	cmd := cli.Command{
		Name:   "push",
		Usage:  "send task to minion(s)",
		Action: execPushCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "environment",
				Value:  "production",
				Usage:  "specify environment to use",
				EnvVar: "GRU_ENVIRONMENT",
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

// Executes the "push" command
func execPushCommand(c *cli.Context) error {
	if len(c.Args()) < 1 {
		return cli.NewExitError(errNoModuleName.Error(), 64)
	}

	// Create the task that we send to our minions
	// The task's command is the module name that will be
	// loaded and processed by the remote minions
	main := c.Args()[0]
	t := task.New(main, c.String("environment"))

	client := newEtcdMinionClientFromFlags(c)

	cFlag := c.String("with-classifier")
	minions, err := parseClassifierPattern(client, cFlag)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	numMinions := len(minions)
	if numMinions == 0 {
		return cli.NewExitError(errNoMinionFound.Error(), 1)
	}

	fmt.Printf("Found %d minion(s) for task processing\n\n", numMinions)

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
	for _, m := range minions {
		err = client.MinionSubmitTask(m, t)
		if err != nil {
			fmt.Printf("Failed to submit task to %s: %s\n", m, err)
			failed++
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
	table.AddRow(t.ID, numMinions-failed, failed, numMinions)
	fmt.Println(table)

	return nil
}
