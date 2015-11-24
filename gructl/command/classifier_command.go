package command

import (
	"errors"
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

	table := uitable.New()
	table.MaxColWidth = 60
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
