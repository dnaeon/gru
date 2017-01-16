// Copyright (c) 2015-2017 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//  1. Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer
//     in this position and unchanged.
//  2. Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in the
//     documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package command

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/gosuri/uitable"
	"github.com/pborman/uuid"
	"github.com/urfave/cli"
)

// NewInfoCommand creates a new sub-command for retrieving
// minion information
func NewInfoCommand() cli.Command {
	cmd := cli.Command{
		Name:   "info",
		Usage:  "get minion info",
		Action: execInfoCommand,
	}

	return cmd
}

// Executes the "info" command
func execInfoCommand(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return cli.NewExitError(errNoMinion.Error(), 64)
	}

	arg := c.Args()[0]
	minion := uuid.Parse(arg)
	if minion == nil {
		return cli.NewExitError(errInvalidUUID.Error(), 64)
	}

	klient := newEtcdMinionClientFromFlags(c)
	name, err := klient.MinionName(minion)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	lastseen, err := klient.MinionLastseen(minion)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	// Ignore errors about missing queue directory
	taskQueue, err := klient.MinionTaskQueue(minion)
	if err != nil {
		if eerr, ok := err.(client.Error); !ok || eerr.Code != client.ErrorCodeKeyNotFound {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	// Ignore errors about missing log directory
	taskLog, err := klient.MinionTaskLog(minion)
	if err != nil {
		if eerr, ok := err.(client.Error); !ok || eerr.Code != client.ErrorCodeKeyNotFound {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	// Ignore errors about missing classifier directory
	classifierKeys, err := klient.MinionClassifierKeys(minion)
	if err != nil {
		if eerr, ok := err.(client.Error); !ok || eerr.Code != client.ErrorCodeKeyNotFound {
			return cli.NewExitError(err.Error(), 1)
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

	return nil
}
