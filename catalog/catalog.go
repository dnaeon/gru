package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

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

// newCatalog creates a new empty catalog
func newCatalog() *Catalog {
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
			id := r.ID()
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
		deps := r.Want()
		for _, dep := range deps {
			if _, ok := resources[dep]; !ok {
				e := fmt.Errorf("Resource %s wants %s, which is not in catalog", name, dep)
				return g, e
			}
			g.AddEdge(nodes[name], nodes[dep])
		}
	}

	return g, nil
}

// Run processes the catalog
func (c *Catalog) Run() error {
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

	for _, node := range sorted {
		r := rMap[node.Name]
		id := r.ID()
		state, err := r.Evaluate()
		if err != nil {
			log.Printf("Failed to evaluate resource '%s': %s", id, err)
			continue
		}

		if state.Want == state.Current {
			continue
		}

		log.Printf("%s is %s, should be %s", id, state.Current, state.Want)
		switch {
		case state.Want == resource.ResourceStatePresent && state.Current == resource.ResourceStateAbsent:
			r.Create()
		case state.Want == resource.ResourceStateAbsent && state.Current != resource.ResourceStateAbsent:
			r.Delete()
		case state.Want == resource.ResourceStateUpdate && state.Current == resource.ResourceStatePresent:
			r.Update()
		default:
			log.Printf("Unknown state '%s' for resource '%s'", state.Want, id)
		}
	}

	return nil
}

// GenerateCatalogDOT generates a DOT file for the resources in catalog
func (c *Catalog) GenerateCatalogDOT(w io.Writer) error {
	g, err := c.resourceGraph()
	if err != nil {
		return err
	}
	g.GenerateDOT("resources", w)

	// Try a topological sort of the graph
	// In case of circular dependencies in the graph
	// generate a DOT file for the remaining nodes in the graph,
	// which would give us the resources causing circular dependencies
	if nodes, err := g.Sort(); err == graph.ErrCircularDependency {
		circularGraph := graph.NewGraph()
		circularGraph.AddNode(nodes...)
		circularGraph.GenerateDOT("resources_circular", w)
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
	c := newCatalog()

	// Discover all modules from the provided module path
	registry, err := module.Discover(path)
	if _, ok := registry[main]; !ok {
		return c, fmt.Errorf("Module %s was not found in the module path", main)
	}

	// A map containing the module names and the actual loaded modules
	moduleNames := make(map[string]*module.Module)
	for n, p := range registry {
		m, err := module.Load(n, p)
		if err != nil {
			return c, err
		}
		moduleNames[n] = m
	}

	// A map containing the modules as graph nodes
	// The graph is used to determine if we have
	// circular module imports and also to provide the
	// proper ordering of loading modules after a
	// topological sort of the graph nodes
	nodes := make(map[string]*graph.Node)
	for n := range moduleNames {
		node := graph.NewNode(n)
		nodes[n] = node
	}

	// Recursively find all imports that the main module has and
	// resolve the dependency graph
	g := graph.NewGraph()
	var createModuleGraph func(m *module.Module) error
	createModuleGraph = func(m *module.Module) error {
		if !g.NodeExists(m.Name) {
			g.AddNode(nodes[m.Name])
		} else {
			return nil
		}

		for _, importName := range m.ModuleImport.Module {
			if _, ok := moduleNames[importName]; !ok {
				return fmt.Errorf("Module %s imports %s, which is not in the module path", m.Name, importName)
			}

			// Build the dependencies of imported modules as well
			createModuleGraph(moduleNames[importName])

			// Finally connect the nodes in the graph
			g.AddEdge(nodes[m.Name], nodes[importName])
		}

		return nil
	}

	//	Build the dependency graph of the module imports
	err = createModuleGraph(moduleNames[main])
	if err != nil {
		return c, err
	}

	// Topologically sort the graph
	// In case of an error it means we have a circular import
	sorted, err := g.Sort()
	if err != nil {
		return c, err
	}

	// Finally add the sorted modules to the catalog
	for _, node := range sorted {
		c.modules = append(c.modules, moduleNames[node.Name])
	}

	return c, nil
}

// MarshalJSON creates a stripped down version of the catalog in JSON,
// which contains all resources from the catalog and is suitable for
// clients to consume in order to create a single-module catalog from it.
func (c *Catalog) MarshalJSON() ([]byte, error) {
	rMap, err := c.createResourceMap()
	if err != nil {
		return nil, err
	}

	resources := make([]resourceMap, 0)
	for _, r := range rMap {
		rJson := resourceMap{
			r.Type(): r,
		}
		resources = append(resources, rJson)
	}

	return json.Marshal(map[string]interface{}{
		"resource": resources,
	})
}
