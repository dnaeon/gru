package catalog

import (
	"fmt"
	"io"

	"github.com/dnaeon/gru/module"
	"github.com/dnaeon/gru/resource"
)

// Catalog type represents a collection of modules and resources
type Catalog struct {
	Modules   []*module.Module
	Resources []resource.Resource
}

// New creates a new empty catalog
func New() *Catalog {
	c := &Catalog{
		Modules:   make([]*module.Module, 0),
		Resources: make([]resource.Resource, 0),
	}

	return c
}

// Run processes the catalog
func (c *Catalog) Run(w io.Writer) error {
	fmt.Fprintf(w, "Loaded %d resources from %d modules\n", len(c.Resources), len(c.Modules))
	for _, r := range c.Resources {
		id := r.ResourceID()

		state, err := r.Evaluate()
		if err != nil {
			fmt.Fprintf(w, "%s: %s\n", id, err)
			continue
		}

		if !resource.StateIsValid(state.Want) || !resource.StateIsValid(state.Current) {
			fmt.Fprintf(w, "Invalid state(s) for resource %s: want %s, current %s\n", id, state.Want, state.Current)
			continue
		}

		// If resource is in the desired state, but out of date
		if state.Want == state.Current {
			if state.Update {
				fmt.Fprintf(w, "%s is out of date\n", r.ResourceID())
				if err := r.Update(w); err != nil {
					fmt.Fprintf(w, "%s error: %s\n", r.ResourceID(), err)
				}
			}
			continue
		}

		fmt.Fprintf(w, "%s is %s, should be %s\n", id, state.Current, state.Want)
		var action func(w io.Writer) error
		if state.Want == resource.StatePresent || state.Want == resource.StateRunning {
			if state.Current == resource.StateAbsent || state.Current == resource.StateStopped {
				action = r.Create
			}
		} else {
			if state.Current == resource.StatePresent || state.Current == resource.StateRunning {
				action = r.Delete
			}
		}

		// Perform the operation
		if err := action(w); err != nil {
			fmt.Fprintf(w, "%s error: %s", r.ResourceID(), err)
		}

		if state.Update {
			if err := r.Update(w); err != nil {
				fmt.Fprintf(w, "%s error %s\n", r.ResourceID(), err)
			}
		}
	}

	return nil
}

// Load creates a catalog from the provided module.
// The module is expected to be found in the module path
// provided by the path argument.
func Load(main, path string) (*Catalog, error) {
	c := New()

	// Discover and load the modules from the provided
	// module path, sort the import graph and
	// finally add the sorted modules to the catalog
	modules, err := module.DiscoverAndLoad(path)
	if err != nil {
		return c, err
	}

	modulesGraph, err := module.ImportGraph(main, path)
	if err != nil {
		return c, err
	}

	modulesSorted, err := modulesGraph.Sort()
	if err != nil {
		return c, err
	}

	for _, node := range modulesSorted {
		c.Modules = append(c.Modules, modules[node.Name])
	}

	// Build the dependency graph for the resources from the
	// loaded modules and sort them
	collection, err := module.ResourceCollection(c.Modules)
	if err != nil {
		return c, err
	}

	collectionGraph, err := collection.DependencyGraph()
	if err != nil {
		return c, err
	}

	collectionSorted, err := collectionGraph.Sort()
	if err != nil {
		return c, err
	}

	for _, node := range collectionSorted {
		c.Resources = append(c.Resources, collection[node.Name])
	}

	return c, nil
}
