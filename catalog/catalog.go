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
	unsorted []resource.Resource `luar:"-"`

	// Sorted contains the list of resources after a topological sort
	sorted []resource.Resource `luar:"-"`

	// Configuration settings
	config *Config `luar:"-"`
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

// New creates a new empty catalog with the provided configuration
func New(config *Config) *Catalog {
	c := &Catalog{
		config:   config,
		sorted:   make([]resource.Resource, 0),
		unsorted: make([]resource.Resource, 0),
	}

	return c
}

// Run processes the catalog
func (c *Catalog) Run() error {
	fmt.Fprintf(c.config.Writer, "Loaded %d resources\n", len(c.sorted))
	for _, r := range c.sorted {
		if err := c.processResource(r); err != nil {
			fmt.Fprintf(c.config.Writer, "%s %s\n", r.ID(), err)
		}
	}

	return nil
}

// processResource processes a single resource
func (c *Catalog) processResource(r resource.Resource) error {
	id := r.ID()
	state, err := r.Evaluate()
	if err != nil {
		return err
	}

	if c.config.DryRun {
		return nil
	}

	// TODO: Skip resources which have failed dependencies

	switch {
	case state.Want == state.Current:
		// Resource is in the desired state
		break
	case state.Want == resource.StatePresent || state.Want == resource.StateRunning:
		// Resource is absent, should be present
		if state.Current == resource.StateAbsent || state.Current == resource.StateStopped {
			fmt.Fprintf(c.config.Writer, "%s is %s, should be %s\n", id, state.Current, state.Want)
			if err := r.Create(); err != nil {
				return err
			}
		}
	case state.Want == resource.StateAbsent || state.Want == resource.StateStopped:
		// Resource is present, should be absent
		if state.Current == resource.StatePresent || state.Current == resource.StateRunning {
			fmt.Fprintf(c.config.Writer, "%s is %s, should be %s\n", id, state.Current, state.Want)
			if err := r.Delete(); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown state %s", state.Want)
	}

	// Update resource if needed
	if state.Update {
		fmt.Fprintf(c.config.Writer, "%s resource is out of date\n", id)
		if err := r.Update(); err != nil {
			return err
		}
	}

	return nil
}

// Load creates a new catalog from the provided configuration
func Load(config *Config) (*Catalog, error) {
	c := New()

	// Inject the configuration for resources
	resource.DefaultConfig = &resource.Config{
		Writer:   config.Writer,
		SiteRepo: config.SiteRepo,
	}

	// Register the resources and catalog in Lua
	resource.LuaRegisterBuiltin(config.L)
	config.L.SetGlobal("catalog", luar.New(config.L, c.unsorted))
	if err := config.L.DoFile(config.Module); err != nil {
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
