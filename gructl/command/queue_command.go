package command

import (
	"fmt"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/cli"
	etcdclient "github.com/coreos/etcd/client"
)

func NewQueueCommand() cli.Command {
	cmd := cli.Command{
		Name:   "queue",
		Usage:  "list minion task queue",
		Action: execQueueCommand,
	}

	return cmd
}

// Executes the "queue" command
func execQueueCommand(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayError(errMissingMinion, 64)
	}

	minion := uuid.Parse(c.Args()[0])
	if minion == nil {
		displayError(errInvalidUUID, 64)
	}

	client := newEtcdMinionClientFromFlags(c)

	// Ignore errors about missing queue directory
	queue, err := client.MinionTaskQueue(minion)
	if err != nil {
		if eerr, ok := err.(etcdclient.Error); !ok || eerr.Code != etcdclient.ErrorCodeKeyNotFound {
			displayError(err, 1)
		}
	}

	for _, t := range queue {
		fmt.Println(t.TaskID)
	}
}
