package command

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/gosuri/uitable"
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
func execReportCommand(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayError(errNoClassifier, 64)
	}

	classifierKey := c.Args()[0]
	client := newEtcdMinionClientFromFlags(c)

	minions, err := client.MinionWithClassifierKey(classifierKey)
	if err != nil {
		displayError(err, 1)
	}

	if len(minions) == 0 {
		return
	}

	report := make(map[string]int)
	for _, minion := range minions {
		classifier, err := client.MinionClassifier(minion, classifierKey)
		if err != nil {
			displayError(err, 1)
		}
		report[classifier.Value] += 1
	}

	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("CLASSIFIER", "VALUE", "MINION(S)")

	for classifierValue, minionCount := range report {
		table.AddRow(classifierKey, classifierValue, minionCount)
	}

	fmt.Println(table)
}
