package catalog

import (
	"errors"
	"fmt"
	"io"

	"github.com/dnaeon/gru/graph"
	"github.com/dnaeon/gru/module"
	"github.com/dnaeon/gru/resource"
)

// ErrEmptyCatalog is returned when no resources were found from the
// loaded modules in the catalog
var ErrEmptyCatalog = errors.New("Catalog is empty")

type resourceMap map[string]resource.Resource

// Catalog type represents a collection of modules loaded from HCL or JSON
type Catalog struct {
	modules []*module.Module
}

// NewCatalog creates a new empty catalog
func NewCatalog() *Catalog {
	c := &Catalog{
		modules: make([]*module.Module, 0),
	}

	return c
}

// createResourceMap creates a map of the unique resource IDs and
// the actual resource instances
func (c *Catalog) createResourceMap() (resourceMap, error) {
	// A map containing the unique resource ID and the
	// module where the resource has been declared
	rModuleMap := make(map[string]string)

	rMap := make(resourceMap)
	for _, m := range c.modules {
		for _, r := range m.Resources {
			id := r.ResourceID()
			if _, ok := rMap[id]; ok {
				return rMap, fmt.Errorf("Duplicate resource %s in %s, previous declaration was in %s", id, m.Name, rModuleMap[id])
			}
			rModuleMap[id] = m.Name
			rMap[id] = r
		}
	}

	if len(rMap) == 0 {
		return rMap, ErrEmptyCatalog
	}

	return rMap, nil
}

// resourceGraph creates a DAG graph for the resources in catalog
func (c *Catalog) resourceGraph() (*graph.Graph, error) {
	// Create a DAG graph of the resources in catalog
	// The generated graph can be topologically sorted in order to
	// determine the proper order of evaluating resources
	// If the graph cannot be sorted, it means we have a
	// circular dependency in our resources
	g := graph.NewGraph()

	resources, err := c.createResourceMap()
	if err != nil {
		return g, err
	}

	// A map containing the resource ids and their nodes in the graph
	// Create a graph nodes for each resource from the catalog
	nodes := make(map[string]*graph.Node)
	for name := range resources {
		node := graph.NewNode(name)
		nodes[name] = node
		g.AddNode(node)
	}

	// Connect the nodes in the graph
	for name, r := range resources {
		before := r.WantBefore()
		after := r.WantAfter()

		// Connect current resource with the ones that happen after it
		for _, dep := range after {
			if _, ok := resources[dep]; !ok {
				e := fmt.Errorf("Resource %s wants %s, which is not in catalog", name, dep)
				return g, e
			}
			g.AddEdge(nodes[name], nodes[dep])
		}

		// Connect current resource with the ones that happen before it
		for _, dep := range before {
			if _, ok := resources[dep]; !ok {
				e := fmt.Errorf("Resource %s wants %s, which is not in catalog", name, dep)
				return g, e
			}
			g.AddEdge(nodes[dep], nodes[name])
		}
	}

	return g, nil
}

// Run processes the catalog
func (c *Catalog) Run(w io.Writer) error {
	resourceErrors := c.Validate()
	if len(resourceErrors) > 0 {
		for _, e := range resourceErrors {
			fmt.Fprint(w, e)
		}
		return errors.New("Failed to validate catalog resources")
	}

	rMap, err := c.createResourceMap()
	if err != nil {
		return err
	}

	// Perform topological sort of the resources graph
	g, err := c.resourceGraph()
	if err != nil {
		return err
	}

	sorted, err := g.Sort()
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Loaded %d resources from catalog\n", c.Len())
	for _, node := range sorted {
		r := rMap[node.Name]
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
				err = r.Update(w)
				if err != nil {
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
		err = action(w)
		if err != nil {
			fmt.Fprintf(w, "%s error: %s", r.ResourceID(), err)
		}

		if state.Update {
			err = r.Update(w)
			if err != nil {
				fmt.Fprintf(w, "%s error %s\n", r.ResourceID(), err)
			}
		}
	}

	return nil
}

// Validate validates the resources from catalog
func (c *Catalog) Validate() []error {
	var foundErrors []error

	rMap, err := c.createResourceMap()
	if err != nil {
		foundErrors = append(foundErrors, err)
	}

	// Validate resources
	for id, r := range rMap {
		err = r.Validate()
		if err != nil {
			foundErrors = append(foundErrors, fmt.Errorf("Failed to validate %s: %s\n", id, err))
		}
	}

	// Check for unknown keys
	for _, m := range c.modules {
		for _, key := range m.UnknownKeys {
			foundErrors = append(foundErrors, fmt.Errorf("Unknown key '%s' in module %s\n", key, m.Name))
		}
	}

	return foundErrors
}

// ResourcesAsDot generates a DOT representation for the resources in catalog
func (c *Catalog) ResourcesAsDot(w io.Writer) error {
	g, err := c.resourceGraph()
	if err != nil {
		return err
	}
	g.AsDot("resources", w)

	// Try a topological sort of the graph
	// In case of circular dependencies in the graph
	// generate a DOT file for the remaining nodes in the graph,
	// which would give us the resources causing circular dependencies
	if nodes, err := g.Sort(); err == graph.ErrCircularDependency {
		circularGraph := graph.NewGraph()
		circularGraph.AddNode(nodes...)
		circularGraph.AsDot("resources_circular", w)
	}

	return nil
}

// Len returns the number of unique resources found in catalog
func (c *Catalog) Len() int {
	resources, err := c.createResourceMap()
	if err != nil {
		return 0
	}

	return len(resources)
}

// Load creates a catalog from the provided module name
func Load(main, path string) (*Catalog, error) {
	c := NewCatalog()

	modules, err := module.DiscoverAndLoad(path)
	if err != nil {
		return c, err
	}

	// Get the import graph of the module and sort it
	g, err := module.ImportGraph(main, path)
	if err != nil {
		return c, err
	}

	sorted, err := g.Sort()
	if err != nil {
		return c, err
	}

	// Finally add the sorted modules to the catalog
	for _, node := range sorted {
		c.modules = append(c.modules, modules[node.Name])
	}

	return c, nil
}
