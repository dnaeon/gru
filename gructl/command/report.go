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

	"github.com/gosuri/uitable"
	"github.com/urfave/cli"
)

// NewReportCommand creates a new sub-command for
// generating reports based on minion classifiers
func NewReportCommand() cli.Command {
	cmd := cli.Command{
		Name:   "report",
		Usage:  "generate classifier report",
		Action: execReportCommand,
	}

	return cmd
}

// Executes the "report" command
func execReportCommand(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return cli.NewExitError(errNoClassifier.Error(), 64)
	}

	classifierKey := c.Args()[0]
	client := newEtcdMinionClientFromFlags(c)

	minions, err := client.MinionWithClassifierKey(classifierKey)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if len(minions) == 0 {
		return nil
	}

	report := make(map[string]int)
	for _, minion := range minions {
		classifier, err := client.MinionClassifier(minion, classifierKey)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		report[classifier.Value]++
	}

	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("CLASSIFIER", "VALUE", "MINION(S)")

	for classifierValue, minionCount := range report {
		table.AddRow(classifierKey, classifierValue, minionCount)
	}

	fmt.Println(table)

	return nil
}
