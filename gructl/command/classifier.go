package command

import (
	"fmt"

	"github.com/coreos/etcd/client"
	"github.com/gosuri/uitable"
	"github.com/pborman/uuid"
	"github.com/urfave/cli"
)

// NewClassifierCommand creates a new sub-command for retrieving
// minion classifiers
func NewClassifierCommand() cli.Command {
	cmd := cli.Command{
		Name:   "classifier",
		Usage:  "list minion classifiers",
		Action: execClassifierCommand,
	}

	return cmd
}

// Executes the "classifier" command
func execClassifierCommand(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return cli.NewExitError(errNoMinion.Error(), 64)
	}

	arg := c.Args()[0]
	minion := uuid.Parse(arg)
	if minion == nil {
		return cli.NewExitError(errInvalidUUID.Error(), 64)
	}

	klient := newEtcdMinionClientFromFlags(c)

	// Ignore errors about missing classifier directory
	classifierKeys, err := klient.MinionClassifierKeys(minion)
	if err != nil {
		if eerr, ok := err.(client.Error); !ok || eerr.Code != client.ErrorCodeKeyNotFound {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	if len(classifierKeys) == 0 {
		return nil
	}

	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("KEY", "VALUE")
	for _, key := range classifierKeys {
		classifier, err := klient.MinionClassifier(minion, key)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		table.AddRow(classifier.Key, classifier.Value)
	}

	fmt.Println(table)

	return nil
}
