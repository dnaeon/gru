package command

import (
	"fmt"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/cli"
	etcdclient "github.com/coreos/etcd/client"
)

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
		displayError(errMissingMinion, 64)
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

	for _, t := range log {
		fmt.Println(t)
	}
}
