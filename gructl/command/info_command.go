package command

import (
	"github.com/codegangsta/cli"
)

func NewInfoCommand() cli.Command {
	cmd := cli.Command{
		Name: "info",
		Usage: "get minion info",
		Action: runInfoCommand,
	}

	return cmd
}

// Executes the "info" command
func runInfoCommand(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayError(errors.New("Must provide a minion uuid"), 64)
	}

	arg := c.Args()[0]
	minion := uuid.Parse(arg)
	if minion == nil {
		displayError(errors.New("Bad minion uuid given"), 64)
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

	fmt.Printf("%-15s: %s\n", "Minion", minion)
	fmt.Printf("%-15s: %s\n", "Name", name)
	fmt.Printf("%-15s: %s\n", "Lastseen", time.Unix(lastseen, 0))
	fmt.Printf("%-15s: %d task(s)\n", "Queue", len(taskQueue))
	fmt.Printf("%-15s: %d task(s)\n", "Processed", len(taskLog))
	fmt.Printf("%-15s: %d key(s)\n", "Classifier", len(classifierKeys))
}
