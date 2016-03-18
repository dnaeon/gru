package command

import (
	"fmt"
	"time"

	"github.com/codegangsta/cli"
	etcdclient "github.com/coreos/etcd/client"
	"github.com/gosuri/uitable"
	"github.com/pborman/uuid"
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
func execLogCommand(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayError(errNoMinion, 64)
	}

	minion := uuid.Parse(c.Args()[0])
	if minion == nil {
		displayError(errInvalidUUID, 64)
	}

	client := newEtcdMinionClientFromFlags(c)

	// Ignore errors about missing log directory
	log, err := client.MinionTaskLog(minion)
	if err != nil {
		if eerr, ok := err.(etcdclient.Error); !ok || eerr.Code != etcdclient.ErrorCodeKeyNotFound {
			displayError(err, 1)
		}
	}

	if len(log) == 0 {
		return
	}

	table := uitable.New()
	table.MaxColWidth = 40
	table.AddRow("TASK", "STATE", "RECEIVED", "PROCESSED")
	for _, taskID := range log {
		t, err := client.MinionTaskResult(minion, taskID)
		if err != nil {
			displayError(err, 1)
		}
		table.AddRow(t.TaskID, t.State, time.Unix(t.TimeReceived, 0), time.Unix(t.TimeProcessed, 0))
	}

	fmt.Println(table)
}
