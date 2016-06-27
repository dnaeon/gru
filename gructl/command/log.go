package command

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/gosuri/uitable"
	"github.com/pborman/uuid"
	"github.com/urfave/cli"
)

// NewLogCommand creates a new sub-command for retrieving the
// log of previously executed tasks by minions
func NewLogCommand() cli.Command {
	cmd := cli.Command{
		Name:   "log",
		Usage:  "list minion task log",
		Action: execLogCommand,
	}

	return cmd
}

// Executes the "log" command
func execLogCommand(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return cli.NewExitError(errNoMinion.Error(), 64)
	}

	minion := uuid.Parse(c.Args()[0])
	if minion == nil {
		return cli.NewExitError(errInvalidUUID.Error(), 64)
	}

	klient := newEtcdMinionClientFromFlags(c)

	// Ignore errors about missing log directory
	log, err := klient.MinionTaskLog(minion)
	if err != nil {
		if eerr, ok := err.(client.Error); !ok || eerr.Code != client.ErrorCodeKeyNotFound {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	if len(log) == 0 {
		return nil
	}

	table := uitable.New()
	table.MaxColWidth = 40
	table.AddRow("TASK", "STATE", "RECEIVED", "PROCESSED")
	for _, id := range log {
		t, err := klient.MinionTaskResult(minion, id)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		table.AddRow(t.ID, t.State, time.Unix(t.TimeReceived, 0), time.Unix(t.TimeProcessed, 0))
	}

	fmt.Println(table)

	return nil
}
