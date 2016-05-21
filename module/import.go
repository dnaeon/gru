package module

import (
	"fmt"
	"io"

	"github.com/dnaeon/gru/graph"
)

// ImportGraph creates a DAG graph of the
// module imports for a given module.
// The resulting DAG graph can be used to determine the
// proper ordering of modules and also to detect whether
// we have circular imports in our modules.
func ImportGraph(main, path string) (*graph.Graph, error) {
	g := graph.New()

	config := &Config{
		Path: path,
	}
	modules, err := DiscoverAndLoad(config)
	if err != nil {
		return g, err
	}

	if _, ok := modules[main]; !ok {
		return g, fmt.Errorf("Module %s not found in module path", main)
	}

	// A map containing the modules as graph nodes
	// The graph can be used to determine if we have
	// circular module imports and also to provide the
	// proper ordering of loading modules after a
	// topological sort of the graph nodes
	nodes := make(map[string]*graph.Node)
	for n := range modules {
		node := graph.NewNode(n)
		nodes[n] = node
	}

	// Recursively find all imports that the main module has
	var buildImportGraphFunc func(m *Module) error
	buildImportGraphFunc = func(m *Module) error {
		// Add the node to the graph if it is not present already
		if _, ok := g.GetNode(m.Name); !ok {
			g.AddNode(nodes[m.Name])
		} else {
			return nil
		}

		// Build the import graph for each imported module
		for _, mi := range m.Imports {
			if _, ok := modules[mi.Name]; !ok {
				return fmt.Errorf("Module %s imports %s, which is not in the module path", m.Name, mi.Name)
			}

			// Build the dependencies of imported modules as well
			buildImportGraphFunc(modules[mi.Name])

			// Finally connect the nodes in the graph
			g.AddEdge(nodes[m.Name], nodes[mi.Name])
		}

		return nil
	}

	if err := buildImportGraphFunc(modules[main]); err != nil {
		return g, err
	}

	return g, nil
}

// ImportGraphAsDot creates a DOT representation of the module imports
func ImportGraphAsDot(main, path string, w io.Writer) error {
	g, err := ImportGraph(main, path)
	if err != nil {
		return err
	}

	g.AsDot("modules", w)

	// Try a topological sort of the graph
	// In case of circular dependencies in the graph
	// generate a DOT for the remaining nodes in the graph,
	// which would give us the modules causing circular dependencies
	if nodes, err := g.Sort(); err == graph.ErrCircularDependency {
		circular := graph.New()
		circular.AddNode(nodes...)
		circular.AsDot("modules_circular", w)
	}

	return nil
}
