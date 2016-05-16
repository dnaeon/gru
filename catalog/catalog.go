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
func (c *Catalog) Run(w io.Writer, opts *resource.Options) error {
	fmt.Fprintf(w, "Loaded %d resources from %d modules\n", len(c.Resources), len(c.Modules))
	for _, r := range c.Resources {
		id := r.ResourceID()

		state, err := r.Evaluate(w, opts)
		if err != nil {
			fmt.Fprintf(w, "%s %s\n", id, err)
			continue
		}

		if opts.DryRun {
			continue
		}

		var resourceErr error
		switch {
		case state.Want == state.Current:
			// Resource is in the desired state
			break
		case state.Want == resource.StatePresent || state.Want == resource.StateRunning:
			// Resource is absent, should be present
			if state.Current == resource.StateAbsent || state.Current == resource.StateStopped {
				fmt.Fprintf(w, "%s is %s, should be %s\n", id, state.Current, state.Want)
				resourceErr = r.Create(w, opts)
			}
		case state.Want == resource.StateAbsent || state.Want == resource.StateStopped:
			// Resource is present, should be absent
			if state.Current == resource.StatePresent || state.Current == resource.StateRunning {
				fmt.Fprintf(w, "%s is %s, should be %s\n", id, state.Current, state.Want)
				resourceErr = r.Delete(w, opts)
			}
		default:
			fmt.Fprintf(w, "%s unknown state(s): want %s, current %s\n", id, state.Want, state.Current)
			continue
		}

		if resourceErr != nil {
			fmt.Fprintf(w, "%s %s\n", id, resourceErr)
		}

		// Update resource if needed
		if state.Update {
			fmt.Fprintf(w, "%s resource is out of date, will be updated\n", id)
			if err := r.Update(w, opts); err != nil {
				fmt.Fprintf(w, "%s %s\n", id, err)
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
