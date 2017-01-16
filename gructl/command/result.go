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

	"github.com/gosuri/uitable"
	"github.com/pborman/uuid"
	"github.com/urfave/cli"
)

// NewResultCommand creates a new sub-command for retrieving
// results of previously executed tasks by minions
func NewResultCommand() cli.Command {
	cmd := cli.Command{
		Name:   "result",
		Usage:  "get task results",
		Action: execResultCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "minion",
				Usage: "get task result for given minion only",
			},
			cli.BoolFlag{
				Name:  "details",
				Usage: "provide more details about the tasks",
			},
		},
	}

	return cmd
}

// Executes the "result" command
func execResultCommand(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return cli.NewExitError(errNoTask.Error(), 64)
	}

	arg := c.Args()[0]
	id := uuid.Parse(arg)
	if id == nil {
		return cli.NewExitError(errInvalidUUID.Error(), 64)
	}

	client := newEtcdMinionClientFromFlags(c)

	// If --minion flag was specified parse the
	// minion uuid and get the task result only
	// from the specified minion, otherwise find
	// all minions which contain the given
	// task and get their results
	var minionWithTask []uuid.UUID

	mFlag := c.String("minion")
	if mFlag == "" {
		// No minion was specified, get all minions
		// with the given task uuid
		m, err := client.MinionWithTaskResult(id)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		minionWithTask = m
	} else {
		// Minion was specified, get task result
		// from the given minion only
		minion := uuid.Parse(mFlag)
		if minion == nil {
			return cli.NewExitError(errInvalidUUID.Error(), 64)
		}
		minionWithTask = append(minionWithTask, minion)
	}

	if len(minionWithTask) == 0 {
		return nil
	}

	// Create table for the task results
	// If the --details flag is specified, then
	// create a table that holds all details about the
	// tasks, otherwise use a simple summary table
	table := uitable.New()
	if c.Bool("details") {
		table.MaxColWidth = 80
		table.Wrap = true
	} else {
		table.MaxColWidth = 40
		table.AddRow("MINION", "RESULT", "STATE")
	}

	for _, minionID := range minionWithTask {
		t, err := client.MinionTaskResult(minionID, id)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if c.Bool("details") {
			table.AddRow("Minion:", minionID)
			table.AddRow("Task ID:", t.ID)
			table.AddRow("State:", t.State)
			table.AddRow("Received:", time.Unix(t.TimeReceived, 0))
			table.AddRow("Processed:", time.Unix(t.TimeProcessed, 0))
			table.AddRow("Result:", t.Result)
		} else {
			table.AddRow(minionID, t.Result, t.State)
		}
	}

	fmt.Println(table)

	return nil
}
