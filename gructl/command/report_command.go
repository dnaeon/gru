package command

import (
	"fmt"

	"github.com/codegangsta/cli"
)

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
		displayError(errMissingClassifier, 64)
	}

	classifierKey := c.Args()[0]
	client := newEtcdMinionClientFromFlags(c)

	fmt.Printf("Generating report for classifier: %s\n", classifierKey)
	minions, err := client.MinionWithClassifierKey(classifierKey)
	if err != nil {
		displayError(err, 1)
	}

	fmt.Printf("Found %d minion(s) with the given classifier", len(minions))

	if len(minions) == 0 {
		return
	}

	report := make(map[string]int)
	for _, minion := range minions {
		classifier, err := client.MinionClassifier(minion, classifierKey)
		if err != nil {
			displayError(err, 1)
		}
		if _, ok := report[classifier.Value]; ok {
			report[classifier.Value] += 1
		} else {
			report[classifier.Value] = 1
		}
	}

	fmt.Println("\n")
	for k, v := range report {
		fmt.Printf("\t%s: %d minion(s)\n", k, v)
	}
}
