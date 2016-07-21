package catalog

import (
	"fmt"
	"log"

	"github.com/dnaeon/gru/resource"
	"github.com/dnaeon/gru/utils"
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

// Catalog type contains a collection of resources
type Catalog struct {
	// Unsorted contains the list of resources created by Lua
	Unsorted []resource.Resource `luar:"-"`

	// Sorted contains the list of resources after a topological sort
	sorted []resource.Resource `luar:"-"`

	// Result contains the results of resource processing and any
	// errors that might have occurred during processing.
	// Keys of the map are the resource ids and their
	// values are the errors returned from resources.
	result map[string]error `luar:"-"`

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
		result:   make(map[string]error),
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
func (c *Catalog) Add(resources ...resource.Resource) {
	for _, r := range resources {
		if r != nil {
			c.Unsorted = append(c.Unsorted, r)
		}
	}
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

	// TODO: Find candidates for concurrent processing

	for _, node := range collectionSorted {
		c.sorted = append(c.sorted, collection[node.Name])
	}

	return nil
}

// Run processes the resources from catalog
func (c *Catalog) Run() error {
	c.config.Logger.Printf("Loaded %d resources\n", len(c.sorted))
	for _, r := range c.sorted {
		id := r.ID()

		// Skip resource
		if err := c.shouldBeSkipped(r); err != nil {
			c.result[id] = err
			c.config.Logger.Printf("%s skipping resource: %s\n", id, err)
			continue
		}

		// Process resource and save result
		if c.result[id] = c.processResource(r); c.result[id] != nil {
			c.config.Logger.Printf("%s %s\n", id, c.result[id])
		}
	}

	return nil
}

// shouldBeSkipped checks if a resource has failed
// dependencies and should be skipped from further processing.
func (c *Catalog) shouldBeSkipped(r resource.Resource) error {
	for _, dep := range r.Dependencies() {
		if c.result[dep] != nil {
			return fmt.Errorf("failed dependency for %s", dep)
		}
	}

	return nil
}

// processResource processes a single resource
func (c *Catalog) processResource(r resource.Resource) error {
	if err := r.Validate(); err != nil {
		return err
	}

	state, err := r.Evaluate()
	if err != nil {
		return err
	}

	if c.config.DryRun {
		return nil
	}

	// Current and wanted states for the resource
	want := utils.NewString(state.Want)
	current := utils.NewString(state.Current)

	// The list of present and absent states for the resource
	present := utils.NewList(r.GetPresentStates()...)
	absent := utils.NewList(r.GetAbsentStates()...)

	var action func() error
	switch {
	case want.IsInList(present) && current.IsInList(absent):
		action = r.Create
	case want.IsInList(absent) && current.IsInList(present):
		action = r.Delete
	case state.Outdated:
		action = r.Update
	}

	if action != nil {
		return action()
	}

	return nil
}
