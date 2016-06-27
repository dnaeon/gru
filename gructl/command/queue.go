package command

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/gosuri/uitable"
	"github.com/pborman/uuid"
	"github.com/urfave/cli"
)

// NewQueueCommand creates a new sub-command for retrieving the
// currently pending tasks for minions
func NewQueueCommand() cli.Command {
	cmd := cli.Command{
		Name:   "queue",
		Usage:  "list minion task queue",
		Action: execQueueCommand,
	}

	return cmd
}

// Executes the "queue" command
func execQueueCommand(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return cli.NewExitError(errNoMinion.Error(), 64)
	}

	minion := uuid.Parse(c.Args()[0])
	if minion == nil {
		return cli.NewExitError(errInvalidUUID.Error(), 64)
	}

	klient := newEtcdMinionClientFromFlags(c)

	// Ignore errors about missing queue directory
	queue, err := klient.MinionTaskQueue(minion)
	if err != nil {
		if eerr, ok := err.(client.Error); !ok || eerr.Code != client.ErrorCodeKeyNotFound {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	if len(queue) == 0 {
		return nil
	}

	table := uitable.New()
	table.MaxColWidth = 40
	table.AddRow("TASK", "STATE", "RECEIVED")
	for _, t := range queue {
		table.AddRow(t.ID, t.State, time.Unix(t.TimeReceived, 0))
	}

	fmt.Println(table)

	return nil
}
