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
	"log"
	"os"
	"runtime"

	"github.com/dnaeon/gru/catalog"
	"github.com/urfave/cli"
	"github.com/yuin/gopher-lua"
)

// NewApplyCommand creates a new sub-command for
// applying configurations on the local system
func NewApplyCommand() cli.Command {
	cmd := cli.Command{
		Name:   "apply",
		Usage:  "apply configuration on the local system",
		Action: execApplyCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "siterepo",
				Value:  "",
				Usage:  "path/url to the site repo",
				EnvVar: "GRU_SITEREPO",
			},
			cli.BoolFlag{
				Name:  "dry-run",
				Usage: "just report what would be done, instead of doing it",
			},
			cli.IntFlag{
				Name:  "concurrency",
				Usage: "number of goroutines used for concurrent processing",
				Value: runtime.NumCPU(),
			},
		},
	}

	return cmd
}

// Executes the "apply" command
func execApplyCommand(c *cli.Context) error {
	if len(c.Args()) < 1 {
		return cli.NewExitError(errNoModuleName.Error(), 64)
	}

	concurrency := c.Int("concurrency")
	if concurrency < 0 {
		concurrency = runtime.NumCPU()
	}

	L := lua.NewState()
	defer L.Close()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	config := &catalog.Config{
		Module:      c.Args()[0],
		DryRun:      c.Bool("dry-run"),
		Logger:      logger,
		SiteRepo:    c.String("siterepo"),
		L:           L,
		Concurrency: concurrency,
	}

	katalog := catalog.New(config)
	if err := katalog.Load(); err != nil {
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	status := katalog.Run()
	status.Summary(logger)

	return nil
}
