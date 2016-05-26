package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/dnaeon/gru/graph"
	"github.com/dnaeon/gru/module"
	"github.com/dnaeon/gru/resource"
)

// NewGraphCommand creates a new sub-command for
// generating the resource DAG graph
func NewGraphCommand() cli.Command {
	cmd := cli.Command{
		Name:   "graph",
		Usage:  "create DOT representation for modules and resources",
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
func execGraphCommand(c *cli.Context) {
	if len(c.Args()) < 1 {
		displayError(errNoModuleName, 64)
	}

	main := c.Args()[0]
	config := &module.Config{
		Path:           filepath.Join(c.String("siterepo"), "modules"),
		ResourceConfig: &resource.Config{},
	}

	// Discover and load all modules
	discovered, err := module.DiscoverAndLoad(config)
	if err != nil {
		displayError(err, 1)
	}

	// Generate DOT representation of the module imports
	modulesGraph, err := module.ImportGraph(main, config.Path)
	if err != nil {
		displayError(err, 1)
	}
	modulesGraph.AsDot(fmt.Sprintf("%s_imports", main), os.Stdout)

	// Attempt a topological sort of the module imports graph.
	// In case of circular dependencies in the graph
	// generate a DOT for the remaining nodes, which would
	// give us the nodes causing the circular dependencies
	modulesSorted, err := modulesGraph.Sort()
	if err == graph.ErrCircularDependency {
		circular := graph.New()
		circular.AddNode(modulesSorted...)
		circular.AsDot(fmt.Sprintf("%s_imports_circular", main), os.Stdout)
		displayError(graph.ErrCircularDependency, 1)
	}

	// Get the sorted modules and create a DOT for their resources
	var modules []*module.Module
	for _, node := range modulesSorted {
		modules = append(modules, discovered[node.Name])
	}

	collection, err := module.ResourceCollection(modules)
	if err != nil {
		displayError(err, 1)
	}

	collectionGraph, err := collection.DependencyGraph()
	if err != nil {
		displayError(err, 1)
	}

	// Create resource nodes in their respective module
	fmt.Fprintf(os.Stdout, "digraph resources {\n")
	fmt.Fprintf(os.Stdout, "\tlabel = \"resources\";\n")
	fmt.Fprintf(os.Stdout, "\tnodesep=1.0;\n")
	fmt.Fprintf(os.Stdout, "\tnode [shape=box];\n")
	fmt.Fprintf(os.Stdout, "\tedge [style=filled];\n")
	for _, m := range modules {
		fmt.Fprintf(os.Stdout, fmt.Sprintf("\tsubgraph cluster_%s {\n", m.Name))
		fmt.Fprintf(os.Stdout, fmt.Sprintf("\t\tlabel = %q;\n", m.Name))
		fmt.Fprintf(os.Stdout, "\t\tcolor = black;\n")
		for _, r := range m.Resources {
			fmt.Fprintf(os.Stdout, fmt.Sprintf("\t\t%q;\n", r.ResourceID()))
		}
		fmt.Fprintf(os.Stdout, "\t}\n")
	}

	// Finally, connect the nodes
	for _, node := range collectionGraph.Nodes {
		if len(node.Edges) == 0 {
			continue
		}

		var edges []string
		for _, edge := range node.Edges {
			edges = append(edges, fmt.Sprintf("%q", edge.Name))
		}

		fmt.Fprintf(os.Stdout, fmt.Sprintf("\t%q -> {%s};\n", node.Name, strings.Join(edges, " ")))
	}
	fmt.Fprintf(os.Stdout, "}\n")

	// Attempt a topological sort of the resource dependencies.
	// In case of circular dependencies in the graph
	// generate a DOT for the remaining nodes, which would
	// give us the resources causing the circular dependencies
	collectionSorted, err := collectionGraph.Sort()
	if err == graph.ErrCircularDependency {
		circular := graph.New()
		circular.AddNode(collectionSorted...)
		circular.AsDot("resources_circular", os.Stdout)
		displayError(graph.ErrCircularDependency, 1)
	}
}
