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

	"github.com/dnaeon/gru/catalog"
	"github.com/dnaeon/gru/graph"
	"github.com/dnaeon/gru/resource"
	"github.com/urfave/cli"
	"github.com/yuin/gopher-lua"
)

// NewGraphCommand creates a new sub-command for
// generating the resource DAG graph
func NewGraphCommand() cli.Command {
	cmd := cli.Command{
		Name:   "graph",
		Usage:  "generate graph representation of resources",
		Action: execGraphCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "siterepo",
				Value:  "",
				Usage:  "path/url to the site repo",
				EnvVar: "GRU_SITEREPO",
			},
		},
	}

	return cmd
}

// Executes the "graph" command
func execGraphCommand(c *cli.Context) error {
	if len(c.Args()) < 1 {
		return cli.NewExitError(errNoModuleName.Error(), 64)
	}

	L := lua.NewState()
	defer L.Close()

	module := c.Args()[0]
	config := &catalog.Config{
		Module:   module,
		DryRun:   true,
		Logger:   log.New(os.Stdout, "", log.LstdFlags),
		SiteRepo: c.String("siterepo"),
		L:        L,
	}

	katalog := catalog.New(config)
	resource.LuaRegisterBuiltin(L)
	if err := L.DoFile(module); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	collection, err := resource.CreateCollection(katalog.Unsorted)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	g, err := collection.DependencyGraph()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	g.AsDot("resources", os.Stdout)
	g.Reversed().AsDot("reversed", os.Stdout)

	sorted, err := g.Sort()
	if err == graph.ErrCircularDependency {
		circular := graph.New()
		circular.AddNode(sorted...)
		circular.AsDot("circular", os.Stdout)
		return cli.NewExitError(graph.ErrCircularDependency.Error(), 1)
	}

	return nil
}
