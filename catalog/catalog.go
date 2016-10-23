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
	status *Status `luar:"-"`

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

// Status type contains status information about processed resources.
type Status struct {
	sync.RWMutex

	// Items contain the status for resources after being processed.
	Items map[string]*StatusItem
}

// StatusItem type represents a single item for a processed resource.
type StatusItem struct {
	// StateChanged field specifies whether or not a resource has changed
	// after being evaluated and processed.
	StateChanged bool

	// Err contains any errors that were encountered during resource
	// evaluation and processing.
	Err error
}

// Summary displays a summary of the resource status.
func (s *Status) Summary(l *log.Logger) {
	s.Lock()
	defer s.Unlock()

	var changed, failed, uptodate int
	for _, item := range s.Items {
		switch {
		case item.StateChanged == true && item.Err == nil:
			changed++
		case item.StateChanged == false && item.Err == nil:
			uptodate++
		default:
			failed++
		}
	}

	l.Printf("%d up-to-date, %d changed, %d failed\n", uptodate, changed, failed)
}

// New creates a new empty catalog with the provided configuration
func New(config *Config) *Catalog {
	c := &Catalog{
		config:     config,
		collection: make(resource.Collection),
		sorted:     make([]*graph.Node, 0),
		reversed:   graph.New(),
		status: &Status{
			Items: make(map[string]*StatusItem),
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
func (c *Catalog) Run() *Status {
	// process executes a single resource
	process := func(r resource.Resource) {
		id := r.ID()
		item := c.execute(r)
		c.status.Lock()
		defer c.status.Unlock()
		c.status.Items[id] = item
		if item.Err != nil {
			c.config.Logger.Printf("%s %s\n", id, item.Err)
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

	return c.status
}

// execute processes a single resource
func (c *Catalog) execute(r resource.Resource) *StatusItem {
	if err := c.hasFailedDependencies(r); err != nil {
		return &StatusItem{Err: err}
	}

	if err := r.Validate(); err != nil {
		return &StatusItem{Err: err}
	}

	if err := r.Initialize(); err != nil {
		return &StatusItem{Err: err}
	}
	defer r.Close()

	state, err := r.Evaluate()
	if err != nil {
		return &StatusItem{Err: err}
	}

	if c.config.DryRun {
		return &StatusItem{}
	}

	// Current and wanted states for the resource
	want := utils.NewString(state.Want)
	current := utils.NewString(state.Current)

	// The list of present and absent states for the resource
	present := utils.NewList(r.GetPresentStates()...)
	absent := utils.NewList(r.GetAbsentStates()...)

	// Process resource
	id := r.ID()
	var action func() error
	switch {
	case want.IsInList(present) && current.IsInList(absent):
		action = r.Create
		c.config.Logger.Printf("%s is %s, should be %s\n", id, current, want)
	case want.IsInList(absent) && current.IsInList(present):
		action = r.Delete
		c.config.Logger.Printf("%s is %s, should be %s\n", id, current, want)
	default:
		// No-op: resource is in sync
	}

	stateChanged := false
	if action != nil {
		stateChanged = true
		if err := action(); err != nil {
			return &StatusItem{StateChanged: true, Err: err}
		}
	}

	// Process resource properties
	for _, p := range r.Properties() {
		synced, err := p.IsSynced()
		if err != nil {
			// Some properties make no sense if the resource is absent, e.g.
			// setting up file permissions requires that the file managed by the
			// resource is present, therefore we ignore errors for properties
			// which make no sense if the resource is absent.
			if err == resource.ErrResourceAbsent {
				continue
			}
			e := fmt.Errorf("unable to evaluate property %s: %s\n", p.Name, err)
			return &StatusItem{StateChanged: true, Err: e}
		}

		if !synced {
			stateChanged = true
			c.config.Logger.Printf("%s property '%s' is out of date\n", id, p.Name())
			if err := p.Set(); err != nil {
				e := fmt.Errorf("unable to set property %s: %s\n", p.Name, err)
				return &StatusItem{StateChanged: true, Err: e}
			}
		}
	}

	if err := c.runTriggers(r); err != nil {
		return &StatusItem{StateChanged: stateChanged, Err: err}
	}

	return &StatusItem{StateChanged: stateChanged, Err: nil}
}

// runTriggers executes the triggers for each
// monitored resource if it's state has changed
func (c *Catalog) runTriggers(r resource.Resource) error {
	c.status.Lock()
	defer c.status.Unlock()

	for subscribed, trigger := range r.SubscribedTo() {
		item := c.status.Items[subscribed]
		if !item.StateChanged {
			continue
		}

		c.config.Logger.Printf("%s running trigger, because %s has changed\n", r.ID(), subscribed)
		c.config.L.Push(trigger)
		if err := c.config.L.PCall(0, 0, nil); err != nil {
			c.config.Logger.Printf("%s trigger exited with an error: %s\n", r.ID(), err)
			return err
		}
	}

	return nil
}

// hasFailedDependencies checks if a resource has failed dependencies.
func (c *Catalog) hasFailedDependencies(r resource.Resource) error {
	c.status.Lock()
	defer c.status.Unlock()

	for _, dep := range r.Dependencies() {
		item := c.status.Items[dep]
		if item.Err != nil {
			return fmt.Errorf("failed dependency for %s", dep)
		}
	}

	return nil
}

// luaLen returns the number of unsorted resources in catalog.
// This method is called from Lua.
func (c *Catalog) luaLen() int {
	return len(c.Unsorted)
}
