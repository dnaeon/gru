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

