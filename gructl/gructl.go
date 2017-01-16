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

package gructl

import (
	"os"
	"time"

	"github.com/dnaeon/gru/gructl/command"
	"github.com/dnaeon/gru/version"
	"github.com/urfave/cli"
)

// Main is the entry point of gructl
func Main() {
	app := cli.NewApp()
	app.Name = "gructl"
	app.Version = version.Version
	app.EnableBashCompletion = true
	app.Usage = "command line tool for managing minions"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "endpoint",
			Value:  "http://127.0.0.1:2379,http://localhost:4001",
			Usage:  "etcd cluster endpoints",
			EnvVar: "GRU_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "username",
			Value:  "",
			Usage:  "username to use for authentication",
			EnvVar: "GRU_USERNAME",
		},
		cli.StringFlag{
			Name:   "password",
			Value:  "",
			Usage:  "password to use for authentication",
			EnvVar: "GRU_PASSWORD",
		},
		cli.DurationFlag{
			Name:   "timeout",
			Value:  time.Second,
			Usage:  "connection timeout per request",
			EnvVar: "GRU_TIMEOUT",
		},
	}

	app.Commands = []cli.Command{
		command.NewApplyCommand(),
		command.NewListCommand(),
		command.NewInfoCommand(),
		command.NewServeCommand(),
		command.NewPushCommand(),
		command.NewClassifierCommand(),
		command.NewReportCommand(),
		command.NewQueueCommand(),
		command.NewLogCommand(),
		command.NewLastseenCommand(),
		command.NewResultCommand(),
		command.NewGraphCommand(),
	}

	app.Run(os.Args)
}
