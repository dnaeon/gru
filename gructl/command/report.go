package command

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/urfave/cli"
)

// NewReportCommand creates a new sub-command for
// generating reports based on minion classifiers
func NewReportCommand() cli.Command {
	cmd := cli.Command{
		Name:   "report",
		Usage:  "generate classifier report",
		Action: execReportCommand,
	}

	return cmd
}

// Executes the "report" command
func execReportCommand(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return cli.NewExitError(errNoClassifier.Error(), 64)
	}

	classifierKey := c.Args()[0]
	client := newEtcdMinionClientFromFlags(c)

	minions, err := client.MinionWithClassifierKey(classifierKey)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if len(minions) == 0 {
		return nil
	}

	report := make(map[string]int)
	for _, minion := range minions {
		classifier, err := client.MinionClassifier(minion, classifierKey)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		report[classifier.Value]++
	}

	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("CLASSIFIER", "VALUE", "MINION(S)")

	for classifierValue, minionCount := range report {
		table.AddRow(classifierKey, classifierValue, minionCount)
	}

	fmt.Println(table)

	return nil
}
