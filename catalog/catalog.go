package catalog

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/dnaeon/gru/graph"
	"github.com/dnaeon/gru/resource"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

var ErrEmptyCatalog = errors.New("Catalog is empty")

type resourceItem struct {
	// Position of the resource declaration in HCL
	position string

	// An instantiated resource after calling it's provider
	r resource.Resource
}

// resourceJsonItem type represents a single resource declaration
// represented in JSON.
// Keys of the map are the resource types, e.g. "package", "service", etc.
// Value of each key is the actual instantiated resource
type resourceJsonItem map[string]resource.Resource

// Catalog type contains resources loaded from a given HCL or JSON input
type Catalog struct {
	// Contains resources to be serialized in JSON
	ResourceJsonItems []resourceJsonItem `json:"resource"`

	// Map containing the unique resource ids of instantiated resources.
	// The map is used to keep track of where resources were declared,
	// for detecting possibly duplicate resource declarations,
	// for building the resource dependency graph and perform
	// topological sort in order to determine the proper order of
	// evaluating the resources
	resourceIdMap map[string]resourceItem
}

// newCatalog creates a new empty catalog
func newCatalog() *Catalog {
	c := &Catalog{
		ResourceJsonItems: make([]resourceJsonItem, 0),
		resourceIdMap:     make(map[string]resourceItem),
	}

	return c
}

// resourceIsRegistered returns true if the resource id already exists in the catalog
func (c *Catalog) resourceIsRegistered(id string) bool {
	_, ok := c.resourceIdMap[id]

	return ok
}

// resourceGraph creates a DAG graph for the resources in catalog
func (c *Catalog) resourceGraph() (*graph.Graph, error) {
	// Create a DAG graph of the currently registered resources
	// The generated graph can be topologically sorted in order to
	// determine the proper order of evaluating resources
	// If the graph cannot be sorted, it means we have a
	// circular dependency in our resources
	g := graph.NewGraph()

	// A map containing the resource ids and their nodes in the graph
	nodes := make(map[string]*graph.Node)

	// Create the graph nodes for each resource from the catalog
	for name := range c.resourceIdMap {
		node := graph.NewNode(name)
		nodes[name] = node
		g.AddNode(node)
	}

	// Connect the nodes in the graph
	for name, ri := range c.resourceIdMap {
		deps := ri.r.Want()
		for _, dep := range deps {
			if !c.resourceIsRegistered(dep) {
				e := fmt.Errorf("Resource '%s' declared at %s wants '%s', which is not found in catalog", name, ri.position, dep)
				return g, e
			}
			g.AddEdge(nodes[name], nodes[dep])
		}
	}

	return g, nil
}

// Run processes the catalog
func (c *Catalog) Run() error {
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
		r := c.resourceIdMap[node.Name].r
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

// GenerateCatalogDOT generates a DOT file of the resources graph from catalog
func (c *Catalog) GenerateCatalogDOT(w io.Writer) error {
	if len(c.resourceIdMap) == 0 {
		return ErrEmptyCatalog
	}

	// Generate the graph for all registered resources
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
	return len(c.resourceIdMap)
}

// hclLoadResources loads all resource declarations from the given HCL input
func (c *Catalog) hclLoadResources(root *ast.ObjectList) error {
	hclResources := root.Filter("resource")
	for _, item := range hclResources.Items {
		position := item.Val.Pos().String()

		// The item is expected to have exactly one key which
		// represents the resource type.
		if len(item.Keys) != 1 {
			e := fmt.Errorf("Invalid resource declaration found at %s", position)
			return e
		}

		// Get the resource type and create the actual resource
		resourceType := item.Keys[0].Token.Value().(string)
		provider, ok := resource.Get(resourceType)
		if !ok {
			e := fmt.Errorf("Unknown resource type '%s' found at %s", resourceType, position)
			return e
		}

		// Create the actual resource by calling it's provider
		r, err := provider(item)
		if err != nil {
			return err
		}

		// Check if we have a duplicate resource declaration
		id := r.ID()
		if ri, ok := c.resourceIdMap[id]; ok {
			e := fmt.Errorf("Duplicate resource declaration for '%s' found at %s, previous declaration was at %s", id, position, ri.position)
			return e
		}

		// Add the new resource to the map of known and unique resources
		ri := resourceItem{
			position: position,
			r:        r,
		}
		c.resourceIdMap[id] = ri

		// Add the resource for JSON serialization as well
		rJson := resourceJsonItem{
			resourceType: r,
		}
		c.ResourceJsonItems = append(c.ResourceJsonItems, rJson)
	}

	return nil
}

// Load reads a catalog from the given HCL or JSON input and creates a catalog
func Load(path string) (*Catalog, error) {
	c := newCatalog()

	input, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	// Parse configuration
	obj, err := hcl.Parse(string(input))
	if err != nil {
		return c, err
	}

	// Top-level node should be an object list
	root, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return c, fmt.Errorf("Missing root node in %s", path)
	}

	err = c.hclLoadResources(root)
	if err != nil {
		return c, err
	}

	return c, nil
}
