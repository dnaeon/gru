package command

import (
	"github.com/codegangsta/cli"
	"github.com/pborman/uuid"
)

func NewUnregisterCommand() cli.Command {
	cmd := cli.Command{
		Name:   "unregister",
		Usage:  "unregister a minion",
		Action: execUnregisterCommand,
	}

	return cmd
}

func execUnregisterCommand(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayError(errMissingMinion, 64)
	}

	for _, arg := range c.Args() {
		minionId := uuid.Parse(arg)
		if minionId == nil {
			displayError(errInvalidUUID, 64)
		}

		// TODO: Actually unregister the minion
		// TODO: Add Unregister() method to the client.Client interface?
	}
}
