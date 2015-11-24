package command

import (
	"fmt"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/cli"
	"github.com/gosuri/uitable"
)

func NewClassifierCommand() cli.Command {
	cmd := cli.Command{
		Name:   "classifier",
		Usage:  "list minion classifiers",
		Action: execClassifierCommand,
	}

	return cmd
}

// Executes the "classifier" command
func execClassifierCommand(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayError(errMissingMinion, 64)
	}

	arg := c.Args()[0]
	minion := uuid.Parse(arg)
	if minion == nil {
		displayError(errInvalidUUID, 64)
	}

	client := newEtcdMinionClientFromFlags(c)
	classifierKeys, err := client.MinionClassifierKeys(minion)
	if err != nil {
		displayError(err, 1)
	}

	if len(classifierKeys) == 0 {
		return
	}

	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("KEY", "VALUE")
	for _, key := range classifierKeys {
		classifier, err := client.MinionClassifier(minion, key)
		if err != nil {
			displayError(err, 1)
		}

		table.AddRow(classifier.Key, classifier.Value)
	}

	fmt.Println(table)
}
