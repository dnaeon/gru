package catalog

import (
	"fmt"
	"io"

	"github.com/dnaeon/gru/resource"
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

// Catalog type contains a collection of resources
type Catalog struct {
	// Unsorted contains the list of resources created by Lua
	unsorted []resource.Resource

	// Sorted contains the list of resources after a topological sort
	sorted []resource.Resource

	// Configuration settings
	config *Config
}

// Config type represents a set of settings to use when
// creating and processing the catalog
type Config struct {
	// Name of the Lua module to load and execute
	Module string

	// Do not take any actions, just report what would be done
	DryRun bool

	// Writer used to log events
	Writer io.Writer

	// Path to the site repo containing module and data files
	SiteRepo string

	// The Lua state
	L *lua.LState
}

// Run processes the catalog
func (c *Catalog) Run() error {
	// Use the same writer as the one used by the resources
	w := c.Config.ModuleConfig.ResourceConfig.Writer

	fmt.Fprintf(w, "Loaded %d resources from %d modules\n", len(c.Resources), len(c.Modules))
	for _, r := range c.Resources {
		id := r.ResourceID()

		state, err := r.Evaluate()
		if err != nil {
			fmt.Fprintf(w, "%s %s\n", id, err)
			continue
		}

		if c.Config.DryRun {
			continue
		}

		// TODO: Skip resources which have failed dependencies

		var resourceErr error
		switch {
		case state.Want == state.Current:
			// Resource is in the desired state
			break
		case state.Want == resource.StatePresent || state.Want == resource.StateRunning:
			// Resource is absent, should be present
			if state.Current == resource.StateAbsent || state.Current == resource.StateStopped {
				fmt.Fprintf(w, "%s is %s, should be %s\n", id, state.Current, state.Want)
				resourceErr = r.Create()
			}
		case state.Want == resource.StateAbsent || state.Want == resource.StateStopped:
			// Resource is present, should be absent
			if state.Current == resource.StatePresent || state.Current == resource.StateRunning {
				fmt.Fprintf(w, "%s is %s, should be %s\n", id, state.Current, state.Want)
				resourceErr = r.Delete()
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
			if err := r.Update(); err != nil {
				fmt.Fprintf(w, "%s %s\n", id, err)
			}
		}
	}

	return nil
}

// Load creates a new catalog from the provided configuration
func Load(config *Config) (*Catalog, error) {
	c := &Catalog{
		config:   config,
		sorted:   make([]resource.Resource, 0),
		unsorted: make([]resource.Resource, 0),
	}

	// Inject the configuration for resources
	resource.DefaultConfig = &resource.Config{
		Writer:   config.Writer,
		SiteRepo: config.SiteRepo,
	}

	// Register the resources and catalog in Lua
	resource.LuaRegisterBuiltin(config.L)
	config.L.SetGlobal("catalog", luar.New(config.L, c.unsorted))
	if err := L.DoFile(config.Module); err != nil {
		return c, err
	}

	// Perform a topological sort of the resources
	collection, err := resource.CreateCollection(c.unsorted)
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
		c.sorted = append(c.sorted, collection[node.Name])
	}

	return c, nil
}
