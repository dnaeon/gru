package command

import (
	"fmt"
	"time"

	"github.com/pborman/uuid"
	"github.com/codegangsta/cli"
	etcdclient "github.com/coreos/etcd/client"
	"github.com/gosuri/uitable"
)

func NewInfoCommand() cli.Command {
	cmd := cli.Command{
		Name:   "info",
		Usage:  "get minion info",
		Action: execInfoCommand,
	}

	return cmd
}

// Executes the "info" command
func execInfoCommand(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayError(errMissingMinion, 64)
	}

	arg := c.Args()[0]
	minion := uuid.Parse(arg)
	if minion == nil {
		displayError(errInvalidUUID, 64)
	}

	client := newEtcdMinionClientFromFlags(c)
	name, err := client.MinionName(minion)
	if err != nil {
		displayError(err, 1)
	}

	lastseen, err := client.MinionLastseen(minion)
	if err != nil {
		displayError(err, 1)
	}

	// Ignore errors about missing queue directory
	taskQueue, err := client.MinionTaskQueue(minion)
	if err != nil {
		if eerr, ok := err.(etcdclient.Error); !ok || eerr.Code != etcdclient.ErrorCodeKeyNotFound {
			displayError(err, 1)
		}
	}

	// Ignore errors about missing log directory
	taskLog, err := client.MinionTaskLog(minion)
	if err != nil {
		if eerr, ok := err.(etcdclient.Error); !ok || eerr.Code != etcdclient.ErrorCodeKeyNotFound {
			displayError(err, 1)
		}
	}

	// Ignore errors about missing classifier directory
	classifierKeys, err := client.MinionClassifierKeys(minion)
	if err != nil {
		if eerr, ok := err.(etcdclient.Error); !ok || eerr.Code != etcdclient.ErrorCodeKeyNotFound {
			displayError(err, 1)
		}
	}

	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("Minion:", minion)
	table.AddRow("Name:", name)
	table.AddRow("Lastseen:", time.Unix(lastseen, 0))
	table.AddRow("Queue:", len(taskQueue))
	table.AddRow("Log:", len(taskLog))
	table.AddRow("Classifiers:", len(classifierKeys))

	fmt.Println(table)
}
