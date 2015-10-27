package command

import (
	"fmt"
	"errors"

	"code.google.com/p/go-uuid/uuid"

	"github.com/codegangsta/cli"
)

func NewClassifierCommands() cli.Command {
	cmd := cli.Command{
		Name: "classifier",
		Usage: "manage classifiers",
		Subcommands: []cli.Command{
			{
				Name: "list",
				Usage: "list classifiers of a minion",
				Action: classifierListCommand,
			},
			{
				Name: "report",
				Usage: "generate a classifier report",
				Action: classifierReportCommand,
			},
		},
	}

	return cmd
}

// Executes the "classifier list" command
func classifierListCommand(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayError(errors.New("Must provide a minion uuid"), 64)
	}

	arg := c.Args()[0]
	minion := uuid.Parse(arg)
	if minion == nil {
		displayError(errors.New("Bad minion uuid given"), 64)
	}

	client := newEtcdMinionClientFromFlags(c)
	classifierKeys, err := client.MinionClassifierKeys(minion)
	if err != nil {
		displayError(err, 1)
	}

	for _, key := range classifierKeys {
		classifier, err := client.MinionClassifier(minion, key)
		if err != nil {
			displayError(err, 1)
		}

		fmt.Printf("%s -> %s\n", classifier.Key, classifier.Value)
	}
}

// Executes the "classifier report" command
func classifierReportCommand(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayError(errors.New("Must provide a classifier key"), 64)
	}

	classifierKey := c.Args()[0]
	client := newEtcdMinionClientFromFlags(c)

	fmt.Printf("Generating report for classifier: %s\n", classifierKey)
	minions, err := client.MinionWithClassifier(classifierKey)
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
