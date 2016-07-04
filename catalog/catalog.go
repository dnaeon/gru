package catalog

import (
	"fmt"
	"log"

	"github.com/dnaeon/gru/resource"
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

// Catalog type contains a collection of resources
type Catalog struct {
	// Unsorted contains the list of resources created by Lua
	Unsorted []resource.Resource `luar:"-"`

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
	Logger *log.Logger

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
		Unsorted: make([]resource.Resource, 0),
	}

	// Inject the configuration for resources
	resource.DefaultConfig = &resource.Config{
		Logger:   config.Logger,
		SiteRepo: config.SiteRepo,
	}

	// Register the catalog type in Lua and also register
	// metamethods for the catalog, so that we can use
	// the catalog in a more Lua-friendly way
	mt := luar.MT(config.L, c)
	mt.RawSetString("__len", luar.New(config.L, (*Catalog).Len))
	config.L.SetGlobal("catalog", luar.New(config.L, c))

	return c
}

// Add adds a resource to the catalog.
// This method is called from Lua when adding new resources
func (c *Catalog) Add(r ...resource.Resource) {
	c.Unsorted = append(c.Unsorted, r...)
}

// Len returns the number of unsorted resources in catalog
func (c *Catalog) Len() int {
	return len(c.Unsorted)
}

// Load loads resources into the catalog
func (c *Catalog) Load() error {
	// Register the resource providers and catalog in Lua
	resource.LuaRegisterBuiltin(c.config.L)
	if err := c.config.L.DoFile(c.config.Module); err != nil {
		return err
	}

	// Perform a topological sort of the resources
	collection, err := resource.CreateCollection(c.Unsorted)
	if err != nil {
		return err
	}

	collectionGraph, err := collection.DependencyGraph()
	if err != nil {
		return err
	}

	collectionSorted, err := collectionGraph.Sort()
	if err != nil {
		return err
	}

	for _, node := range collectionSorted {
		c.sorted = append(c.sorted, collection[node.Name])
	}

	return nil
}

// Run processes the catalog
func (c *Catalog) Run() error {
	c.config.Logger.Printf("Loaded %d resources\n", len(c.sorted))
	for _, r := range c.sorted {
		if err := c.processResource(r); err != nil {
			c.config.Logger.Printf("%s %s\n", r.ID(), err)
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
			c.config.Logger.Printf("%s is %s, should be %s\n", id, state.Current, state.Want)
			if err := r.Create(); err != nil {
				return err
			}
		}
	case state.Want == resource.StateAbsent || state.Want == resource.StateStopped:
		// Resource is present, should be absent
		if state.Current == resource.StatePresent || state.Current == resource.StateRunning {
			c.config.Logger.Printf("%s is %s, should be %s\n", id, state.Current, state.Want)
			if err := r.Delete(); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown state %s", state.Want)
	}

	// Update resource if needed
	if state.Update {
		c.config.Logger.Printf("%s resource is out of date\n", id)
		if err := r.Update(); err != nil {
			return err
		}
	}

	return nil
}
