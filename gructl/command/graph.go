package command

import (
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
		Usage:  "create DOT representation of resources",
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
		Logger:   resource.DefaultLogger,
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

	collectionGraph, err := collection.DependencyGraph()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	collectionGraph.AsDot("resources", os.Stdout)

	collectionSorted, err := collectionGraph.Sort()
	if err == graph.ErrCircularDependency {
		circular := graph.New()
		circular.AddNode(collectionSorted...)
		circular.AsDot("resources_circular", os.Stdout)
		return cli.NewExitError(graph.ErrCircularDependency.Error(), 1)
	}

	return nil
}
