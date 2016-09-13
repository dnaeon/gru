package catalog

import (
	"fmt"
	"log"
	"sync"

	"github.com/dnaeon/gru/graph"
	"github.com/dnaeon/gru/resource"
	"github.com/dnaeon/gru/utils"
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

// Catalog type contains a collection of resources
type Catalog struct {
	// Unsorted contains the list of resources created by Lua
	Unsorted []resource.Resource `luar:"-"`

	// Collection contains the unsorted resources as a collection
	collection resource.Collection `luar:"-"`

	// Sorted contains the resources after a topological sort.
	sorted []*graph.Node `luar:"-"`

	// Reversed contains the resource dependency graph in reverse
	// order. It is used for finding the reverse dependencies of
	// resources.
	reversed *graph.Graph `luar:"-"`

	// Status contains status information about resources
	status *status `luar:"-"`

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

	// Number of goroutines to use for concurrent processing
	Concurrency int
}

// status type contains status information about processed resources
type status struct {
	sync.RWMutex

	// Items contain the result of resource processing and any
	// errors that might have occurred during processing.
	// Keys of the map are the resource ids and their
	// values are the errors returned by resources.
	items map[string]error
}

// set sets the status for a resource
func (s *status) set(id string, err error) {
	s.Lock()
	defer s.Unlock()
	s.items[id] = err
}

// get retrieves the status of a resource
func (s *status) get(id string) (error, bool) {
	s.Lock()
	defer s.Unlock()
	err, ok := s.items[id]

	return err, ok
}

// isSynced returns a boolean indicating whether a
// resource is up to date
func (s *status) isSynced(id string) bool {
	s.Lock()
	defer s.Unlock()

	return s.items[id] == resource.ErrInSync
}

// hasChanged returns a boolean indicating whether a
// resource state has changed after processing
func (s *state) hasChanged(id string) bool {
	s.Lock()
	defer s.Unlock()

	return s.items[id] == nil
}

// hasFailed returns a boolean indicating whether a
// resource has failed during processing
func (s *status) hasFailed(id string) bool {
	s.Lock()
	defer s.Unlock()

	return s.items[id] != nil && s.items[id] != resource.ErrInSync
}

// New creates a new empty catalog with the provided configuration
func New(config *Config) *Catalog {
	c := &Catalog{
		config:     config,
		collection: make(resource.Collection),
		sorted:     make([]*graph.Node, 0),
		reversed:   graph.New(),
		status: &status{
			items: make(map[string]error),
		},
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
	mt.RawSetString("__len", luar.New(config.L, (*Catalog).luaLen))
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

	reversed := collectionGraph.Reversed()

	sorted, err := collectionGraph.Sort()
	if err != nil {
		return err
	}

	// Set catalog fields
	c.collection = collection
	c.sorted = sorted
	c.reversed = reversed

	c.config.Logger.Printf("Loaded %d resources\n", len(c.sorted))

	return nil
}

// Run processes the resources from catalog
func (c *Catalog) Run() error {
	// process executes a single resource
	process := func(r resource.Resource) {
		id := r.ID()
		err := c.execute(r)
		c.status.set(id, err)
		if c.status.hasFailed(id) {
			c.config.Logger.Printf("%s %s\n", id, err)
		}
	}

	// Start goroutines for concurrent processing
	var wg sync.WaitGroup
	ch := make(chan resource.Resource, 1024)
	c.config.Logger.Printf("Starting %d goroutines for concurrent processing\n", c.config.Concurrency)
	for i := 0; i < c.config.Concurrency; i++ {
		wg.Add(1)
		worker := func() {
			defer wg.Done()
			for r := range ch {
				c.config.Logger.Printf("%s is concurrent", r.ID())
				process(r)
			}
		}
		go worker()
	}

	// Process the resources
	for _, node := range c.sorted {
		r := c.collection[node.Name]
		switch {
		// Resource is concurrent and is an isolated node
		case r.IsConcurrent() && len(r.Dependencies()) == 0 && len(c.reversed.Nodes[r.ID()].Edges) == 0:
			ch <- r
			continue
		// Resource is concurrent and has no reverse dependencies
		case r.IsConcurrent() && len(c.reversed.Nodes[r.ID()].Edges) == 0:
			ch <- r
			continue
		// Resource is not concurrent
		default:
			process(r)
		}
	}

	close(ch)
	wg.Wait()

	// Print summary report
	if !c.config.DryRun {
		var changed, failed, uptodate int
		for _, err := range c.status.items {
			switch err {
			case nil:
				changed++
			case resource.ErrInSync:
				uptodate++
			default:
				failed++
			}
		}
		c.config.Logger.Printf("Resource summary is %d up-to-date, %d changed, %d failed\n", uptodate, changed, failed)
	}

	return nil
}

// execute processes a single resource
func (c *Catalog) execute(r resource.Resource) error {
	// Check if the resource has failed dependencies
	for _, dep := range r.Dependencies() {
		if c.status.hasFailed(dep) {
			return fmt.Errorf("failed dependency for %s", dep)
		}
	}

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

	id := r.ID()
	var action func() error
	switch {
	case want.IsInList(present) && current.IsInList(absent):
		action = r.Create
		c.config.Logger.Printf("%s is %s, should be %s\n", id, current, want)
	case want.IsInList(absent) && current.IsInList(present):
		action = r.Delete
		c.config.Logger.Printf("%s is %s, should be %s\n", id, current, want)
	case state.Outdated:
		action = r.Update
		c.config.Logger.Printf("%s is out of date\n", id)
	default:
		c.config.Logger.Printf("%s is up to date\n", id)
		return resource.ErrInSync
	}

	return action()
}

// luaLen returns the number of unsorted resources in catalog.
// This method is called from Lua.
func (c *Catalog) luaLen() int {
	return len(c.Unsorted)
}
